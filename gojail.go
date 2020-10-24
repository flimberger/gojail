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
	"strconv"
	sys "syscall"

	"purplekraken.com/pkg/gojail/syscall"
)

const (
	errmsglen  = 1024
	maxnamelen = 256 // MAXHOSTNAMELEN on FreeBSD, defined in include/sys/param.h
)

type ParamType int

const (
	String ParamType = 0
	Int    ParamType = 1
	Raw    ParamType = 2
)

type JailParam interface {
	Name() []byte
	Data() []byte
	Type() ParamType
}

type jailParam struct {
	name  []byte
	data  []byte
	ptype ParamType
}

func (jp jailParam) Name() []byte {
	return jp.name
}

func (jp jailParam) Data() []byte {
	return jp.data
}

func (jp jailParam) Type() ParamType {
	return jp.ptype
}

func NewStringParam(name, value string) JailParam {
	return jailParam{
		name:  []byte(name),
		data:  []byte(value),
		ptype: String,
	}
}

func NewIntParam(name string, value int) JailParam {
	buf := make([]byte, 4)
	hostByteOrder.PutUint32(buf, uint32(value))
	return jailParam{
		name:  []byte(name),
		data:  buf,
		ptype: Int,
	}
}

// Error message from the jail subsystem.
// Represents an error returned as the "errmsg" parameter from JailGet or JailSet.
type JailErr struct {
	errmsg string
}

func (je *JailErr) Error() string {
	return je.errmsg
}

func makeJailErr(errmsg []byte) error {
	return &JailErr{
		errmsg: string(errmsg),
	}
}

func intToBytes(i int) []byte {
	b := make([]byte, 4)
	hostByteOrder.PutUint32(b, uint32(i))
	return b
}

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

// Returns the JID of the jail identified by name.
func GetId(name string) (int, error) {
	var iov [4][]byte
	if jid, err := strconv.Atoi(name); err == nil {
		if jid == 0 {
			return jid, nil
		}
		iov[0] = []byte("jid")
		iov[1] = intToBytes(jid)
	} else {
		iov[0] = []byte("name")
		iov[1] = []byte(name)
	}
	iov[2] = []byte("errmsg")
	iov[3] = make([]byte, errmsglen)
	jid, err := syscall.JailGet(iov[:], 0)
	if err != nil {
		if iov[3][0] != 0 {
			err = makeJailErr(iov[3])
		} else {
			err = asSyscallError("jail_get", err)
		}
		return -1, err
	}
	return jid, nil
}

// Returns the name of the jail identified by jid.
func GetName(jid int) (string, error) {
	var iov [6][]byte
	iov[0] = []byte("jid")
	iov[1] = intToBytes(jid)
	iov[2] = []byte("name")
	iov[3] = make([]byte, maxnamelen)
	iov[4] = []byte("errmsg")
	iov[5] = make([]byte, errmsglen)
	jid, err := syscall.JailGet(iov[:], 0)
	if err != nil {
		if iov[3][0] != 00 {
			err = makeJailErr(iov[5])
		} else {
			err = asSyscallError("jail_get", err)
		}
		return "", err
	}
	return string(iov[3]), err
}

// Attach the current process to the jail identified by jid.
// See jail_attach(2) for further information.
func Attach(jid int) error {
	return asSyscallError("jail_attach", syscall.JailAttach(jid))
}

// Remove the jail idenified by jid.
// See jail_remove(2) for further information.
func Remove(jid int) error {
	return asSyscallError("jail_remove", syscall.JailRemove(jid))
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

func SetParams(params []JailParam, flags int) (int, error) {
	p := paramsToBytes(params)
	jid, err := syscall.JailSet(p, flags)
	return jid, asSyscallError("jail_set", err)
}

func GetParams(params []JailParam, flags int) (int, error) {
	p := paramsToBytes(params)
	jid, err := syscall.JailGet(p, flags)
	return jid, asSyscallError("jail_get", err)
}
