MIN_COVERAGE = 50

test:
	go test $$(go list ./...) -race -count=1 -cover -coverprofile=coverage.txt.tmp && cat coverage.txt.tmp | grep -v ".gen.go" > coverage.txt && go tool cover -func=coverage.txt \
	| grep total | tee /dev/stderr | sed 's/\%//g' | awk '{err=0;c+=$$3}{if (c > 0 && c < $(MIN_COVERAGE)) {printf "=== FAIL: Coverage failed at %.2f%%\n", c; err=1}} END {exit err}'

format:
	gci write $$(find . -type f -name '*.go' -not -path "./pkg/proto/*" -not -name "*.gen.go" -not -path "*/mock/*") -s "Standard,Default,Prefix(github.com/donmikel/wow)"

generate:
	go generate ./...

lint:
	golangci-lint run --deadline=5m -v

gosec:
	gosec -exclude-generated -exclude=G104 -fmt=json -exclude-dir=.go ./...

lint_docker:
	docker run --rm -v $(GOPATH)/pkg/mod:/go/pkg/mod:ro -v `pwd`:/`pwd`:ro -w /`pwd` golangci/golangci-lint:v1.46.2-alpine golangci-lint run --fix --deadline=5m -v

build: build_server build_client build_guide

build_docker: build_server_docker build_client_docker build_guide_docker

build_server:
	go build --tags netcgo -o ./bin/server ./applications/server/cmd/

build_client:
	go build --tags netcgo -o ./bin/client ./applications/client/cmd/

build_guide:
	go build --tags netcgo -o ./bin/guide ./applications/guide/cmd/

build_server_docker:
	docker build --tag=server:latest --file=docker/Dockerfile.server .

build_client_docker:
	docker build --tag=client:latest --file=docker/Dockerfile.client .

build_guide_docker:
	docker build --tag=guide:latest --file=docker/Dockerfile.guide .

up:
	docker-compose -f docker/docker-compose.yml up -d --build

down:
	docker-compose -f docker/docker-compose.yml down
