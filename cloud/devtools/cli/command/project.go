package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func ListProjects() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	projects, err := client.ListProjects(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, p := range projects {
		fmt.Printf("%s  %s\n", p.ID, p.Name)
	}
}
