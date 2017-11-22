// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package github

import "github.com/mattermost/mattermost-server/model"

type Webhook interface {
	SlackAttachment() (*model.SlackAttachment, error)
	renderText(tplBody string) (string, error)
}

var webhooks = make(map[string]Webhook)

func RegisterWebhook(name string, newWebhook Webhook) {
	webhooks[name] = newWebhook
}

func GetWebhook(name string) Webhook {

	webhook, ok := webhooks[name]
	if ok {
		return webhook
	}
	return nil
}
