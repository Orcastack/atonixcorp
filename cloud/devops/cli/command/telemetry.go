package commands

import (
	"context"
	"fmt"

	"atonixcorp/cloud/core/keystone"
	"atonixcorp/cloud/core/telemetry"
)

func ListMetrics() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	tm, err := telemetry.NewClient(ks, telemetry.Config{
		Region:    "RegionOne",
		Interface: "public",
	})
	if err != nil {
		fmt.Println("Telemetry client error:", err)
		return
	}

	metrics, err := tm.ListMetrics(ctx)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	fmt.Println("Metrics:")
	for _, m := range metrics {
		fmt.Printf("- %s (%s)\n", m.Name, m.Unit)
	}
}

func ListEvents() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	tm, err := telemetry.NewClient(ks, telemetry.Config{
		Region:    "RegionOne",
		Interface: "public",
	})
	if err != nil {
		fmt.Println("Telemetry client error:", err)
		return
	}

	events, err := tm.ListEvents(ctx)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	fmt.Println("Events:")
	for _, e := range events {
		fmt.Printf("- %s @ %s\n", e.EventType, e.Timestamp)
	}
}
