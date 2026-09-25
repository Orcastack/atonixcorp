package commands

import (
	"context"
	"fmt"

	"atonixcorp/cloud/core/keystone"
	"atonixcorp/cloud/core/nova"
)

func ListCompute() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	nv, _ := nova.NewClient(ks, nova.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	servers, _ := nv.ListServers(ctx)
	fmt.Println("Servers:", servers)
}

func CreateCompute(args []string) {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	nv, _ := nova.NewClient(ks, nova.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	srv, err := nv.CreateServer(ctx, nova.CreateServerRequest{
		Name:      "example-server",
		FlavorRef: "flavor-id",
		ImageRef:  "image-id",
		NetworkID: "network-id",
	})
	if err != nil {
		fmt.Println("Create server failed:", err)
		return
	}

	fmt.Println("Created server:", srv.ID)
}
