# Guidelines for AI coding agents

# Common sense

- Act as a experienced senior software engineer
- Adhere to DRY principal – reuse common code where reasonably possible

## Tests

### Mocking

- Use `mock_dal` package for mocking `dal` interfaces
- Use `go.uber.org/mock/gomock` to generate mocks and use them

## Pre-submit checks

- Run `golangci-lint run` - 0 issues expected
- Run `go test ./...` - 0 tests expected to fail
