# Stockyard Cutoff

**Link shortener and redirect manager — your domain, click tracking, QR codes, vanity URLs**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9330:9330 -v cutoff_data:/data ghcr.io/stockyard-dev/stockyard-cutoff
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9330` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9330` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `CUTOFF_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 10 links, 1k clicks tracked | Unlimited links and clicks |
| Price | Free | $1.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Creator & Small Business

## License

Apache 2.0
