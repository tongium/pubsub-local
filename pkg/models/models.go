package models

type Subscriber struct {
	ID                     string `yaml:"id"`
	EnabledMessageOrdering bool   `yaml:"enabled_message_ordering"`
}

type Topic struct {
	ID          string       `yaml:"id"`
	Subscribers []Subscriber `yaml:"subscribers"`
}

type Settings struct {
	Topics []Topic `yaml:"topics"`
}

type MessageRecord struct {
	PublishTime string            `json:"publish_time"`
	MessageID   string            `json:"message_id"`
	Attributes  map[string]string `json:"attributes"`
	OrderingKey string            `json:"ordering_key,omitempty"`
	Data        any               `json:"data"`
}
