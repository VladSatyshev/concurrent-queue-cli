package client

const (
	url  = "http://127.0.0.1"
	port = "8000"

	intPrefix = "/int"
	v1Prefix  = "/v1"

	baseV1    = url + ":" + port + v1Prefix
	baseIntV1 = baseV1 + intPrefix

	QueuesEndpoint        = baseIntV1 + "/queues"
	MessagesEndpoint      = baseV1 + "/queues/%s/messages"
	SubscriptionsEndpoint = baseV1 + "/queues/%s/subscriptions"
)

type Queue struct {
	Name           string                  `json:"name"`
	MaxLength      uint                    `json:"maxLength"`
	MaxSubscribers uint                    `json:"maxSubscribers"`
	Subscribers    map[string]struct{}     `json:"subscribers,omitempty"`
	Messages       map[string]QueueMessage `json:"messages,omitempty"`
}

type QueueMessage struct {
	Body   map[string]interface{}
	SeenBy map[string]struct{}
}
