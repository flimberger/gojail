// SPDX-License-Identifier: BSD-2-Clause-FreeBSD
//
// Copyright (c) 2020 Florian Limberger <flo@purplekraken.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package gojail

import (
	"fmt"
	"net"
	"unsafe"

	"purplekraken.com/pkg/gojail/syscall"
)

// Opaque wrapper for syscall.JailStruct, so the user doesn't have to bother with
// the details.
type Jail struct {
	Path string
	Hostname string
	Jailname string
	IPs []net.IP
}

func rawbytes(s string) *byte {
	return (*byte)(unsafe.Pointer(&[]byte(s)[0]));
}

func NewJail(path, hostname, jailname string, ips []net.IP) (*Jail, error) {
	if len(path) == 0 || len(hostname) == 0 {
		return nil, fmt.Errorf("invalid parameter: path or hostname is nil")
	}
	return &Jail{
		Path: path,
		Hostname: hostname,
		Jailname: jailname,
		IPs: ips,
	}, nil
}

func (jail *Jail) Call() error {
	var ip4s []net.IP
	var ip6s []net.IP
	for _, ip := range jail.IPs {
		if ip4 := ip.To4(); ip4 != nil {
			ip4s = append(ip4s, ip4)
		} else if ip6 := ip.To16(); ip6 != nil {
			ip6s = append(ip6s, ip6)
		} else {
			panic("invalid IP address")
		}
	}
	nip4s := uint32(len(ip4s))
	nip6s := uint32(len(ip6s))
	var rawip4s *[4]byte
	var rawip6s *[16]byte
	if nip4s != 0 {
		rawip4s = (*[4]byte)(unsafe.Pointer(&ip4s[0]))
	}
	if nip6s != 0 {
		rawip6s = (*[16]byte)(unsafe.Pointer(&ip6s[0]))
	}
	j := syscall.JailStruct{
		Version: syscall.JAIL_API_VERSION,
		Path: rawbytes(jail.Path),
		Hostname: rawbytes(jail.Hostname),
		Jailname: rawbytes(jail.Jailname),
		IP4s: nip4s,
		IP6s: nip6s,
		IP4: rawip4s,
		IP6: rawip6s,
	}
	return syscall.Jail(&j)
}
