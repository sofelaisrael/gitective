# Gitective 🔍

> Investigating your Git history...

A Go CLI tool for digging into your Git history.

## Install

```bash
go install github.com/sofelaisrael/gitective/cmd/gitective@latest
```

## Usage

```bash
go run ./cmd/gitective
# or after build
./gitective
```

Output:
```
Gitective
Investigating your Git history...
```

## Project Structure

```
.
├── cmd/gitective/main.go  # CLI entrypoint
├── go.mod
└── README.md
```

## Development

```bash
go run ./cmd/gitective
go build -o gitective ./cmd/gitective
```

## License

MIT
