// SPDX-License-Identifier: BSD-2-Clause-FreeBSD
//
// Copyright (c) 2020 Ricardo Pchevuzinske Katz <ricardo.katz@gmail.com>
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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: removing id")
		os.Exit(2)
	}
	jid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid ID provided: ", err.Error())
	}
	err = gojail.Remove(jid)

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
	fmt.Printf("Removed jail with ID: %d\n", jid)

}

func doError(msg string) {
	fmt.Fprintln(os.Stderr, "Error parsing the arguments: ", msg)
	os.Exit(1)
}
