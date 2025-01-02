package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type QueueClient interface {
	GetByName(ctx context.Context, queueName string) (interface{}, error)
	GetAll(ctx context.Context) (interface{}, error)
	AddMessage(ctx context.Context, queueName string, jsonBody string) (interface{}, error)
	AddSubscriber(ctx context.Context, queueName string, subscriberName string) (interface{}, error)
	ConsumeMessages(ctx context.Context, queueName string, subscriberName string) (interface{}, error)
}

type queueClientImpl struct {
	client http.Client
}

func NewQueueClient(client http.Client) QueueClient {
	return &queueClientImpl{
		client: client,
	}
}

func (c *queueClientImpl) GetByName(ctx context.Context, queueName string) (interface{}, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/%s", QueuesEndpoint, queueName))
	if err != nil {
		return Queue{}, err
	}
	defer resp.Body.Close()

	var res Queue
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return Queue{}, err
	}

	return res, nil
}

func (c *queueClientImpl) GetAll(ctx context.Context) (interface{}, error) {
	resp, err := c.client.Get(QueuesEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res []Queue
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *queueClientImpl) AddMessage(ctx context.Context, queueName string, msg string) (interface{}, error) {
	jsonBody := map[string]string{
		"msg": msg,
	}
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(fmt.Sprintf(MessagesEndpoint, queueName), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *queueClientImpl) AddSubscriber(ctx context.Context, queueName string, subscriberName string) (interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(SubscriptionsEndpoint, queueName), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Subscriber", subscriberName)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *queueClientImpl) ConsumeMessages(ctx context.Context, queueName string, subscriberName string) (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(MessagesEndpoint, queueName), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Subscriber", subscriberName)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
