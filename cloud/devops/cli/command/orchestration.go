package commands

import (
	"context"
	"fmt"

	"atonixcorp/cloud/core/heat"
	"atonixcorp/cloud/core/keystone"
)

func ListStacks() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	ht, _ := heat.NewClient(ks, heat.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	stacks, _ := ht.ListStacks(ctx)
	fmt.Println("Stacks:", stacks)
}
