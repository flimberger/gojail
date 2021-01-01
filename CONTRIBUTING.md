# Contributing to gojail

## Issues

If you report issues,
please make sure to include the versions of FreeBSD and Go you are using.

## Code and More

### What

The `gojail` package should be small and self-contained,
if the feature does not concern jails,
it should go into another package.
If it is not in [`libjail`](https://www.freebsd.org/cgi/man.cgi?query=jail_getname&sektion=3&manpath=FreeBSD+12.2-RELEASE+and+Ports),
it should probably not be in `gojail`.

Implementation details from the `syscall` package should not leak out into the main package.

That said,
bug fixes and doc improvements are always appreciated!

### How

To contribute to gojail,
you can submit a pull request on GitHub,
or simply send patches via email.

If possible,
split contributions into logical changes and submit them separately.
This makes errors easier to trace back later.

Please be sure that the code is properly formatted before submitting.
You can use `go fmt` to format it for you.

### Copyright Tracking

We use the (Developer's Certificate of Origin)[https://developercertificate.org/] to track the copyright of contributions.
Please make sure you can certify the text of the next section below.
If you can,
add a line like the following to the end of your commit message:

    Signed-off-by: Jane Doe <jane.doe@example.com>

This is done automatically if you use `git commit -s` to commit your changes.

#### Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
