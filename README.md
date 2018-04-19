dcrvotetimes
============

dcrvotetimes provides a simple tool for Decred that calculates the mean time it
took all tickets on the blockchain the vote.  In order to provide the most
accurate value, it only considers the time between the moment each ticket was
added to live ticket pool and thus became eligible to vote and them time it
voted.

This requires a dcrd instance with the JSON-RPC server running.  Specify the
RPC username and password on the command line to this tool via the `-rpcuser`
and `-rpcpass` flags, respectively.

Additional details about every ticket can be displayed with the `-verbose` flag
if desired.

## Installation and updating

### Windows/Linux/BSD/POSIX - Build from source

Building or updating from source requires the following build dependencies:

- **Go 1.9 or 1.10**

  Installation instructions can be found here: http://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

- **Dep**

  Dep is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/golang/dep/cmd/dep`

Unfortunately, the use of `dep` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `dep`, and then `go` is used to build and
install the project.

**Getting the source**:

For a first time installation, the project and dependency sources can be
obtained manually with `git` and `glide` (create directories as needed):

```
git clone https://github.com/davecgh/dcrvotetimes $GOPATH/src/github.com/davecgh/dcrvotetimes
cd $GOPATH/src/github.com/davecgh/dcrvotetimes
dep ensure
```

To update an existing source tree, pull the latest changes and install the
matching dependencies:

```
cd $GOPATH/src/github.com/davecgh/dcrvotetimes
git pull
dep ensure
```

**Building/Installing**:

The `go` tool is used to build or install (to `GOPATH`) the project.  Some
example build instructions are provided below (all must run from the
`dcrvotetimes` project directory).

To build a `dcrvotetimes` executable and install it to `$GOPATH/bin/`:

```
go install
```

To build a `dcrvotetimes` executable and place it in the current directory:

```
go build
```

## Issue Tracker

The [integrated github issue tracker](https://github.com/davecgh/dcrvotetimes/issues)
is used for this project.

## License

dcrvotetimes is licensed under the [copyfree](http://copyfree.org) ISC License.
