package main

import (
	"context"
	"fmt"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func main() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: "https://api.atcloud.world/v1",
		Token:    "YOUR_TOKEN",
	})

	// List compute instances
	compute, _ := client.ListCompute(context.Background())
	fmt.Println("Compute:", compute)

	// Identity info
	id, _ := client.IdentityInfo(context.Background())
	fmt.Println("Identity:", id)

	// Telemetry
	tele, _ := client.Telemetry(context.Background())
	fmt.Println("Telemetry:", tele)
}
