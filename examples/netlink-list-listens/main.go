package main

import (
	"log"
	"fmt"

	"github.com/yskszk63/netlink-list-listens"
)

func main() {
        l, err := netlinklistlistens.ListListens()
        if err != nil {
                log.Fatal(err)
        }
        for _, addr := range l {
                fmt.Printf("%s\n", addr)
        }
}
