GOTOOL :=	go

all:	build
.PHONY:	all

build:
	${GOTOOL} build -v ./...
.PHONY:	build

test:
	${GOTOOL} test ./...
.PHONY:	test

clean:
	${GOTOOL} clean -x ./...
.PHONY: clean
