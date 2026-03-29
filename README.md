# go-mod-add

A CLI tool that forcefully updates Go module dependencies by removing them from `go.mod` first, then re-adding the specified versions via `go get -u`.

## Why?

`go get` may refuse to downgrade or change a dependency if constraint resolution prevents it. This tool works around that by dropping the `require` entry before re-fetching, ensuring the exact version you specify is installed.

## Installation

```bash
go install go-mod-add@latest
```

Or build from source:

```bash
go build -o go-mod-add .
```

## Usage

```bash
go-mod-add <module@version> [<module@version> ...]
```

### Examples

Update a single dependency:

```bash
go-mod-add golang.org/x/mod@v0.34.0
```

Update multiple dependencies at once:

```bash
go-mod-add golang.org/x/mod@v0.34.0 github.com/stretchr/testify@v1.9.0
```

## What it does

For each specified dependency the tool:

1. **Removes** the existing `require` entry from `go.mod` (using `golang.org/x/mod/modfile`).
2. **Runs** `go get -u <module>@<version>` to fetch the desired version.
3. **Runs** `go mod tidy` to clean up.
4. **Runs** `go mod vendor` if a `vendor/` directory exists.

## Development

This project uses [Mage](https://magefile.org/) as a build tool.

```bash
# Build the binary
mage build

# Format source code
mage format
```

## License

MIT
