# Pre-commit Setup Guide

## Overview

This guide will help you set up `pre-commit` hooks for your project. Pre-commit hooks are useful for automatically running checks before committing code to ensure code quality, security,  and consistency.
TechOps has enabled general golang linters and gitleaks which should be enabled on each commit.

## Prerequisites

- Python 3.6 or higher
- `pip` (Python package installer)
- `git` installed and configured

## Installation

To install `pre-commit`, follow these steps:

1. **Install pre-commit**
   You can install `pre-commit` using `pip`:

```
   pip install pre-commit
```

## Sample Usage
This is a sample run on a basic commit:
```
$ git commit -m "add precommit"
Check Yaml...........................................(no files to check)Skipped
Fix End of Files.........................................................Passed
Trim Trailing Whitespace.................................................Passed
Check for added large files..............................................Passed
go fmt...............................................(no files to check)Skipped
go imports...........................................(no files to check)Skipped
golangci-lint........................................(no files to check)Skipped
Detect hardcoded secrets.................................................Passed
```

You can also run `pre-commit` test on all files:
```
pre-commit run --all-files
```
