package commands

import (
	"context"
	"fmt"

	"atonixcorp/cloud/core/keystone"
)

func IdentityInfo() {
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

	tok := ks.Token()

	fmt.Println("User ID:", tok.UserID)
	fmt.Println("Project ID:", tok.ProjectID)
	fmt.Println("Roles:", tok.Roles)
}
