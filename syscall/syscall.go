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

// This file contains code from the Go programming language, available under the following license:
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Low-level implementation of jail-related syscalls.
package syscall // import "purplekraken.com/pkg/gojail/syscall"

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	JAIL_CREATE = 0x01 // Create jail if it does not exist
	JAIL_UPDATE = 0x02 // Update parameters of existing jail
	JAIL_ATTACH = 0x04 // Attach to jail upon creation
	JAIL_DYING  = 0x08 // Allow getting a dying jail
)

// Do the interface allocations only once for common
// Errno values.
var (
	errEAGAIN       error = syscall.EAGAIN
	errEEXIST       error = syscall.EEXIST
	errEFAULT       error = syscall.EFAULT
	errEINVAL       error = syscall.EINVAL
	errENAMETOOLONG error = syscall.ENAMETOOLONG
	errENOENT       error = syscall.ENOENT
	errEPERM        error = syscall.EPERM
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.EAGAIN:
		return errEAGAIN
	case syscall.EEXIST:
		return errEEXIST
	case syscall.EFAULT:
		return errEFAULT
	case syscall.EINVAL:
		return errEINVAL
	case syscall.ENAMETOOLONG:
		return errENAMETOOLONG
	case syscall.ENOENT:
		return errENOENT
	case syscall.EPERM:
		return errEPERM
	}
	return e
}

func syscall1(sysnum uintptr, jid int) error {
	_, _, e := unix.Syscall(sysnum, uintptr(jid), 0, 0)
	return errnoErr(e)
}

func JailAttach(jid int) error {
	return syscall1(unix.SYS_JAIL_ATTACH, jid)
}

func JailRemove(jid int) error {
	return syscall1(unix.SYS_JAIL_REMOVE, jid)
}

var _zero uintptr

func bytes2iovec(bs [][]byte) []syscall.Iovec {
	iovecs := make([]syscall.Iovec, len(bs))
	for i, b := range bs {
		iovecs[i].SetLen(len(b))
		if len(b) > 0 {
			iovecs[i].Base = &b[0]
		} else {
			iovecs[i].Base = (*byte)(unsafe.Pointer(&_zero))
		}
	}
	return iovecs
}

func syscall2(sysnum uintptr, params [][]byte, flags int) (int, error) {
	iovs := bytes2iovec(params)
	var p unsafe.Pointer
	if len(iovs) > 0 {
		p = unsafe.Pointer(&iovs[0])
	} else {
		p = unsafe.Pointer(&_zero)
	}
	jid, _, e := unix.Syscall(sysnum, uintptr(p), uintptr(len(iovs)), uintptr(flags))
	return int(jid), errnoErr(e)
}

func JailGet(params [][]byte, flags int) (int, error) {
	return syscall2(unix.SYS_JAIL_GET, params, flags)
}

func JailSet(params [][]byte, flags int) (int, error) {
	return syscall2(unix.SYS_JAIL_SET, params, flags)
}
