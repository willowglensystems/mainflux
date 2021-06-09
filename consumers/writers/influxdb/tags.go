package influxdb

import (
	"git.willowglen.ca/sq/third-party/mainflux.git/pkg/transformers/json"
	"git.willowglen.ca/sq/third-party/mainflux.git/pkg/transformers/senml"
)

type tags map[string]string

func senmlTags(msg senml.Message) tags {
	return tags{
		"channel":   msg.Channel,
		"subtopic":  msg.Subtopic,
		"publisher": msg.Publisher,
		"name":      msg.Name,
	}
}

func jsonTags(msg json.Message) tags {
	return tags{
		"channel":   msg.Channel,
		"subtopic":  msg.Subtopic,
		"publisher": msg.Publisher,
	}
}
