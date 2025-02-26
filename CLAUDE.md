# CLAUDE.md - FanyiHub Development Guide

## Build Commands
- Build app: `wails build`
- Development mode: `wails dev`
- Run specific tests: `go test ./[package]` (e.g., `go test ./config`)
- Lint code: `golangci-lint run`

## Code Style Guidelines
- **Packages**: Organized by functionality (config, hotkey, langdetect, llm)
- **Imports**: Standard Go import grouping (std lib first, then external)
- **Error Handling**: Return errors with context using `fmt.Errorf()`
- **Naming**: 
  - Use CamelCase for exported functions/variables
  - Use lowerCamelCase for unexported functions/variables
  - Brief but descriptive package/type names
- **Comments**: Document all exported functions and types
- **Logging**: Use `slog` for structured logging
- **Types**: Prefer explicit types over interface{} or any
- **Config**: Follow existing patterns in config/config.go for app configuration

Remember to run `wails build` to rebuild the app after making changes.