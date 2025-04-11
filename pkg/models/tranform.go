package models

import (
	"strings"
)

func (t *WebhookPayload) ToNotifications(topic string) (notifications []*Notification) {
	if len(t.Alerts) == 0 {
		var n Notification

		n.Topic = topic
		n.Message = t.Message
		n.Title = t.Title

		return []*Notification{&n}
	}

	for _, alert := range t.Alerts {
		var n Notification

		n.Topic = topic
		n.Message = t.Message
		n.Title = t.Title

		priority, ok := alert.Labels["priority"]
		if ok {
			n.Priority = toPriority(priority)
		}

		if alert.GeneratorURL != "" {
			n.Actions = append(n.Actions, &Action{
				Action: "view",
				Label:  "Open Alert in Grafana",
				Url:    alert.GeneratorURL,
			})
		}

		if alert.SilenceURL != "" {
			n.Actions = append(n.Actions, &Action{
				Action: "view",
				Label:  "Silence Alert",
				Url:    alert.SilenceURL,
				Clear:  true,
			})
		}

		notifications = append(notifications, &n)
	}

	return notifications
}

func toPriority(label string) int {
	switch strings.TrimSpace(strings.ToLower(label)) {
	case "1", "min", "minimum", "lowest":
		return 1
	case "2", "low":
		return 2
	case "4", "high":
		return 4
	case "5", "max", "urgent", "highest":
		return 5
	default:
		return 3
	}
}
