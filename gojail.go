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
	"fmt"
	"net"
	"os"
	"strconv"
	sys "syscall"

	"golang.org/x/sys/unix"
	"purplekraken.com/pkg/gojail/syscall"
)

const (
	errmsglen  = 1024
	maxnamelen = 256 // MAXHOSTNAMELEN on FreeBSD, defined in include/sys/param.h
)

type Flags int

const (
	CreateFlag     Flags = syscall.JAIL_CREATE
	UpdateFlag     Flags = syscall.JAIL_UPDATE
	AttachFlag     Flags = syscall.JAIL_ATTACH
	AllowDyingFlag Flags = syscall.JAIL_DYING
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

func NewStringParam(name, value string) (JailParam, error) {
	nameb, err := unix.ByteSliceFromString(name)
	if err != nil {
		return nil, err
	}
	valueb, err := unix.ByteSliceFromString(value)
	if err != nil {
		return nil, err
	}
	return jailParam{
		name:  nameb,
		data:  valueb,
		ptype: String,
	}, nil
}

func NewIntParam(name string, value int) (JailParam, error) {
	nameb, err := unix.ByteSliceFromString(name)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 4)
	hostByteOrder.PutUint32(buf, uint32(value))
	return jailParam{
		name:  nameb,
		data:  buf,
		ptype: Int,
	}, nil
}

// TODO: The IP must be added to some interface
// otherwise is just an IP :)
// inet 192.168.0.222 netmask 0xffffffff broadcast 192.168.0.222
func NewIPParam(value string) (JailParam, error) {
	ip := net.ParseIP(value)
	if ip == nil {
		return nil, fmt.Errorf("Invalid IP address provided")
	}

	var nameb []byte
	var buf []byte
	if ip4 := ip.To4(); ip4 != nil {
		nameb = byteSliceFromStringOrDie("ip4.addr")
		buf = ip4
	} else {
		nameb = byteSliceFromStringOrDie("ip6.addr")
		buf = ip
	}

	return jailParam{
		name:  nameb,
		data:  buf,
		ptype: Raw,
	}, nil
}

// Error message from the jail subsystem.
// Represents an error returned as the "errmsg" parameter from JailGet or JailSet.
type JailErr struct {
	errmsg string
}

func (je *JailErr) Error() string {
	return je.errmsg
}

// Error returned by GetId and GetName if the specified jail does not exist.
var NoJail error = &JailErr{errmsg: "No such jail"}

func makeJailErr(errmsg []byte) error {
	return &JailErr{
		errmsg: string(errmsg),
	}
}

func byteSliceFromStringOrDie(s string) []byte {
	b, err := unix.ByteSliceFromString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func intToBytes(i int) []byte {
	b := make([]byte, 4)
	hostByteOrder.PutUint32(b, uint32(i))
	return b
}

// Converts errno to an instance of os.SyscallError using errno if retval is
// not zero.
func asSyscallError(name string, err error) error {
	if err != nil {
		if errno, ok := err.(sys.Errno); ok {
			err = os.NewSyscallError(name, errno)
		}
	}
	return err
}

// Returns the JID of the jail identified by name, -1 if it doesn't exist.
func GetId(name string) (int, error) {
	var iov [4][]byte
	if jid, err := strconv.Atoi(name); err == nil {
		if jid == 0 {
			return jid, nil
		}
		iov[0] = byteSliceFromStringOrDie("jid")
		iov[1] = intToBytes(jid)
	} else {
		iov[0] = byteSliceFromStringOrDie("name")
		iov[1], err = unix.ByteSliceFromString(name)
		if err != nil {
			return -1, err
		}
	}
	iov[2] = byteSliceFromStringOrDie("errmsg")
	iov[3] = make([]byte, errmsglen)
	jid, err := syscall.JailGet(iov[:], 0)
	if err != nil {
		// The jail does not exist, but that is not really an error.
		// Checking the kind of error is tedious for the users, so
		// differentiate here.
		// Attention: jail_get(2) returns ENOENT on three occasions:
		// 1. The jail referred to by a jid or name parameter does
		//    not exist.
		// 2. The jail referred to by a jid is not accessible by the
		//    process, because the process is in a different jail.
		// 3. The lastjid parameter is greater than the highest
		//    current jail ID.
		// We don't care for the second case, because the situation
		// is equivalent to the first case, for processes in a jail
		// other jails do not exist, but we need to be careful with
		// the third case.
		// In this function, there is no "lastjid" parameter, so
		// everything is fine, but this is not the general case.
		if syserr, ok := err.(sys.Errno); ok && syserr == sys.ENOENT {
			err = NoJail
		} else {
			err = asSyscallError("jail_get", err)
		}
	} else if jid == -1 && iov[3][0] != 0 {
		err = makeJailErr(iov[3])
	}
	return jid, err
}

// Returns the name of the jail identified by jid.
func GetName(jid int) (string, error) {
	var iov [6][]byte
	iov[0] = byteSliceFromStringOrDie("jid")
	iov[1] = intToBytes(jid)
	iov[2] = byteSliceFromStringOrDie("name")
	iov[3] = make([]byte, maxnamelen)
	iov[4] = byteSliceFromStringOrDie("errmsg")
	iov[5] = make([]byte, errmsglen)
	jid, err := syscall.JailGet(iov[:], 0)
	if err != nil {
		if syserr, ok := err.(sys.Errno); ok && syserr == sys.ENOENT {
			err = NoJail
		} else {
			err = asSyscallError("jail_get", err)
		}
	} else if jid == -1 && iov[5][0] != 00 {
		err = makeJailErr(iov[5])
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

func SetParams(params []JailParam, flags Flags) (int, error) {
	p := paramsToBytes(params)
	jid, err := syscall.JailSet(p, int(flags))
	return jid, asSyscallError("jail_set", err)
}

func GetParams(params []JailParam, flags Flags) (int, error) {
	p := paramsToBytes(params)
	jid, err := syscall.JailGet(p, int(flags))
	return jid, asSyscallError("jail_get", err)
}
