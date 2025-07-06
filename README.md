# Gomer

Telegram bot for team management and reporting.

## Quick Start

### Using Docker Compose

1. Copy environment template:
```bash
cp env.example .env
```

2. Configure your environment variables in `.env` file

3. Run:
```bash
docker compose up -d
```

### Local Development

1. Copy environment template:
```bash
cp env.example .env
```

2. Configure your environment variables in `.env` file

3. Load environment and run:
```bash
source scripts/load-env.sh
go run ./cmd/gomer/main.go
```

## Configuration

See [CONFIGURATION.md](CONFIGURATION.md) for detailed configuration options.
