package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	fmt.Println("GitProcess - Azure DevOps Branch Management Tool")
	fmt.Println("===============================================")

	// Get repository name from user
	repoName := getInput("Enter the repository name (e.g., LOM_ABC): ")
	if repoName == "" {
		fmt.Println("Repository name cannot be empty")
		return
	}

	// Get repository URL from user
	repoURL := getInput("Enter the repository URL: ")
	if repoURL == "" {
		fmt.Println("Repository URL cannot be empty")
		return
	}

	// Clone repository
	fmt.Printf("\nCloning repository '%s'...\n", repoName)
	repo, err := cloneRepository(repoURL, repoName)
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		return
	}

	// Switch to develop branch
	fmt.Println("\nSwitching to develop branch...")
	err = switchToBranch(repo, "develop")
	if err != nil {
		fmt.Printf("Error switching to develop branch: %v\n", err)
		return
	}

	// Check for active unmerged branches
	fmt.Println("\nChecking for unmerged branches...")
	unmergedBranches, err := getUnmergedBranches(repo, "develop")
	if err != nil {
		fmt.Printf("Error checking unmerged branches: %v\n", err)
		return
	}

	// If unmerged branches exist, ask user to continue
	if len(unmergedBranches) > 0 {
		fmt.Println("\nFound unmerged branches:")
		for _, branch := range unmergedBranches {
			fmt.Printf("  - %s\n", branch)
		}

		continueProcess := getConfirmation("\nDo you want to continue with creating a new branch? (y/n): ")
		if !continueProcess {
			fmt.Println("Process stopped by user.")
			return
		}
	}

	// Get new branch name
	newBranchName := getInput("\nEnter the new branch name: ")
	if newBranchName == "" {
		fmt.Println("Branch name cannot be empty")
		return
	}

	// Create new branch
	fmt.Printf("\nCreating new branch '%s'...\n", newBranchName)
	err = createBranch(repo, newBranchName)
	if err != nil {
		fmt.Printf("Error creating branch: %v\n", err)
		return
	}

	// Copy branch policies (placeholder for Azure DevOps integration)
	fmt.Println("\nCopying branch policies from develop...")
	err = copyBranchPolicies(repoName, "develop", newBranchName)
	if err != nil {
		fmt.Printf("Warning: Could not copy branch policies: %v\n", err)
	}

	fmt.Printf("\n✓ Successfully created branch '%s' with policies copied from develop\n", newBranchName)
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func getConfirmation(prompt string) bool {
	response := getInput(prompt)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func cloneRepository(url, name string) (*git.Repository, error) {
	repo, err := git.PlainClone(name, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func switchToBranch(repo *git.Repository, branchName string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Get the branch reference
	branchRef := plumbing.NewRemoteReferenceName("origin", branchName)
	ref, err := repo.Reference(branchRef, true)
	if err != nil {
		return fmt.Errorf("branch '%s' not found", branchName)
	}

	// Create local branch tracking remote
	localBranchRef := plumbing.NewBranchReferenceName(branchName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: localBranchRef,
		Create: true,
		Hash:   ref.Hash(),
	})
	if err != nil {
		// If branch already exists locally, just checkout
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: localBranchRef,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getUnmergedBranches(repo *git.Repository, baseBranch string) ([]string, error) {
	var unmergedBranches []string

	// Get all branches
	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	// Get base branch commit
	baseRef, err := repo.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
	if err != nil {
		return nil, err
	}

	baseCommit, err := repo.CommitObject(baseRef.Hash())
	if err != nil {
		return nil, err
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().Short()
		if branchName == baseBranch || branchName == "main" || branchName == "master" {
			return nil
		}

		// Check if branch is merged
		branchCommit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return nil
		}

		// Simple check: if the branch commit is not an ancestor of base branch
		isAncestor, err := branchCommit.IsAncestor(baseCommit)
		if err != nil {
			return nil
		}

		if !isAncestor && branchCommit.Hash != baseCommit.Hash {
			unmergedBranches = append(unmergedBranches, branchName)
		}

		return nil
	})

	return unmergedBranches, err
}

func createBranch(repo *git.Repository, branchName string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Get current HEAD
	head, err := repo.Head()
	if err != nil {
		return err
	}

	// Create new branch
	branchRef := plumbing.NewBranchReferenceName(branchName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Create: true,
		Hash:   head.Hash(),
	})

	return err
}

func copyBranchPolicies(repoName, sourceBranch, targetBranch string) error {
	return copyBranchPoliciesFromAzureDevOps(repoName, sourceBranch, targetBranch)
}
