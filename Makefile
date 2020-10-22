PKG :=	purplekraken.com/pkg/gojail

GOTOOL :=	go

all:	build
.PHONY:	all

build:
	${GOTOOL} build cmd
.PHONY:	build

test:
	${GOTOOL} test ${PKG}/syscall
.PHONY:	test
