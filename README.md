# json-formatter

A JSON formatting, minification, validation and query tool written in Go.

## Build

```bash
go build -o json-formatter main.go
```

## Usage

```bash
./json-formatter fmt data.json
./json-formatter fmt data.json --indent 4
./json-formatter min data.json
./json-formatter validate data.json
./json-formatter keys data.json
./json-formatter get data.json .user.email
./json-formatter get data.json .items.0.name

# Stdin
cat data.json | ./json-formatter fmt
cat data.json | ./json-formatter get .name
```

## Made in [discord.gg/codeshare](https://discord.gg/codeshare) · [astra-dev.com.tr](https://astra-dev.com.tr)

## License

MIT
