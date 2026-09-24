# CI Engine

A lightweight continuous integration engine with custom YAML syntax for building Go and Rust projects into binaries.

## Overview

This CI engine provides a minimal, opinionated build system that uses **topological sorting** to resolve job dependencies and execute them in the correct order. Build output is **streamed directly to stdout** for real-time visibility into the compilation process.

## Features

* **Custom YAML Pipeline Syntax** — Define build jobs with dependencies, language detection, and environment configuration
* **Topological Job Ordering** — Automatically resolves job dependencies and executes them in the optimal order
* **Real-Time Stream Output** — Build logs are streamed to stdout as they execute, eliminating wait time for build completion

## Supported Languages

* Go
* Rust

## Quick Start

Create a `.pipeline.yaml` file in your project root:

```yaml
pipeline:
  build:
    name: build
    language: go
    environment: dev

  lint:
    name: lint
    language: go
    environment: dev
    needs: [test]

  test:
    name: test
    language: go
    environment: dev
    needs: [build]

policy:
  network: false
```

yamlRun the pipeline:

```bash
ci-engine run .  (assuming binary name is ci-engine)

#Localhost
go run . (to run it locall after cloning)
```

bash

## Pipeline Configuration

### Pipeline Section

Define build jobs with the following fields:

| Field           | Type   | Description                                               |
| --------------- | ------ | --------------------------------------------------------- |
| `name`        | string | Unique identifier for the job                             |
| `language`    | string | Build language (go, rust)                                 |
| `environment` | string | Environment context (dev, prod, etc.)                     |
| `needs`       | array  | List of job names that must complete before this job runs |

### Policy Section

Control execution constraints:

| Field       | Type    | Description                                                    |
| ----------- | ------- | ------------------------------------------------------------   |
| `network`   | boolean | Enable/disable network access during builds (default: false)   |

## How It Works

1. **Parse YAML** — Configuration is loaded and validated
2. **Dependency Resolution** — Jobs are organized using topological sorting
3. **Execution** — Jobs execute in dependency order with real-time output streaming
4. **Binary Output** — Compiled binaries are written to the workspace directory

## Example Workflow

Given the configuration above, execution order is:

1. `build` (no dependencies)
2. `test` (waits for `build`)
3. `lint` (waits for `test`)

Each job's output streams to stdout as it completes.

Each Pipeline is different
