package netlinklistlistens

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
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
)

type inetDiagReqV2 C.inet_diag_req_v2

func (req *inetDiagReqV2) pack() []byte {
	return (*(*[C.sizeof_inet_diag_req_v2]byte)(unsafe.Pointer(req)))[:]
}

func (req *inetDiagReqV2) Len() int {
	return C.sizeof_inet_diag_req_v2
}

func (req *inetDiagReqV2) Serialize() []byte {
	return req.pack()
}

type inetDiagMsg C.inet_diag_msg

func unpackInetDiagMsg(b []byte) *inetDiagMsg {
	return (*inetDiagMsg)(unsafe.Pointer(&b[0]))
}

func htons(v C.ushort) uint16 {
	b := *(*[2]byte)(unsafe.Pointer(&v))
	return binary.BigEndian.Uint16(b[:])
}

func (msg *inetDiagMsg) src() (*netip.AddrPort, error) {
	src := msg.id.idiag_src

	var ip netip.Addr
	switch msg.idiag_family {
	case syscall.AF_INET:
		{
			b := *(*[4]byte)(unsafe.Pointer(&src[0]))
			ip = netip.AddrFrom4(b)
		}
	case syscall.AF_INET6:
		{
			b := *(*[16]byte)(unsafe.Pointer(&src[0]))
			ip = netip.AddrFrom16(b)
		}
	default:
		return nil, errors.New(fmt.Sprintf("unexpected family %d", msg.idiag_family))
	}

	sport := msg.id.idiag_sport
	addr := netip.AddrPortFrom(ip, htons(sport))
	return &addr, nil
}

func inetDiag(cb func(*inetDiagMsg) error) error {
	families := []C.uchar{
		syscall.AF_INET,
		syscall.AF_INET6,
	}

	for _, family := range families {
		req := nl.NewNetlinkRequest(nl.SOCK_DIAG_BY_FAMILY, syscall.NLM_F_DUMP)
		req.AddData(&inetDiagReqV2{
			sdiag_family:   C.uchar(family),
			sdiag_protocol: syscall.IPPROTO_TCP,
			idiag_states:   1 << C.TCP_LISTEN,
		})
		res, err := req.Execute(syscall.NETLINK_INET_DIAG, 0)
		if err != nil {
			return err
		}

		for _, data := range res {
			msg := unpackInetDiagMsg(data)
			if err := cb(msg); err != nil {
				return err
			}
		}
	}

	return nil
}

func ListListens() ([]netip.AddrPort, error) {
	ret := []netip.AddrPort{}
	err := inetDiag(func(msg *inetDiagMsg) error {
		src, err := msg.src()
		if err != nil {
			return err
		}
		ret = append(ret, *src)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
