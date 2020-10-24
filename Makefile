PKG :=	purplekraken.com/pkg/gojail
CMD :=	gojail

GOTOOL :=	go

all:	gojail
.PHONY:	all

gojail:
	${GOTOOL} build -o ${CMD} ${PKG}/cmd
.PHONY:	gojail

test:
	${GOTOOL} test ${PKG}/syscall
.PHONY:	test

clean:
	${GOTOOL} clean ${PKG}
	rm -f ${CMD}
.PHONY: clean
