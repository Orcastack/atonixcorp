package command

import (
	"context"
	"fmt"
	"os"
	"strconv"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func CreateStorage(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: atcloud create storage <name> <sizeGB>")
		return
	}

	size, _ := strconv.Atoi(args[1])

	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	req := sdk.CreateVolumeRequest{
		Name:   args[0],
		SizeGB: size,
	}

	res, err := client.CreateStorage(context.Background(), req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Volume created: %s (%dGB)\n", res.ID, res.SizeGB)
}

func DeleteStorage(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: atcloud delete storage <id>")
		return
	}

	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	err := client.DeleteStorage(context.Background(), args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Volume deleted:", args[0])
}

func ListStorage() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	list, err := client.ListStorage(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, v := range list {
		fmt.Printf("%s  %-10dGB  %s\n", v.ID, v.SizeGB, v.State)
	}
}
