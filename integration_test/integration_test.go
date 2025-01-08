package test

import (
	"encoding/json"
	"testing"

	"github.com/VladSatyshev/concurrent-queue-cli/client"
	"github.com/VladSatyshev/concurrent-queue-cli/integration_test/config"
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

		res := cmdExecutor.Execute(t, "queue show")
		var actualQueues []client.Queue
		err := json.Unmarshal(res, &actualQueues)
		assert.Nil(t, err)
		assert.Equal(t, len(expectedQueues), len(actualQueues))
		for _, aq := range actualQueues {
			for _, eq := range expectedQueues {
				if aq.Name == eq.Name {
					assert.Equal(t, eq, aq)
				}
			}
		}
	})
	t.Run("Shoud show created queues2", func(t *testing.T) {
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

		res := cmdExecutor.Execute(t, "queue show")
		var actualQueues []client.Queue
		err := json.Unmarshal(res, &actualQueues)
		assert.Nil(t, err)
		assert.Equal(t, len(expectedQueues), len(actualQueues))
		for _, aq := range actualQueues {
			for _, eq := range expectedQueues {
				if aq.Name == eq.Name {
					assert.Equal(t, eq, aq)
				}
			}
		}
	})

}
