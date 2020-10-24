PKG :=	purplekraken.com/pkg/gojail

GOTOOL :=	go

all:	build
.PHONY:	all

build:
	#${GOTOOL} build cmd
	${GOTOOL} build
.PHONY:	build

test:
	${GOTOOL} test ${PKG}/syscall
.PHONY:	test

clean:
	${GOTOOL} clean ${PKG}
.PHONY: clean
