package main

import (
	"fmt"
	"log"
	"net/netip"
	"time"

	"github.com/yskszk63/netlink-list-listens"
)

func update(m *map[netip.AddrPort]struct{}) ([]netip.AddrPort, []netip.AddrPort, error) {
	l, err := netlinklistlistens.ListListens()
	if err != nil {
		return nil, nil, err
	}

	old := *m
	new := make(map[netip.AddrPort]struct{})
	*m = new
	add := make([]netip.AddrPort, 0)
	rm := make([]netip.AddrPort, 0)

	for _, addr := range l {
		new[addr] = struct{}{}

		_, exists := old[addr]
		if !exists {
			add = append(add, addr)
			continue
		}
		delete(old, addr)
	}

	for k := range old {
		rm = append(rm, k)
	}

	return add, rm, nil
}

func main() {
	m := make(map[netip.AddrPort]struct{})

	d := time.Second * 1
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		add, rm, err := update(&m)
		if err != nil {
			log.Fatal(err)
		}

		for _, a := range add {
			fmt.Printf("ADD\t%s\n", a)
		}

		for _, r := range rm {
			fmt.Printf("REMOVE\t%s\n", r)
		}

		<- ticker.C
	}
}
