# olog - OpenTelemetry Logging Facade

[![Go Reference](https://pkg.go.dev/badge/github.com/pellared/olog.svg)](https://pkg.go.dev/github.com/pellared/olog)
[![go.mod](https://img.shields.io/github/go-mod/go-version/pellared/olog)](go.mod)
[![LICENSE](https://img.shields.io/github/license/pellared/olog)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pellared/olog)](https://goreportcard.com/report/github.com/pellared/olog)
[![Codecov](https://codecov.io/gh/pellared/olog/branch/main/graph/badge.svg)](https://codecov.io/gh/pellared/olog)

⭐ `Star` this repository if you find it valuable and worth maintaining.

👁 `Watch` this repository to get notified about new releases, issues, etc.

## Description

The `olog` package provides an ergonomic frontend API for OpenTelemetry structured logging.

It is designed to provide a more user-friendly interface while using the OpenTelemetry Logs API as the backend.

### Features

1. **Simple API**: Easy-to-use methods like `Debug()`, `Info()`, `Warn()`, `Error()` similar to popular logging libraries
2. **Level-specific enabled checks**: Built-in support for common log levels with `DebugEnabled()`, `InfoEnabled()`, etc. for performance
3. **Structured logging**: Support for key-value pairs using the alternating syntax (similar to `slog`)
4. **Event logging**: Dedicated `Event()` and `EventAttr()` methods for semantic events
5. **Context support**: All methods accept `context.Context` for trace correlation
6. **Logger composition**: `With()` and `WithAttr()` methods for attribute composition
7. **Performance oriented**: Level-specific enabled checks to avoid expensive operations
8. **Automatic package detection**: Auto-detects caller's package name when logger name is not specified
9. **Type safety**: Support for both argument-based and strongly-typed attribute APIs

## Background

This module directly addresses the concerns raised in
[opentelemetry-specification#4661](https://github.com/open-telemetry/opentelemetry-specification/issues/4661):

1. **User-friendly API**: Provides simple methods like `Info()`, `Debug()` instead of complex record construction
2. **Performance**: Includes level-specific enabled methods (`DebugEnabled()`, `InfoEnabled()`, etc.) that are ergonomic and efficient
3. **Level support**: Built-in support for common logging levels with both generic and specific enabled checks
4. **Structured logging focus**: Emphasizes key-value pairs for structured data with both argument-based and typed APIs
5. **Ergonomic design**: Follows familiar patterns from `slog` and other popular logging libraries
6. **Zero-configuration**: Automatic package name detection means less boilerplate setup code
7. **Type safety**: Optional strongly-typed attribute API for better performance and compile-time safety

## Contributing

Feel free to create an issue,
join the [discussions](https://github.com/pellared/olog/discussions/2),
or propose a pull request.

Please follow the [Code of Conduct](CODE_OF_CONDUCT.md).

This module follows several key design principles:

1. **Ergonomic API**: Provides simple methods that are easy to use and understand
2. **Performance First**: Includes `Enabled()` checks and optimizations to minimize overhead
3. **Structured Logging**: Emphasizes key-value pairs over string formatting
4. **Compatibility**: Uses OpenTelemetry Logs API as the backend for full compatibility
5. **Composability**: Supports logger composition through `With()`
6. **Familiar Patterns**: Similar to `slog` design patterns that Go developers already know
