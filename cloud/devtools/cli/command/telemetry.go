package command

import (
	"context"
	"fmt"
	"os"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func ListTelemetry() {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	snap, err := client.Telemetry(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("CPU: %.2f%%\n", snap.CPUUsage)
	fmt.Printf("Memory: %.2f%%\n", snap.MemoryUsage)
	fmt.Printf("Disk: %.2f%%\n", snap.DiskUsage)
}
