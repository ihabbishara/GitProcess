# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GitProcess is a Go application that automates Azure DevOps branch management workflows. It provides an interactive CLI tool for:

- Cloning repositories locally
- Switching to develop branch
- Detecting unmerged branches
- Creating new branches with user confirmation
- Copying branch policies from develop to new branches

## Development Commands

### Build and Run
```bash
# Build the application
go build -o gitprocess

# Run the application
./gitprocess

# Run directly with Go
go run .
```

### Dependencies
```bash
# Install/update dependencies
go mod tidy

# Add new dependencies
go get <package>
```

## Architecture

### Core Components

**main.go**: Entry point and main workflow orchestration
- Interactive prompts for user input
- Git operations using go-git library
- Workflow coordination

**azuredevops.go**: Azure DevOps API integration
- Branch policy management
- REST API client for Azure DevOps
- PAT (Personal Access Token) authentication

### Key Functions

**Git Operations** (main.go):
- `cloneRepository()`: Clone repo using go-git
- `switchToBranch()`: Switch to develop branch
- `getUnmergedBranches()`: Detect branches not merged to develop
- `createBranch()`: Create new branch from current HEAD

**Azure DevOps Integration** (azuredevops.go):
- `getBranchPolicies()`: Fetch policies from source branch
- `createBranchPolicy()`: Apply policies to target branch
- `getRepositoryID()`: Get repo ID for API calls

## Configuration

### Environment Variables
- `AZDO_ORG`: Azure DevOps organization
- `AZDO_PROJECT`: Azure DevOps project name  
- `AZDO_PAT`: Personal Access Token for API access

### Interactive Prompts
If environment variables are not set, the application will prompt for:
- Repository name and URL
- Azure DevOps organization and project
- Personal Access Token

## Dependencies

- `github.com/go-git/go-git/v5`: Pure Go Git implementation
- Standard library: `net/http`, `encoding/json`, `bufio`, `os`

## Error Handling

The application includes error handling for:
- Git operations (clone, branch switching, branch creation)
- Azure DevOps API calls (authentication, policy retrieval/creation)
- User input validation