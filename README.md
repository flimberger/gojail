# gojail

Implementation of the FreeBSD
[`jail(2)` API](https://www.freebsd.org/cgi/man.cgi?query=jail&apropos=0&sektion=2&manpath=FreeBSD+12.1-RELEASE&arch=default&format=html)
for Go.

# Packages

The `gojail` package provides high-level access to the `jail(2)` API,
while `gojail/syscall` implements the low-level system call interface.
