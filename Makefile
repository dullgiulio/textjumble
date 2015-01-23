GOPATH := $(CURDIR):$(GOPATH)

default: new quality build

build:
	mkdir -p bin
	go build -o bin/tj tj

clean:
	rm -rf bin

quality: fmt vet lint

fmt:
	go fmt tj 

vet:
	go vet tj

lint: bin/golint
	bin/golint src/tj/*.go

bin/golint:
	go get golang.org/x/tools/cmd/vet
	go build -o bin/golint golang.org/x/tools/cmd/vet 
