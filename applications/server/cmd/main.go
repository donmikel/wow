package main

import (
	"context"
	"flag"
	"fmt"
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

	"github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/applications/server/adapters/inmemory"
	"github.com/donmikel/wow/applications/server/config"
	"github.com/donmikel/wow/applications/server/handlers/challenge"
	"github.com/donmikel/wow/applications/server/handlers/quotes"
	"github.com/donmikel/wow/applications/server/interfaces"
	"github.com/donmikel/wow/applications/server/service"
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
	PrometheusSubsystem = "server"
)

// Kubernetes (rolling update) doesn't wait until a pod is out of rotation before sending SIGTERM,
// and external LB could still route traffic to a non-existing pod resulting in a surge of 50x API errors.
// It's recommended to wait for 5 seconds before terminating the program; see references
// https://github.com/kubernetes-retired/contrib/issues/1140, https://youtu.be/me5iyiheOC8?t=1797.
const preStopWait = 5 * time.Second

// Shutdown timeout for http servers.
const shutdownTimeout = 5 * time.Second

const guideTourLength = 5

// Truncation value should be adjusted to the average time of the tour.
const timeStampTruncationInterval = 1 * time.Minute

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
		// Create a new server with options.
		s := easytcp.NewServer(&easytcp.ServerOption{
			Packer: easytcp.NewDefaultPacker(), // use default packer
			Codec:  &easytcp.JsonCodec{},
		})

		ctx, cancel := context.WithCancel(context.Background())
		{
			g.Add(func() error {
				var svc server.Service
				{
					var wowRepository interfaces.WOWRepository
					quoutes := []string{
						"The fool doth think he is wise, but the wise man knows himself to be a fool.",
						"It is better to remain silent at the risk of being thought a fool, than to talk and remove all doubt of it.",
						"Whenever you find yourself on the side of the majority, it is time to reform (or pause and reflect).",
						"When someone loves you, the way they talk about you is different. You feel safe and comfortable.",
						"Knowing yourself is the beginning of all wisdom.",
						"The only true wisdom is in knowing you know nothing.",
						"The saddest aspect of life right now is that science gathers knowledge faster than society gathers wisdom.",
					}
					wowRepository = inmemory.NewWOWRepository(quoutes)

					svc = service.NewService(wowRepository)
				}

				guideKeys := make([]string, len(cfg.Guides))
				guides := make([]string, len(cfg.Guides))
				for i, guide := range cfg.Guides {
					guideKeys[i] = guide.Key
					guides[i] = guide.Host
				}

				serverHistMetric := kitprom.NewHistogramFrom(stdprom.HistogramOpts{
					Namespace: PrometheusNamespace,
					Subsystem: PrometheusSubsystem,
					Name:      "server_duration_seconds",
					Help:      "Server requests duration in seconds by method.",
				}, []string{"method"})

				getTourLength := func() int {
					return guideTourLength
				}
				getTimeStamp := func() int64 {
					return time.Now().UTC().Truncate(timeStampTruncationInterval).Unix()
				}

				var challengeTcpServer challenge.ServerInterface
				var quotesTcpServer quotes.ServerInterface
				{
					challengeTcpServer = challenge.NewTCPServer(guides, guideKeys, getTourLength, getTimeStamp, logger)
					challengeTcpServer = challenge.NewInstrumentingMiddleware(serverHistMetric, challengeTcpServer)
					challenge.RegisterHandlers(challengeTcpServer, s)

					quotesTcpServer = quotes.NewTCPServer(svc, logger)
					quotesTcpServer = quotes.NewInstrumentingMiddleware(serverHistMetric, quotesTcpServer)

					pow := challenge.NewPOWMiddleware(guideKeys, getTourLength, getTimeStamp, logger)
					quotes.RegisterHandlers(quotesTcpServer, s, pow.GuideTour)
				}

				// Listen and serve.
				if err := s.Run(cfg.TCP.TCPAddr); err != nil && err != easytcp.ErrServerStopped {
					level.Error(logger).Log("msg", fmt.Sprintf("serve error: %v", err.Error()))
				}

				return nil
			}, func(_ error) {
				if err = s.Stop(); err != nil {
					level.Info(logger).Log("msg", "tcp server stopped err: %w", err)
				}

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
				ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
				defer cancel()

				shErr := eOps.Shutdown(ctx)
				if shErr != nil {
					level.Error(logger).Log("msg", "ops server shut down with error", "err", shErr)
				}

				level.Info(logger).Log("msg", "ops server was interrupted")
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
