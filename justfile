test:
	go test ./... -v

install: build_cli copy_files
	cd cmd/aicmd && go install

copy_files:
  ./scripts/install.sh

build_cli:
  ./scripts/increment_version.sh
  go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/aicmd/main.go)" -o aicmdtools cmd/aicmd/main.go
