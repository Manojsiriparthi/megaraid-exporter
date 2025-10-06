# Contributing to MegaRAID Exporter

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Development Setup

1. **Prerequisites**
   - Go 1.19 or later
   - Git
   - MegaRAID controller with MegaCLI64 installed
   - Root privileges for testing

2. **Install MegaCLI64**
   ```bash
   # Download from LSI/Broadcom website
   # Or install via package manager
   sudo apt-get install megacli  # Debian/Ubuntu
   # Verify installation
   sudo megacli64 -v
   ```

3. **Clone and Setup**
   ```bash
   git clone https://github.com/yourusername/megaraid-exporter.git
   cd megaraid-exporter
   go mod download
   ```

4. **Build and Test**
   ```bash
   make build
   make test
   ```

## Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused
- Handle MegaCLI command errors gracefully

## MegaCLI Integration Guidelines

### Command Execution
- Always use `-NoLog` flag to prevent log spam
- Set appropriate timeouts for commands
- Handle controller-specific command variations
- Parse output defensively (MegaCLI format can vary)

### Error Handling
- Distinguish between command errors and parsing errors
- Log MegaCLI command output on failures
- Provide meaningful error messages to users
- Implement retry logic for transient failures

### Performance Considerations
- Cache command results when appropriate
- Limit concurrent MegaCLI executions
- Use efficient parsing algorithms
- Monitor memory usage during metric collection

## Testing

### Unit Tests
```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/megacli
```

### Integration Tests
```bash
# Test with real hardware (requires MegaCLI)
sudo go test -tags=integration ./...
```

### Test Data
- Use mock MegaCLI output for unit tests
- Include various controller configurations
- Test error conditions and edge cases
- Validate metric output format

## Submitting Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write clean, documented code
   - Add tests for new functionality
   - Update documentation if needed
   - Test with real MegaRAID hardware when possible

3. **Test Your Changes**
   ```bash
   make test
   make lint
   sudo ./megaraid-exporter --log-level debug  # Manual test
   ```

4. **Commit and Push**
   ```bash
   git add .
   git commit -m "Add: brief description of changes"
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request**
   - Provide clear description of changes
   - Reference any related issues
   - Include test results if applicable
   - Ensure CI passes

## Reporting Issues

When reporting bugs, please include:
- Operating system and version
- Go version
- MegaRAID controller model and firmware version
- MegaCLI version (`megacli64 -v`)
- Steps to reproduce
- Error messages or logs
- Expected vs actual behavior
- Relevant MegaCLI command output

## Feature Requests

For new features:
- Describe the use case
- Explain why it would be valuable
- Provide examples if possible
- Consider backward compatibility
- Include relevant MegaCLI commands if applicable

## Common MegaCLI Commands for Development

```bash
# Controller information
sudo megacli64 -AdpAllInfo -aALL -NoLog

# List all controllers
sudo megacli64 -adpCount -NoLog

# Virtual drives
sudo megacli64 -LDInfo -Lall -aALL -NoLog

# Physical drives
sudo megacli64 -PDList -aALL -NoLog

# BBU information
sudo megacli64 -AdpBbuCmd -aALL -NoLog

# Event log
sudo megacli64 -AdpEventLog -GetEvents -f /tmp/events.log -aALL -NoLog

# Foreign configuration
sudo megacli64 -CfgForeign -Scan -aALL -NoLog
```

## Documentation

- Update README for user-facing changes
- Add inline code comments
- Update configuration examples
- Include metric descriptions
- Update troubleshooting guides

Thank you for contributing!
