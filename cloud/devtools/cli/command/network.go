package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func CreateNetwork(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: atcloud create network <name> <cidr>")
		return
	}

	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	req := sdk.CreateNetworkRequest{
		Name: args[0],
		CIDR: args[1],
	}

	res, err := client.CreateNetwork(context.Background(), req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Network created: %s (%s)\n", res.ID, res.CIDR)
}

func DeleteNetwork(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: atcloud delete network <id>")
		return
	}

	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	err := client.DeleteNetwork(context.Background(), args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Network deleted:", args[0])
}

func ListNetwork() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	list, err := client.ListNetwork(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, n := range list {
		fmt.Printf("%s  %-15s  %s\n", n.ID, n.CIDR, n.State)
	}
}
