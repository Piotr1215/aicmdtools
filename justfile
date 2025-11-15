default:
  just --list

test:
  go test ./... -v

install: build_cli build_graph build_chat build_fix
  @cd cmd/aicmd && go install
  @cd cmd/aicompgraph && go install
  @cd cmd/aichat && go install
  @cd cmd/aifix && go install

copy_files:
  @./scripts/install.sh

build_graph:
  @./scripts/increment_version.sh
  @go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aicompgraph/main.go)" -o aicompgraph cmd/aicompgraph/main.go

build_cli:
  @./scripts/increment_version.sh
  @go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aicmd/main.go)" -o aicmd cmd/aicmd/main.go

build_chat:
  @./scripts/increment_version.sh
  @go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aichat/main.go)" -o aichat cmd/aichat/main.go

build_fix:
  @./scripts/increment_version.sh
  @go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aifix/main.go)" -o aifix cmd/aifix/main.go
