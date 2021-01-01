# gojail

Implementation of the FreeBSD
[`jail(2)`](https://www.freebsd.org/cgi/man.cgi?query=jail&sektion=2&manpath=FreeBSD+12.2-RELEASE+and+Ports)
and [`jailparam(3)`](https://www.freebsd.org/cgi/man.cgi?query=jailparam&sektion=3&manpath=FreeBSD+12.2-RELEASE+and+Ports)
APIs for Go.

## Packages

The `gojail` package provides high-level access to the `jail(2)` API,
while `gojail/syscall` implements the low-level system call interface.
The latter should be treated as an implementation detail and not be used by regular consumers of the API.
