NAME    = crec
LDFLAGS = -ldflags="-s -w"
GOOS    = linux
GOARCH  = arm64

.PHONY: clean
clean:
	rm -rf bin/* Gopkg.lock

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/$(NAME) main.go;
