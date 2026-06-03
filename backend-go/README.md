# Project backend-go

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## OpenTelemetry (Honeycomb)

Set these environment variables before running the API to export traces:

- `OTEL_EXPORTER_OTLP_ENDPOINT` (example: `https://api.honeycomb.io`)
- `OTEL_EXPORTER_OTLP_HEADERS` (example: `x-honeycomb-team=YOUR_API_KEY`)
- `OTEL_SERVICE_NAME` (example: `backend-go`)
