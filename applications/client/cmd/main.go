package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/DarthPestilane/easytcp"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/labstack/echo/v4"
	"github.com/oklog/run"
	stdprom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/donmikel/wow/applications/client"
	"github.com/donmikel/wow/applications/client/adapters/tcp/server"
	"github.com/donmikel/wow/applications/client/config"
	"github.com/donmikel/wow/applications/client/interfaces"
	"github.com/donmikel/wow/applications/client/service"
)

// exitCode is a process termination code.
type exitCode int

// Possible process termination codes are listed below.
const (
	// exitSuccess is code for successful program termination.
	exitSuccess exitCode = 0
	// exitFailure is code for unsuccessful program termination.
	exitFailure exitCode = 1
)

// Constants for prometheus metrics names. The full name is built like this:
// <namespace>_<subsystem>_<name>.
const (
	PrometheusNamespace = "wow"
	PrometheusSubsystem = "client"
)

// Kubernetes (rolling update) doesn't wait until a pod is out of rotation before sending SIGTERM,
// and external LB could still route traffic to a non-existing pod resulting in a surge of 50x API errors.
// It's recommended to wait for 5 seconds before terminating the program; see references
// https://github.com/kubernetes-retired/contrib/issues/1140, https://youtu.be/me5iyiheOC8?t=1797.
const preStopWait = 5 * time.Second

// Shutdown timeout for http servers.
const shutdownTimeout = 5 * time.Second

var (
	// version is the service version from git tag.
	version = ""
)

func main() {
	os.Exit(int(gracefulMain()))
}

// gracefulMain releases resources gracefully upon termination.
// When we call os.Exit defer statements do not run resulting in unclean process shutdown.
// nolint
func gracefulMain() exitCode {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	configPath := fs.String("config", "", "path to the config file")
	v := fs.Bool("v", false, "Show version")
	cpuprofile := fs.String("cpuprofile", "", "write cpu profile to `file`")

	err := fs.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		return exitSuccess
	}
	if err != nil {
		logger.Log("msg", "parsing cli flags failed", "err", err)
		return exitFailure
	}

	if *v {
		if version == "" {
			level.Error(logger).Log("Version not set")
		} else {
			level.Info(logger).Log("Version: %s\n", version)
		}

		return exitSuccess
	}

	logger.Log("cpuprofile", *cpuprofile)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			logger.Log("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Log("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	logger.Log("configPath", *configPath)

	cfg, err := config.Parse(*configPath)
	if err != nil {
		logger.Log("msg", "cannot parse service config", "err", err)
		return exitFailure
	}

	err = cfg.Validate()
	if err != nil {
		logger.Log("msg", "config validation failed", "err", err)
		return exitFailure
	}

	// It's nice to be able to see panics in Logs, hence we monitor for panics after
	// logger has been bootstrapped.
	defer monitorPanic(logger)

	eOps := echo.New()
	eOps.HideBanner = true
	// Expose endpoint for healthcheck.
	eOps.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	// Expose the registered Prometheus metrics via HTTP.
	eOps.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	var g run.Group
	{
		ctx, cancel := context.WithCancel(context.Background())
		{
			g.Add(func() error {
				conn, err := net.Dial("tcp", cfg.Server.TCPAddr)
				defer conn.Close()

				if err != nil {
					return fmt.Errorf("tcp server connection error: %w", err)
				}

				var svc client.Service
				{
					var si interfaces.Server
					si, err = server.NewProvider(
						conn,
						easytcp.NewDefaultPacker(),
						&easytcp.JsonCodec{},
						logger,
					)
					if err != nil {
						return fmt.Errorf("server provider creation error: %w", err)
					}

					serverHistMetric := kitprom.NewHistogramFrom(stdprom.HistogramOpts{
						Namespace: PrometheusNamespace,
						Subsystem: PrometheusSubsystem,
						Name:      "server_duration_seconds",
						Help:      "Server requests duration in seconds by method.",
					}, []string{"method", "error"})

					si = server.NewInstrumentingMiddleware(serverHistMetric, si)

					svc = service.NewService(si, easytcp.NewDefaultPacker(), &easytcp.JsonCodec{}, logger)
				}

				go func() {
					ticker := time.NewTicker(10 * time.Second)
					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:
							quote, err := svc.GetQuote()
							if err != nil {
								level.Error(logger).Log("msg", "get quote error: %w", err)
								continue
							}
							level.Info(logger).Log("msg", "quote received",
								"quote", quote,
							)
						case <-ctx.Done():
							return
						}
					}
				}()

				<-ctx.Done()
				return nil
			}, func(_ error) {
				level.Info(logger).Log("msg", "tcp server was interrupted")
				cancel()
			})
		}
		{
			g.Add(func() error {
				sig := make(chan os.Signal, 1)
				signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
				select {
				case <-ctx.Done():
					return nil
				case s := <-sig:
					level.Info(logger).Log("msg", fmt.Sprintf("signal received (waiting %v before terminating): %v", preStopWait, s))
					time.Sleep(preStopWait)
					level.Info(logger).Log("msg", "terminating...")

					return nil
				}
			}, func(_ error) {
				level.Info(logger).Log("msg", "program was interrupted")
				cancel()
			})
		}
		{
			g.Add(func() error {
				level.Info(logger).Log("msg", "ops server is starting", "addr", cfg.Ops.HTTPAddr)

				err = eOps.Start(cfg.Ops.HTTPAddr)
				if err != nil {
					return fmt.Errorf("ops server start error: %w", err)
				}

				return nil
			}, func(_ error) {
				level.Info(logger).Log("msg", "ops server was interrupted")
				ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
				defer cancel()

				shErr := eOps.Shutdown(ctx)
				if shErr != nil {
					level.Error(logger).Log("msg", "ops server shut down with error", "err", shErr)
				}
			})
		}

	}
	err = g.Run()
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("actors stopped with err: %v", err))
		return exitFailure
	}

	level.Info(logger).Log("msg", "actors stopped without errors")

	return exitSuccess
}

// monitorPanic monitors panics and reports them somewhere (e.g. logs, ...).
func monitorPanic(logger log.Logger) {
	if rec := recover(); rec != nil {
		err := fmt.Sprintf("panic: %v \n stack trace: %s", rec, debug.Stack())
		level.Error(logger).Log("err", err)
		panic(err)
	}
}
