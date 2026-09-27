package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func IdentityInfo() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	info, err := client.IdentityInfo(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("User: %s (%s)\n", info.Username, info.UserID)
	fmt.Println("Roles:", info.Roles)
	fmt.Println("Projects:", info.Projects)
}
