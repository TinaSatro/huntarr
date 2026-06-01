package outreach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const claudeAPI = "https://api.anthropic.com/v1/messages"

type ClaudeGenerator struct {
	apiKey string
	client *http.Client
}

func NewClaudeGenerator(apiKey string) *ClaudeGenerator {
	return &ClaudeGenerator{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *ClaudeGenerator) Generate(jobTitle, company, jobDescription string) (Message, error) {
	prompt := fmt.Sprintf(`You are helping a Senior/Staff SRE and Infrastructure Architect reach out to recruiters.

Write a short, professional recruiter outreach message for this job:
COMPANY: %s
TITLE: %s
DESCRIPTION: %s

The message should be:
- 3-4 sentences max
- Mention relevant experience: Kubernetes, RKE2, airgap edge systems, infrastructure architecture
- Sound human, not templated
- End with a call to action

Respond ONLY with a raw JSON object. No markdown. No backticks. No explanation. Just the JSON:
{"subject": "...", "body": "..."}`,
		company, jobTitle, jobDescription)

	body, _ := json.Marshal(map[string]any{
		"model":      "claude-sonnet-4-5",
		"max_tokens": 1000,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, err := http.NewRequest("POST", claudeAPI, bytes.NewReader(body))
	if err != nil {
		return Message{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Message{}, err
	}

	if len(result.Content) == 0 {
		return Message{}, fmt.Errorf("empty response from Claude")
	}

	var msg struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := json.Unmarshal([]byte(result.Content[0].Text), &msg); err != nil {
		return Message{}, fmt.Errorf("parse response: %w", err)
	}

	return Message{
		Subject: msg.Subject,
		Body:    msg.Body,
		Company: company,
		JobURL:  "",
	}, nil
}
