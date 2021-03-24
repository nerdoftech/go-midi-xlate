
X32_BIN = x32-delay
MAC_DIR = build/mac/
X32_SRC = ./cmd/x32-delay

test:
	go test -v -cover ./pkg/readhandlers/
	go test -v -cover ./pkg/x32/

build:
	pwd
	mkdir -vp $(MAC_DIR)
	env GOARCH=amd64 GOOS=darwin go build -o $(MAC_DIR)$(X32_BIN)-macos-64bit $(X32_SRC)

clean:
	rm -vrf build