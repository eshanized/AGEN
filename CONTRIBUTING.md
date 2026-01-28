# Contributing to AGEN

First off, thanks for taking the time to contribute! 🎉

The following is a set of guidelines for contributing to AGEN. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

- **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/eshanized/agen/issues).
- If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/eshanized/agen/issues/new). Be sure to include a **title and clear description**, as well as as much relevant information as possible.
- Include a **code sample** or an **executable test case** demonstrating the expected behavior that is not occurring.

### Suggesting Enhancements

- Open a new issue and describe the feature you would like to see.
- Explain *why* this feature would be useful to most AGEN users.

### Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes (`go test ./...`).
5. Make sure your code lints (`go vet ./...`).

## Styleguides

### Go Styleguide

- We follow the [official Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Use `gofmt` to format your code.

### Commit Messages

- Use the present tense ("Add feature" not "Added feature").
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...").
- Limit the first line to 72 characters or less.
- Reference issues and pull requests liberally after the first line.
