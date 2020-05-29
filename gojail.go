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
	"unsafe"

	"purplekraken.com/pkg/gojail/syscall"
)

func Attach(jid int) error {
	return syscall.JailAttach(jid)
}

func Remove(jid int) error {
	return syscall.JailRemove(jid)
}

type JailParam struct {
	Name string
	Data []byte
}

func rawbytes(b []byte) *byte {
	return (*byte)(unsafe.Pointer(&b[0]));
}

func rawstring(s string) *byte {
	return (*byte)(unsafe.Pointer(&[]byte(s)[0]));
}

func Create(params []JailParam, flags int) error {
	iovs []syscall.Iovec
	for _, param := range params {
		iov := syscall.Iovec{
			Base: rawstring(param.Name),
			Len: uint64(len(param.Name)),
		};
		iovs = append(iovs, iov)
		iov = syscall.Iovec{
			Base: rawbytes(param.Data),
			Len: uint64(len(param.Data)),
		};
	}
	return syscall.JailSet((*syscall.Iovec)(unsafe.Pointer(&iovs[0])), len(iovs), flags)
}
