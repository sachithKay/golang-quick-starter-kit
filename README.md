# Ghost - Go Quickstart Project

This is a Golang quickstart project implementing production-grade architectural patterns, setting up robust CI/CD pipelines with Docker, and a unified Makefile.

## Prerequisites

- **Go 1.20+**
- **Docker** (for containerized environments)
- **buf** CLI (to generate Protocol Buffer files if used)

## Getting Started

To run the application locally, you can use the provided Makefile commands:

### Run Locally

```bash
make run
```

### Run Tests

Execute all tests in the project:

```bash
make test
```

### Build Production Binary

Build the executable binary into the `bin/` directory:

```bash
make build
```

### Generate Protobuf Files

To compile and generate protobuf files via `buf`:

```bash
make generate
```

### Clean up Binaries

Clean the generated binaries:

```bash
make clean
```
