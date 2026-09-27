package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func CreateCompute(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: atcloud create compute <name> <plan>")
		return
	}

	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	req := sdk.CreateComputeRequest{
		Name: args[0],
		Plan: args[1],
	}

	res, err := client.CreateCompute(context.Background(), req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Compute created: %s (%s)\n", res.ID, res.State)
}

func DeleteCompute(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: atcloud delete compute <id>")
		return
	}

	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	err := client.DeleteCompute(context.Background(), args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Compute deleted:", args[0])
}

func ListCompute() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	list, err := client.ListCompute(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, c := range list {
		fmt.Printf("%s  %-10s  %s\n", c.ID, c.Plan, c.State)
	}
}
