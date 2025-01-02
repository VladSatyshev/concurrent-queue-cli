package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VladSatyshev/concurrent-queue-cli/client"
	"github.com/VladSatyshev/concurrent-queue-cli/internal/cli"
)

func main() {
	client := client.NewQueueClient(http.Client{
		Timeout: time.Second * 10,
	})

	app := cli.NewCliApp(client)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
