# go-ping-imgur

A tiny Go CLI that pings Imgur images with simple HTTP `HEAD` requests.

No Imgur API key is required. This project is intentionally minimal — useful for learning Go, trying out [Cobra](https://github.com/spf13/cobra), and checking whether your hosted images still respond.

## Why not the Imgur API?

Imgur's developer registration has been closed for new applications for some time. Direct HTTP requests to image URLs still work fine and are enough for a simple "is this image still alive?" check.

## Install

```bash
go install ./cmd/ping-imgur
```

Or build locally:

```bash
go build -o bin/ping-imgur ./cmd/ping-imgur
```

## Usage

Ping a bare image ID:

```bash
ping-imgur ping cvWgXFc
```

Ping a direct CDN URL:

```bash
ping-imgur ping https://i.imgur.com/cvWgXFc.jpg
```

Ping many targets from a file:

```bash
ping-imgur ping -f images.example.txt
```

Set a request timeout:

```bash
ping-imgur ping -t 5s cvWgXFc.jpg
```

### Input formats

- Full Imgur URL: `https://i.imgur.com/abc123.jpg`, `https://imgur.com/abc123`
- Image ID with extension: `abc123.jpg` (uses the CDN)
- Bare image ID: `abc123` (uses the gallery page)

You can combine positional targets with `-f` (file targets are appended after CLI args).

### What “OK” means

- **CDN URL** (`abc123.jpg` or `i.imgur.com/...`): the image endpoint returned a successful HTTP status.
- **Bare ID** (gallery page): the gallery HTML responded — not a guarantee the image file still exists.
- Redirects to non-Imgur hosts are blocked. Imgur may return **429** if you ping too quickly from one IP.

## Example output

```text
[OK] cvWgXFc.jpg -> https://i.imgur.com/cvWgXFc.jpg (200) in 180ms
[OK] https://i.imgur.com/cvWgXFc.jpg (200) in 95ms

2/2 targets responded OK
```

## Project layout

```text
cmd/ping-imgur/     CLI entrypoint (Cobra commands)
internal/imgur/     URL normalization and HTTP ping logic
images.example.txt  Sample input file
```

## A note on image retention

Periodic HTTP requests may count as views, but Imgur does not guarantee indefinite hosting — especially for anonymous or inactive content. Treat this as a fun learning tool, not a backup strategy.

## License

MIT
