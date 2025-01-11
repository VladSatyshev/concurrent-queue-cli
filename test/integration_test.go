package test

import (
	"encoding/json"
	"testing"

	"github.com/VladSatyshev/concurrent-queue-cli/client"
	"github.com/VladSatyshev/concurrent-queue-cli/test/config"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationWithServer(t *testing.T) {
	buildCli(t, buildPath)

	t.Run("Shoud show created queues", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue1",
				Length:            1,
				SubscribersAmount: 1,
			},
			{
				Name:              "testQueue2",
				Length:            2,
				SubscribersAmount: 2,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		expectedQueues := []client.Queue{
			{

				Name:           "testQueue1",
				MaxLength:      1,
				MaxSubscribers: 1,
			},
			{
				Name:           "testQueue2",
				MaxLength:      2,
				MaxSubscribers: 2,
			},
		}

		res := cmdExecutor.Execute(t, `queue show`)
		var actualQueues []client.Queue
		err := json.Unmarshal(res, &actualQueues)
		assert.Nil(t, err)
		assertEqualQueues(t, expectedQueues, actualQueues)
	})

	t.Run("Shoud get queue by name", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            1,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		expectedQueues := []client.Queue{
			{
				Name:           "testQueue",
				MaxLength:      1,
				MaxSubscribers: 1,
			},
		}

		res := cmdExecutor.Execute(t, `queue show detail --name testQueue`)

		actualQueues := make([]client.Queue, 1)
		err := json.Unmarshal(res, &actualQueues[0])
		assert.Nil(t, err)
		assertEqualQueues(t, expectedQueues, actualQueues)
	})

	t.Run("Shoud be able to subscribe to queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            1,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		expectedQueues := []client.Queue{
			{
				Name:           "testQueue",
				MaxLength:      1,
				MaxSubscribers: 1,
				Subscribers:    map[string]struct{}{"user": {}},
			},
		}

		cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user`)
		res := cmdExecutor.Execute(t, `queue show detail --name testQueue`)

		actualQueues := make([]client.Queue, 1)
		err := json.Unmarshal(res, &actualQueues[0])
		assert.Nil(t, err)
		assertEqualQueues(t, expectedQueues, actualQueues)
	})

	t.Run("Shoud not be able to subscribe to queue twice", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            1,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user`)
		res := cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user`)
		assert.Contains(t, string(res), "user user has already subscribed to queue")
	})

	t.Run("Shoud not be able to subscribe to full queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            1,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user1`)
		res := cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user2`)
		assert.Contains(t, string(res), "too many subscribers")
	})

	t.Run("Shoud not be able to subscribe to not existent queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		res := cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user`)
		assert.Contains(t, string(res), "Queue not found")
	})

	t.Run("Shoud be able to add message to queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            2,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		expectedMessages := []map[string]interface{}{
			{"msg": "test1"},
			{"msg": "test2"},
		}

		cmdExecutor.Execute(t, `queue add --name testQueue --message test1`)
		cmdExecutor.Execute(t, `queue add --name testQueue --message test2`)
		res := cmdExecutor.Execute(t, `queue show detail --name testQueue`)

		actualQueues := make([]client.Queue, 1)
		err := json.Unmarshal(res, &actualQueues[0])
		assert.Nil(t, err)
		assert.Equal(t, 1, len(actualQueues))
		actualMessagesBodies := make([]map[string]interface{}, 0)
		for _, v := range actualQueues[0].Messages {
			actualMessagesBodies = append(actualMessagesBodies, v.Body)
		}
		assertEqualMessages(t, expectedMessages, actualMessagesBodies)
	})

	t.Run("Shoud not be able to add message to full queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            1,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		cmdExecutor.Execute(t, `queue add --name testQueue --message test1`)
		res := cmdExecutor.Execute(t, `queue add --name testQueue --message test2`)
		assert.Contains(t, string(res), "too many messages")
	})

	t.Run("Shoud be able to consume messages from queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            2,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		expectedMessages := []map[string]interface{}{
			{"msg": "test1"},
			{"msg": "test2"},
		}

		expectedMessagesEmpty := []map[string]interface{}{}

		cmdExecutor.Execute(t, `queue add --name testQueue --message test1`)
		cmdExecutor.Execute(t, `queue add --name testQueue --message test2`)
		cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user`)
		res := cmdExecutor.Execute(t, `queue consume --name testQueue --subscriber user`)

		actualMessages := map[string]interface{}{}
		err := json.Unmarshal(res, &actualMessages)
		assert.Nil(t, err)
		actualMessagesBodies := make([]map[string]interface{}, 0)
		for _, v := range actualMessages {
			actualMessagesBodies = append(actualMessagesBodies, v.(map[string]interface{}))
		}
		assertEqualMessages(t, expectedMessages, actualMessagesBodies)

		resEmpty := cmdExecutor.Execute(t, `queue consume --name testQueue --subscriber user`)
		actualMessagesEmpty := map[string]interface{}{}
		err = json.Unmarshal(resEmpty, &actualMessagesEmpty)
		assert.Nil(t, err)
		actualMessagesBodiesEmpty := make([]map[string]interface{}, 0)
		for _, v := range actualMessagesEmpty {
			actualMessagesBodiesEmpty = append(actualMessagesBodiesEmpty, v.(map[string]interface{}))
		}
		assertEqualMessages(t, expectedMessagesEmpty, actualMessagesBodiesEmpty)
	})

	t.Run("Not a subscriber should not be able to consume messages from queue", func(t *testing.T) {
		qConfigs := []config.QueueConfig{
			{
				Name:              "testQueue",
				Length:            1,
				SubscribersAmount: 1,
			},
		}
		cmdExecutor, teardown := configureEnvironment(t, qConfigs)
		defer teardown()

		cmdExecutor.Execute(t, `queue add --name testQueue --message test1`)
		cmdExecutor.Execute(t, `queue subscribe --name testQueue --subscriber user`)
		res := cmdExecutor.Execute(t, `queue consume --name testQueue --subscriber not_a_sub`)
		assert.Contains(t, string(res), "doesn't have subscriber")
	})

}
