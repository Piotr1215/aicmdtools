test:
	go test ./... -v

install: copy_files
	cd cmd/goai && go install

copy_files:
  ./install.sh

build_cli:
  go build cmd/goai/main.go

