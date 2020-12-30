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

package main

import (
	"fmt"
	"os"
	"strconv"

	"purplekraken.com/pkg/gojail"
	"purplekraken.com/pkg/gojail/syscall"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "updating: gojail name hostname securelevel ipaddr")
		os.Exit(2)
	}

	params := []gojail.JailParam{}
	name, err := gojail.NewStringParam("name", os.Args[1])
	if err != nil {
		panic(err)
	}
	params = append(params, name)

	hostname, err := gojail.NewStringParam("host.hostname", os.Args[2])
	if err != nil {
		doError(err.Error())
	}
	params = append(params, hostname)

	secureint, err := strconv.Atoi(os.Args[3])
	if err != nil || secureint < 0 || secureint > 3 {
		doError("Invalid securelevel provided, must be a number between 0 and 3")
	}

	securelevel, err := gojail.NewIntParam("securelevel", secureint)
	if err != nil {
		doError(err.Error())
	}
	params = append(params, securelevel)

	ip4, err := gojail.NewIPParam(os.Args[4])
	if err != nil {
		doError(err.Error())
	}
	params = append(params, ip4)

	jid, err := gojail.SetParams(params, syscall.JAIL_UPDATE)

	if err != nil {
		if je, ok := err.(*gojail.JailErr); ok {
			fmt.Fprintln(os.Stderr, "gojail: errmsg:", je)
		} else if sce, ok := err.(*os.SyscallError); ok {
			fmt.Fprintln(os.Stderr, "gojail: syscall:", sce)
		} else {
			fmt.Fprintln(os.Stderr, "gojail:", err)
		}
		os.Exit(1)
	}
	fmt.Printf("Updated Jail with ID: %d\n", jid)

}

func doError(msg string) {
	fmt.Fprintln(os.Stderr, "Error parsing the arguments: ", msg)
	os.Exit(1)
}
