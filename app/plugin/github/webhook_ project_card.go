// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package github

import (
	"bytes"
	"text/template"

	"github.com/mattermost/mattermost-server/model"
)

const (
	GITHUB_PROJECT_CARD_WEBHOOK = "project_card"
)

// ProjectCardWebhook - Handles Project Card Webhook events from Github.
type ProjectCardWebhook struct {
	Action       string
	Sender       User
	Repository   Repository
	Organization Organization
	ProjectCard  ProjectCard `json:"project_card"`
}

func init() {
	RegisterWebhook(GITHUB_PROJECT_CARD_WEBHOOK, &ProjectCardWebhook{})
}

// SlackAttachment Returns the text to be placed in the resulting post or an empty string if nothing should be
// posted.
func (w *ProjectCardWebhook) SlackAttachment() (*model.SlackAttachment, error) {

	text, err := w.renderText("" +
		"{{.ProjectCard.Creator.Login}} {{.Action}} Project Card #{{.ProjectCard.ID}} " +
		"")
	if err != nil {
		return nil, err
	}

	var fields []*model.SlackAttachmentField

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Title",
		Value: w.ProjectCard.Title(),
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Created At",
		Value: w.ProjectCard.CreatedAt,
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Updated At",
		Value: w.ProjectCard.UpdatedAt,
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

func (w *ProjectCardWebhook) renderText(tplBody string) (string, error) {
	tpl, err := template.New("post").Parse(tplBody)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, struct {
		*ProjectCardWebhook
		Title string
	}{
		ProjectCardWebhook: w,
		Title:              w.ProjectCard.Title(),
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
