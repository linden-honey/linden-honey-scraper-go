# Linden Honey Scraper

> Lyrics scraper service powered by Golang

[![build](https://img.shields.io/github/workflow/status/linden-honey/linden-honey-scraper-go/CI)](https://github.com/linden-honey/linden-honey-scraper-go/actions?query=workflow%3ACI)
[![version](https://img.shields.io/github/go-mod/go-version/linden-honey/linden-honey-scraper-go)](https://go.dev/)
[![coverage](https://img.shields.io/codecov/c/github/linden-honey/linden-honey-scraper-go)](https://codecov.io/github/linden-honey/linden-honey-scraper-go)
[![tag](https://img.shields.io/github/tag/linden-honey/linden-honey-scraper-go.svg)](https://github.com/linden-honey/linden-honey-scraper-go/tags)
[![reference](https://pkg.go.dev/badge/github.com/linden-honey/linden-honey-scraper-go.svg)](https://pkg.go.dev/github.com/linden-honey/linden-honey-scraper-go)

## Technologies

- [Golang](https://go.dev/)

## Usage

### Local

Build application artifacts:

```bash
make build
```

Run tests:

```bash
make test
```

Run application:

```bash
make run
```

### Docker

Bootstrap full project using docker compose:

```bash
docker compose up
```

Bootstrap project excluding some services using docker compose:

```bash
docker compose up --scale [SERVICE=0...]
```

Stop and remove containers, networks, images:

```bash
docker compose down
```

## Application instance

[https://linden-honey-scraper-go.herokuapp.com](https://linden-honey-scraper-go.herokuapp.com)
