# sneat-go

[![Go CI](https://github.com/sneat-co/sneat-go-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/sneat-co/sneat-go-backend/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sneat-co/sneat-go-backend)](https://goreportcard.com/report/github.com/sneat-co/sneat-go-backend)

Go lang backend for sneat apps:

- https://sneat.app/
- https://dailyscrums.app/ - free open source tool to run your stand-up meetings

<!-- dev-approach:v1 -->
## Our approach to development

We build with our own tooling:

- **[SpecScore](https://specscore.md)** — specify requirements as `SpecScore.md` artifacts
- **[SpecStudio](https://specscore.studio)** — author & manage specs across their lifecycle
- **[inGitDB](https://ingitdb.com)** — store structured data in Git where applicable
- **[DALgo](https://dalgo.io)** — data access layer for Go
- **[cover100.dev](https://cover100.dev)** — drive toward 100% test coverage
- **[DataTug](https://datatug.io)** — query & explore data
<!-- /dev-approach -->

## 3-d party dependencies

- AWS - to send emails

## How to run?

Read https://github.com/sneat-co/sneat-devenv

## Testing chatbots locally

There is a dedicated section regards how to test Telegram bots locally in [src/bots](src/sneatgaeapp/bots).
