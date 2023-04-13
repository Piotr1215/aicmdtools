test:
	go test ./... -v

install: build_cli build_graph copy_files
  cd cmd/aicmd && go install
  cd cmd/aicompgraph && go install

copy_files:
  ./scripts/install.sh

build_graph:
  ./scripts/increment_version.sh
  go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aicompgraph/main.go)" -o aicompgraph cmd/aicompgraph/main.go

build_cli:
  ./scripts/increment_version.sh
  go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aicmd/main.go)" -o aicmd cmd/aicmd/main.go
