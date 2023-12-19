package main

import "time"

type Issue struct {
	HtmlUrl     string       `json:"html_url"`
	Number      int          `json:"number"`
	Title       string       `json:"title"`
	Labels      []Label      `json:"labels"`
	State       string       `json:"state"`
	Comments    int          `json:"comments"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Body        string       `json:"body"`
	Reactions   Reactions    `json:"reactions"`
	PullRequest *PullRequest `json:"pull_request"`
}

type Label struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type Reactions struct {
	TotalCount int `json:"total_count"`
}

type PullRequest struct {
	HtmlUrl string `json:"html_url"`
}
