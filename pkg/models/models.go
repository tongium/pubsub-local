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
	// AutoDiscover creates a capture subscription for topics created at
	// runtime (e.g. by your application). Defaults to true when omitted.
	AutoDiscover *bool `yaml:"auto_discover"`
	// DiscoverInterval is how often to look for new topics (Go duration, default "2s").
	DiscoverInterval string  `yaml:"discover_interval"`
	Topics           []Topic `yaml:"topics"`
}

func (s Settings) AutoDiscoverEnabled() bool {
	return s.AutoDiscover == nil || *s.AutoDiscover
}

type MessageRecord struct {
	PublishTime string            `json:"publish_time"`
	MessageID   string            `json:"message_id"`
	Attributes  map[string]string `json:"attributes"`
	OrderingKey string            `json:"ordering_key,omitempty"`
	Data        any               `json:"data"`
}
