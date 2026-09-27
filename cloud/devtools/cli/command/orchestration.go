package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func ListEvents() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	events, err := client.ListEvents(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, e := range events {
		fmt.Printf("[%s] %s\n", e.Type, e.Message)
	}
}

func ListMetrics() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	metrics, err := client.ListMetrics(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, m := range metrics {
		fmt.Printf("%s = %.2f\n", m.Name, m.Value)
	}
}
