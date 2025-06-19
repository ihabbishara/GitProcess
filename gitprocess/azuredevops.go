package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AzureDevOpsConfig struct {
	Organization string
	Project      string
	Repository   string
	PAT          string // Personal Access Token
}

type BranchPolicy struct {
	ID          int                    `json:"id,omitempty"`
	Type        PolicyType             `json:"type"`
	Settings    map[string]interface{} `json:"settings"`
	IsEnabled   bool                   `json:"isEnabled"`
	IsBlocking  bool                   `json:"isBlocking"`
	Scope       []Scope                `json:"scope"`
}

type PolicyType struct {
	ID   string `json:"id"`
	Name string `json:"displayName"`
}

type Scope struct {
	RefName        string `json:"refName"`
	MatchKind      string `json:"matchKind"`
	RepositoryID   string `json:"repositoryId"`
}

func getAzureDevOpsConfig() (*AzureDevOpsConfig, error) {
	config := &AzureDevOpsConfig{
		Organization: os.Getenv("AZDO_ORG"),
		Project:      os.Getenv("AZDO_PROJECT"),
		PAT:          os.Getenv("AZDO_PAT"),
	}

	if config.Organization == "" {
		config.Organization = getInput("Enter Azure DevOps Organization: ")
	}
	
	if config.Project == "" {
		config.Project = getInput("Enter Azure DevOps Project: ")
	}
	
	if config.PAT == "" {
		config.PAT = getInput("Enter Personal Access Token (PAT): ")
	}

	if config.Organization == "" || config.Project == "" || config.PAT == "" {
		return nil, fmt.Errorf("missing required Azure DevOps configuration")
	}

	return config, nil
}

func getBranchPolicies(config *AzureDevOpsConfig, repositoryID, branchName string) ([]BranchPolicy, error) {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/policy/configurations?repositoryId=%s&refName=refs/heads/%s&api-version=7.1-preview.1",
		config.Organization, config.Project, repositoryID, branchName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("", config.PAT)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get branch policies: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Value []BranchPolicy `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

func createBranchPolicy(config *AzureDevOpsConfig, policy BranchPolicy) error {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/policy/configurations?api-version=7.1-preview.1",
		config.Organization, config.Project)

	// Remove ID for creation
	policy.ID = 0

	jsonData, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth("", config.PAT)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create branch policy: %s - %s", resp.Status, string(body))
	}

	return nil
}

func getRepositoryID(config *AzureDevOpsConfig, repositoryName string) (string, error) {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s?api-version=7.1-preview.1",
		config.Organization, config.Project, repositoryName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth("", config.PAT)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get repository: %s - %s", resp.Status, string(body))
	}

	var repo struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return "", err
	}

	return repo.ID, nil
}

func copyBranchPoliciesFromAzureDevOps(repoName, sourceBranch, targetBranch string) error {
	config, err := getAzureDevOpsConfig()
	if err != nil {
		return err
	}

	// Get repository ID
	repositoryID, err := getRepositoryID(config, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repository ID: %v", err)
	}

	// Get source branch policies
	sourcePolicies, err := getBranchPolicies(config, repositoryID, sourceBranch)
	if err != nil {
		return fmt.Errorf("failed to get source branch policies: %v", err)
	}

	if len(sourcePolicies) == 0 {
		fmt.Printf("No policies found on branch '%s'\n", sourceBranch)
		return nil
	}

	fmt.Printf("Found %d policies on branch '%s'\n", len(sourcePolicies), sourceBranch)

	// Copy each policy to target branch
	for _, policy := range sourcePolicies {
		// Update scope to target branch
		for i := range policy.Scope {
			policy.Scope[i].RefName = fmt.Sprintf("refs/heads/%s", targetBranch)
		}

		err := createBranchPolicy(config, policy)
		if err != nil {
			fmt.Printf("Warning: Failed to copy policy '%s': %v\n", policy.Type.Name, err)
			continue
		}

		fmt.Printf("✓ Copied policy: %s\n", policy.Type.Name)
	}

	return nil
}