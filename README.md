# Single Page Application Server

## Installation

- Golang is required, minimum version: `1.23`

```sh
go install github.com/banan-tech/spaserver@latest
```

## Running

```sh
$GOHOME/bin/spaserver <ROOT-DIR>
```
Default port is `3000`, to specify another port:

```sh
$GOHOME/bin/spaserver -port 5000 <ROOT-DIR>
```
