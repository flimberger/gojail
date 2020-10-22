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

// A Low-level implementation of the `jail(2)` and related syscalls.
// Despite the name, they are not implemented as raw syscalls, but rather as
// wrappers for the libc.

package syscall

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

/*
#include <sys/param.h>
#include <sys/jail.h>
*/
import "C"

type Iovec unix.Iovec

const (
	JAIL_CREATE = C.JAIL_CREATE
	JAIL_UPDATE = C.JAIL_UPDATE
	JAIL_ATTACH = C.JAIL_ATTACH
	JAIL_DYING = C.JAIL_DYING
	JAIL_SET_MASK = C.JAIL_SET_MASK
	JAIL_GET_MASK = C.JAIL_GET_MASK
	JAIL_SYS_DISABLE = C.JAIL_SYS_DISABLE
	JAIL_SYS_NEW = C.JAIL_SYS_NEW
	JAIL_SYS_INHERIT = C.JAIL_SYS_INHERIT
)

// Converts errno to an instance of os.SyscallError using errno if retval is
// not zero.
func asError(name string, err error) error {
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			return os.NewSyscallError(name, syscall.Errno(errno))
		}
		return err
	}
	return nil
}

func JailAttach(jid int) error {
	_, err := C.jail_attach(C.int(jid))
	return asError("jail_attach", err)
}

func JailRemove(jid int) error {
	_, err := C.jail_remove(C.int(jid))
	return asError("jail_remove", err)
}

func JailGet(iov *Iovec, niov uint, flags int) error {
	unsafeptr := (*C.struct_iovec)(unsafe.Pointer(iov))
	_, err := C.jail_get(unsafeptr, C.uint(niov), C.int(flags))
	return asError("jail_get", err)
}

func JailSet(iov []Iovec, flags int) error {
	niov := uint(len(iov))
	unsafeptr := (*C.struct_iovec)(unsafe.Pointer(&iov[0]))
	_, err := C.jail_set(unsafeptr, C.uint(niov), C.int(flags))
	return asError("jail_set", err)
}
