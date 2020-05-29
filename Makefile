PKG := purplekraken.com/pkg/gojail

GOTOOL := go

all:	build
.PHONY:	all

build:
	${GOTOOL} build ${PKG}/syscall
.PHONY:	build

test:
	${GOTOOL} test ${PKG}/syscall
.PHONY:	test
