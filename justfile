test:
	go test ./... -v

install: build_cli copy_files
	cd cmd/goai && go install

copy_files:
  ./install.sh

build_cli:
  ./increment_version.sh
  go build -ldflags "-X main.version=$(grep -oP 'version = "\K[^"]+' ./cmd/goai/main.go)" -o goai cmd/goai/main.go
