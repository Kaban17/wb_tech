# L2-8: Go Time Application

This is a simple Go application that fetches and displays the current time from an NTP server.


### Building

To build the application, run the following command:

```sh
make build
```

This will create an executable file named `myapp` in the project directory.

### Running

To run the application, use the following command:

```sh
make run
```

## Development

### Static Analysis

To run the static analysis tools (`go vet` and `golint`), you can use the following targets:

```sh
# Run go vet
make vet

# Run golint
make lint
```

If you don't have `golint` installed, you can install it by running:

```sh
make install-lint
```

The default `make` command will run both `vet` and `lint`:

```sh
make
```
