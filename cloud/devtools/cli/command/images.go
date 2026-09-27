package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func ListImages() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	imgs, err := client.ListImages(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, img := range imgs {
		fmt.Printf("%s  %-20s  %s\n", img.ID, img.Name, img.OS)
	}
}
