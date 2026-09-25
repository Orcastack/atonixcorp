package commands

import (
	"context"
	"fmt"

	"atonixcorp/cloud/core/keystone"
)

func ListProjects() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	if err := ks.AuthenticatePassword(ctx); err != nil {
		fmt.Println("Auth failed:", err)
		return
	}

	projects, err := ks.ListProjects(ctx)
	if err != nil {
		fmt.Println("Failed to list projects:", err)
		return
	}

	fmt.Println("Projects:")
	for _, p := range projects {
		fmt.Printf("- %s (%s)\n", p.Name, p.ID)
	}
}
