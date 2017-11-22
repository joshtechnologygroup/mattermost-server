package github

import (
	"strings"
	"time"
)

type Item struct {
	ID        int64 `json:"id"`
	Title     string
	URL       string `json:"html_url"`
	State     string
	User      User
	Assignee  *User
	Assignees []User    `json:"assignees"`
	Labels    []Label   `json:"labels"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClosedAt  time.Time `json:"closed_at"`
}

func (p Item) GetAssignees() string {
	strs := make([]string, len(p.Assignees))
	for i, v := range p.Assignees {
		strs[i] = v.Login
	}
	return strings.Join(strs, ", ")
}

func (p Item) GetLabels() string {
	strs := make([]string, len(p.Labels))
	for i, v := range p.Labels {
		strs[i] = v.Name
	}
	return strings.Join(strs, ", ")
}

type Repository struct {
	Name string `json:"full_name"`
	URL  string `json:"html_url"`
}

type Organization struct {
	Login     string
	AvatarURL string `json:"avatar_url"`
}

type User struct {
	Login     string
	URL       string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
}

type Label struct {
	ID    int64 `json:"id"`
	Name  string
	Color string
}

type PullRequest struct {
	Item
}

type Issue struct {
	Item
	Body string
}

type ProjectCard struct {
	ID         int64  `json:"id"`
	ContentURL string `json:"content_url"`
	Note       string
	Creator    User
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (p ProjectCard) Title() string {
	title := p.Note
	if title == "" {
		title = p.ContentURL
	}
	return title
}
