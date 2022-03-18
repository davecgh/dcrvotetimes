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

- **Go 1.17 or 1.18**

  Installation instructions can be found here: https://go.dev/doc/install.
  Ensure Go was installed properly and is a supported version:
  ```sh
  $ go version
  $ go env GOROOT GOPATH
  ```
  NOTE: `GOROOT` and `GOPATH` must not be on the same path. Since Go 1.8 (2016),
  `GOROOT` and `GOPATH` are set automatically, and you do not need to change
  them. However, you still need to add `$GOPATH/bin` to your `PATH` in order to
  run binaries installed by `go get` and `go install` (On Windows, this happens
  automatically).

  Unix example -- add these lines to .profile:

  ```
  PATH="$PATH:/usr/local/go/bin"  # main Go binaries ($GOROOT/bin)
  PATH="$PATH:$HOME/go/bin"       # installed Go projects ($GOPATH/bin)
  ```

**Building/Installing/Updating**:

Run the follow command to build the latest release version of the `dcrvotetimes`
executable from source and install it to `$GOPATH/bin/`:

```
go install github.com/davecgh/dcrvotetimes@latest
```

## Issue Tracker

The [integrated github issue tracker](https://github.com/davecgh/dcrvotetimes/issues)
is used for this project.

## License

dcrvotetimes is licensed under the [copyfree](http://copyfree.org) ISC License.
