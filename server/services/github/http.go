package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPGitHubService implements GitHubService using HTTP API calls.
type HTTPGitHubService struct {
	httpClient *http.Client
	baseURL    string
}

// NewHTTPGitHubService creates a new HTTPGitHubService.
func NewHTTPGitHubService() Service {
	return &HTTPGitHubService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

// Key represents a GitHub SSH key.
type Key struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

// GetSSHKeys retrieves SSH public keys for the given GitHub usernames.
func (s *HTTPGitHubService) GetSSHKeys(ctx context.Context, githubIDs []string) ([]string, error) {
	var allKeys []string

	for _, githubID := range githubIDs {
		keys, err := s.getUserSSHKeys(ctx, githubID)
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
	}

	return allKeys, nil
}

// getUserSSHKeys retrieves SSH keys for a single GitHub user.
func (s *HTTPGitHubService) getUserSSHKeys(ctx context.Context, githubID string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/keys", s.baseURL, githubID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request for user %s: %w", githubID, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch keys for user %s: %w", githubID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, NewUserNotFoundError(githubID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d for user %s", resp.StatusCode, githubID)
	}

	var keys []Key
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, fmt.Errorf("decode response for user %s: %w", githubID, err)
	}

	sshKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		sshKeys = append(sshKeys, key.Key)
	}

	return sshKeys, nil
}
