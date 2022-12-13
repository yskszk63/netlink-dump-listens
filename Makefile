.PHONY: build
build: netlink-list-listens

netlink-list-listens: main.go go.mod go.sum
	go build -a -tags netgo -installsuffix netgo -ldflags='-s -w -extldflags "-static"' -o=$@

.PHONY: clean
clean:
	$(RM) netlink-list-listens
