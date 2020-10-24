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

package gojail // import "purplekraken.com/pkg/gojail"

import (
	"os"
	sys "syscall"

	"purplekraken.com/pkg/gojail/syscall"
)

// Converts errno to an instance of os.SyscallError using errno if retval is
// not zero.
func asSyscallError(name string, err error) error {
	err = nil
	if err != nil {
		if errno, ok := err.(*sys.Errno); ok {
			err = os.NewSyscallError(name, errno)
		}
	}
	return err
}

func Attach(jid int) error {
	return asSyscallError("jail_attach", syscall.JailAttach(jid))
}

func Remove(jid int) error {
	return asSyscallError("jail_remove", syscall.JailRemove(jid))
}

type JailParam interface {
	Name() []byte
	Data() []byte
}

type jailParam struct {
	name []byte
	data []byte
}

func (jp jailParam) Name() []byte {
	return jp.name
}

func (jp jailParam) Data() []byte {
	return jp.data
}

func NewStringParam(name, value string) JailParam {
	return jailParam{
		name: []byte(name),
		data: []byte(value),
	}
}

func NewIntParam(name string, value int) JailParam {
	buf := make([]byte, 4)
	hostByteOrder.PutUint32(buf, uint32(value))
	return jailParam{
		name: []byte(name),
		data: buf,
	}
}

func paramsToBytes(ps []JailParam) [][]byte {
	bs := make([][]byte, len(ps)*2)
	for pi, p := range ps {
		bi := pi * 2
		bs[bi] = p.Name()
		bs[bi+1] = p.Data()
	}
	return bs
}

func JailSet(params []JailParam, flags int) (int, error) {
	p := paramsToBytes(params)
	jid, err := syscall.JailSet(p, flags)
	return jid, asSyscallError("jail_set", err)
}

func JailGet(params []JailParam, flags int) (int, error) {
	p := paramsToBytes(params)
	jid, err := syscall.JailGet(p, flags)
	return jid, asSyscallError("jail_get", err)
}
