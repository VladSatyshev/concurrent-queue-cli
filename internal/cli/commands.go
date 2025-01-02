package cli

import (
	"github.com/VladSatyshev/concurrent-queue-cli/client"
	"github.com/urfave/cli/v2"
)

type commandExecutor struct {
	client client.QueueClient
}

func NewCommandExecutor(client client.QueueClient) *commandExecutor {
	return &commandExecutor{
		client: client,
	}
}

func (e *commandExecutor) GetAllQueues(c *cli.Context) error {
	queues, err := e.client.GetAll(c.Context)
	if err != nil {
		return err
	}

	if err := prettyPrint(queues); err != nil {
		return err
	}
	return nil
}

func (e *commandExecutor) GetQueue(c *cli.Context) error {
	name := c.String(NameFlag)
	queue, err := e.client.GetByName(c.Context, name)
	if err != nil {
		return err
	}

	if err := prettyPrint(queue); err != nil {
		return err
	}
	return nil
}

func (e *commandExecutor) AddMessage(c *cli.Context) error {
	name := c.String(NameFlag)
	msg := c.String(MessageFlag)

	retMsg, err := e.client.AddMessage(c.Context, name, msg)
	if err != nil {
		return err
	}

	if err := prettyPrint(retMsg); err != nil {
		return err
	}
	return nil
}

func (e *commandExecutor) AddSubscriber(c *cli.Context) error {
	name := c.String(NameFlag)
	subscriber := c.String(SubscriberFlag)

	msg, err := e.client.AddSubscriber(c.Context, name, subscriber)
	if err != nil {
		return err
	}

	if err := prettyPrint(msg); err != nil {
		return err
	}
	return nil
}

func (e *commandExecutor) ConsumeMessages(c *cli.Context) error {
	name := c.String(NameFlag)
	subscriber := c.String(SubscriberFlag)

	msgs, err := e.client.ConsumeMessages(c.Context, name, subscriber)
	if err != nil {
		return err
	}

	if err := prettyPrint(msgs); err != nil {
		return err
	}
	return nil
}
