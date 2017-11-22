// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package github

import (
	"bytes"
	"text/template"

	"github.com/mattermost/mattermost-server/model"
)

const (
	GITHUB_PULL_REQUEST_WEBHOOK = "pull_request"
)

// PullRequestWebhook - Handles Pull Request webhook events from Github.
type PullRequestWebhook struct {
	Action       string
	Number       int
	Sender       User
	Repository   Repository
	Organization Organization
	PullRequest  PullRequest `json:"pull_request"`
}

func init() {
	RegisterWebhook(GITHUB_PULL_REQUEST_WEBHOOK, &PullRequestWebhook{})
}

// SlackAttachment Returns the text to be placed in the resulting post or an empty string if nothing should be
// posted.
func (w *PullRequestWebhook) SlackAttachment() (*model.SlackAttachment, error) {

	text, err := w.renderText("" +
		"{{.PullRequest.User.Login}} {{.Action}} Pull Request #{{.Number}} " +
		"[{{.PullRequest.Title}}]({{.PullRequest.URL}})" +
		"")
	if err != nil {
		return nil, err
	}

	var fields []*model.SlackAttachmentField
	if w.PullRequest.Assignee != nil {
		fields = append(fields, &model.SlackAttachmentField{
			Title: "Assignees",
			Value: w.PullRequest.GetAssignees(),
			Short: true,
		})
	}
	fields = append(fields, &model.SlackAttachmentField{
		Title: "State",
		Value: w.PullRequest.State,
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Created At",
		Value: w.PullRequest.CreatedAt,
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Updated At",
		Value: w.PullRequest.UpdatedAt,
		Short: true,
	})

	return &model.SlackAttachment{
		Fallback: text,
		Color:    "#95b7d0",
		Pretext:  "",
		Text:     text,
		Fields:   fields,
	}, nil
}

func (w *PullRequestWebhook) renderText(tplBody string) (string, error) {
	tpl, err := template.New("post").Parse(tplBody)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, struct {
		*PullRequestWebhook
	}{
		PullRequestWebhook: w,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
