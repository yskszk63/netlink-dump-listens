.PHONY: build-examples
build-examples: netlink-list-listens watch-listens

netlink-list-listens: examples/netlink-list-listens/main.go netlink-list-listens.go go.mod go.sum
	go build -a -tags netgo -installsuffix netgo -ldflags='-s -w -extldflags "-static"' -o=$@ ./examples/netlink-list-listens

watch-listens: examples/watch-listens/main.go netlink-list-listens.go go.mod go.sum
	go build -a -tags netgo -installsuffix netgo -ldflags='-s -w -extldflags "-static"' -o=$@ ./examples/watch-listens

.PHONY: clean
clean:
	$(RM) netlink-list-listens watch-listens
