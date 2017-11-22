// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package github

import (
	"bytes"
	"text/template"

	"github.com/mattermost/mattermost-server/model"
)

const (
	GITHUB_ISSUES_WEBHOOK = "issues"
)

// IssuesWebhook - Handles Issue Webhook events from Github.
type IssuesWebhook struct {
	Action       string
	Sender       User
	Repository   Repository
	Organization Organization
	Issue        Issue
}

func init() {
	RegisterWebhook(GITHUB_ISSUES_WEBHOOK, &IssuesWebhook{})
}

//SlackAttachment  Returns the text to be placed in the resulting post or an empty string if nothing should be
// posted.
func (w *IssuesWebhook) SlackAttachment() (*model.SlackAttachment, error) {

	text, err := w.renderText("" +
		"{{.Issue.User.Login}} {{.Action}} issue #{{.Issue.ID}} " +
		"[{{.Issue.Title}}]({{.Issue.URL}})" +
		"")
	if err != nil {
		return nil, err
	}

	var fields []*model.SlackAttachmentField

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Assignees",
		Value: w.Issue.GetAssignees(),
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Labels",
		Value: w.Issue.GetLabels(),
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Created At",
		Value: w.Issue.CreatedAt,
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Updated At",
		Value: w.Issue.UpdatedAt,
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

func (w *IssuesWebhook) renderText(tplBody string) (string, error) {
	tpl, err := template.New("post").Parse(tplBody)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, struct {
		*IssuesWebhook
	}{
		IssuesWebhook: w,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
