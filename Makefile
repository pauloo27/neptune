BINARY_NAME = neptune 

build:
	go build -v

# requires a valid cross-compiling environment
windows:
	PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/lib/pkgconfig CGO_ENABLED=1 CC=x86_64-w64-mingw32-cc GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -v

windows_pack: windows
	upx ./$(BINARY_NAME).exe

run: build
	./$(BINARY_NAME) 

install: build
	sudo cp ./$(BINARY_NAME) /usr/bin/

update_mod:
	go build -v -mod=mod

# (build but with a smaller binary)
dist:
	go build -ldflags="-w -s" -gcflags=all=-l -v

# (even smaller binary)
pack: dist
	upx ./$(BINARY_NAME)
