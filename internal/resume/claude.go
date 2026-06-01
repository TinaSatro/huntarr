package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const claudeAPI = "https://api.anthropic.com/v1/messages"

type ClaudeAdapter struct {
	apiKey string
	client *http.Client
}

func NewClaudeAdapter(apiKey string) *ClaudeAdapter {
	return &ClaudeAdapter{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *ClaudeAdapter) Adapt(base BaseResume, jobTitle, company, jobDescription string) (AdaptedResume, error) {
	prompt := fmt.Sprintf(`You are a resume writing expert.

Given this base resume:
HEADLINE: %s
SUMMARY: %s

Rewrite the headline and summary to match this job:
COMPANY: %s
TITLE: %s
DESCRIPTION: %s

Respond ONLY with a raw JSON object. No markdown. No backticks. No explanation. Just the JSON:
{"headline": "...", "summary": "..."}`,
		base.Headline, base.Summary, company, jobTitle, jobDescription)

	body, _ := json.Marshal(map[string]any{
		"model":      "claude-sonnet-4-5",
		"max_tokens": 1000,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, err := http.NewRequest("POST", claudeAPI, bytes.NewReader(body))
	if err != nil {
		return AdaptedResume{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return AdaptedResume{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AdaptedResume{}, err
	}

	if len(result.Content) == 0 {
		return AdaptedResume{}, fmt.Errorf("empty response from Claude")
	}

	var adapted struct {
		Headline string `json:"headline"`
		Summary  string `json:"summary"`
	}
	if err := json.Unmarshal([]byte(result.Content[0].Text), &adapted); err != nil {
		return AdaptedResume{}, fmt.Errorf("parse response: %w", err)
	}

	return AdaptedResume{
		Headline: adapted.Headline,
		Summary:  adapted.Summary,
		JobURL:   "",
		JobTitle: jobTitle,
		Company:  company,
	}, nil
}
