package redmine

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/epistax1s/gomer/internal/config"
	"github.com/epistax1s/gomer/internal/log"
)

type RedmineClient struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

func NewRedmineClient(conf *config.RedmineConfig) *RedmineClient {
	return &RedmineClient{
		BaseURL: conf.BaseURL,
		ApiKey:  conf.ApiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *RedmineClient) GetTimeEntries(redUserID int, date string) ([]TimeEntry, error) {
	apiURL := fmt.Sprintf("%s/time_entries.json", c.BaseURL)

	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	query.Set("user_id", fmt.Sprintf("%d", redUserID))
	query.Set("spent_on", date)
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Redmine-API-Key", c.ApiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Error("Error while requesting Redmine", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var response TimeEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.TimeEntries, nil
}

func (c *RedmineClient) GetIssue(issueID int) (string, error) {
	url := fmt.Sprintf("%s/issues/%d.json", c.BaseURL, issueID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Redmine-API-Key", c.ApiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("issue %d: статус %d: %s", issueID, resp.StatusCode, string(body))
	}

	var result struct {
		Issue struct {
			Subject string `json:"subject"`
		} `json:"issue"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Issue.Subject, nil
}
