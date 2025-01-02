package cli

import (
	"github.com/VladSatyshev/concurrent-queue-cli/client"
	"github.com/urfave/cli/v2"
)

func NewCliApp(client client.QueueClient) *cli.App {
	cmdExec := NewCommandExecutor(client)

	return &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "queue",
				Usage: "Queue commands",
				Subcommands: []*cli.Command{
					{
						Name:   "show",
						Usage:  "Show queues",
						Action: cmdExec.GetAllQueues,
						Subcommands: []*cli.Command{
							{
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "name",
										Usage:    "Name of queue",
										Required: true,
									},
								},
								Name:   "detail",
								Usage:  "Show queue details by name",
								Action: cmdExec.GetQueue,
							},
						},
					},
					{
						Name:   "add",
						Usage:  "Add message to queue",
						Action: cmdExec.AddMessage,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     NameFlag,
								Usage:    "Name of queue",
								Required: true,
							},
							&cli.StringFlag{
								Name:     MessageFlag,
								Usage:    "Message to add to queue",
								Required: true,
							},
						},
					},
					{
						Name:   "subscribe",
						Usage:  "Subscribe to queue",
						Action: cmdExec.AddSubscriber,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     NameFlag,
								Usage:    "Name of queue",
								Required: true,
							},
							&cli.StringFlag{
								Name:     SubscriberFlag,
								Usage:    "Subscriber name",
								Required: true,
							},
						},
					},
					{
						Name:   "consume",
						Usage:  "Consume messages from queue",
						Action: cmdExec.ConsumeMessages,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     NameFlag,
								Usage:    "Name of queue",
								Required: true,
							},
							&cli.StringFlag{
								Name:     SubscriberFlag,
								Usage:    "Subscriber name",
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}
