build:
	go build -v

run: build
	./neptune

update_mod:
	go build -v -mod=mod

# (build but with a smaller binary)
dist:
	go build -ldflags="-w -s" -gcflags=all=-l -v

# (even smaller binary)
pack: dist
	upx ./neptune
