package main

/*
#include <linux/netlink.h>
#include <linux/inet_diag.h>
#include <netinet/tcp.h>
#include <arpa/inet.h>

typedef struct inet_diag_req_v2 inet_diag_req_v2;
typedef struct inet_diag_msg inet_diag_msg;
*/
import "C"

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
)

type inet_diag_req_v2 C.inet_diag_req_v2

func (req *inet_diag_req_v2) Len() int {
	return C.sizeof_inet_diag_req_v2
}

func (req *inet_diag_req_v2) Serialize() []byte {
	return (*(*[C.sizeof_inet_diag_req_v2]byte)(unsafe.Pointer(req)))[:]
}

func inetDiag(cb func(*C.inet_diag_msg)) error {
	families := []C.uchar{
		syscall.AF_INET,
		syscall.AF_INET6,
	}

	for _, family := range families {
		req := nl.NewNetlinkRequest(nl.SOCK_DIAG_BY_FAMILY, syscall.NLM_F_DUMP)
		req.AddData(&inet_diag_req_v2{
			sdiag_family:   C.uchar(family),
			sdiag_protocol: syscall.IPPROTO_TCP,
			idiag_states:   1 << C.TCP_LISTEN,
		})
		res, err := req.Execute(syscall.NETLINK_INET_DIAG, 0)
		if err != nil {
			return err
		}

		for _, data := range res {
			msg := (*C.inet_diag_msg)(unsafe.Pointer(&data[0]))
			cb(msg)
		}
	}

	return nil
}

func dump(msg *C.inet_diag_msg) {
	sport := C.ntohs(msg.id.idiag_sport)

	ip := make([]byte, 40)
	ret := C.inet_ntop(C.int(msg.idiag_family), unsafe.Pointer(&msg.id.idiag_src), (*C.char)(unsafe.Pointer(&ip)), C.uint(len(ip)))
	if ret == nil {
		ip = []byte("ERROR")
	}

	if msg.idiag_family == syscall.AF_INET {
		fmt.Printf("%s:%d\n", C.GoString((*C.char)(unsafe.Pointer(&ip))), sport)
	} else {
		fmt.Printf("[%s]:%d\n", C.GoString((*C.char)(unsafe.Pointer(&ip))), sport)
	}
}

func main() {
	if err := inetDiag(dump); err != nil {
		log.Fatal(err)
	}
}
