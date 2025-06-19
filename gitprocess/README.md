# GitProcess

A Go CLI application for automated Azure DevOps branch management workflows.

## Features

- **Repository Cloning**: Clone Azure DevOps repositories locally
- **Branch Management**: Switch to develop branch and create new branches
- **Unmerged Branch Detection**: Identify branches not merged to develop
- **Interactive Prompts**: User-friendly confirmation dialogs
- **Policy Management**: Copy branch policies from develop to new branches
- **Azure DevOps Integration**: Full REST API integration for policy management

## Prerequisites

- Go 1.19 or higher
- Azure DevOps Personal Access Token (PAT) with appropriate permissions:
  - Code (read)
  - Project and team (read)
  - Build (read)
  - Release (read)

## Installation

```bash
git clone <repository-url>
cd gitprocess
go build -o gitprocess
```

## Usage

Run the application:
```bash
./gitprocess
```

The application will prompt you for:
1. Repository name (e.g., "LOM_ABC")
2. Repository URL
3. Azure DevOps organization (or set `AZDO_ORG` environment variable)
4. Azure DevOps project (or set `AZDO_PROJECT` environment variable)
5. Personal Access Token (or set `AZDO_PAT` environment variable)

## Environment Variables

Set these to avoid interactive prompts:

```bash
export AZDO_ORG="your-organization"
export AZDO_PROJECT="your-project"
export AZDO_PAT="your-personal-access-token"
```

## Workflow

1. **Clone Repository**: Downloads repository to local directory
2. **Switch to Develop**: Checks out the develop branch
3. **Check Unmerged Branches**: Identifies branches not merged to develop
4. **User Confirmation**: Prompts to continue if unmerged branches exist
5. **Create New Branch**: Creates branch with user-provided name
6. **Copy Policies**: Copies all branch policies from develop to new branch

## Azure DevOps Integration

The application integrates with Azure DevOps REST API to:
- Retrieve repository information
- Fetch branch policies from the develop branch
- Apply the same policies to newly created branches

Supported policy types include:
- Minimum number of reviewers
- Required reviewers
- Build validation
- Status checks
- Comment requirements

## Error Handling

The application provides clear error messages for:
- Authentication failures
- Network connectivity issues
- Missing repositories or branches
- API permission errors
- Git operation failures

## Development

See [CLAUDE.md](CLAUDE.md) for development guidance and architecture details.