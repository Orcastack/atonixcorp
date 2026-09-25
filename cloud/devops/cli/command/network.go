package commands

import (
	"context"
	"fmt"

	"atonixcorp/cloud/core/keystone"
	"atonixcorp/cloud/core/neutron"
)

func ListNetwork() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	nt, _ := neutron.NewClient(ks, neutron.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	nets, _ := nt.ListNetworks(ctx)
	fmt.Println("Networks:", nets)
}

func CreateNetwork(args []string) {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	nt, _ := neutron.NewClient(ks, neutron.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	net, err := nt.CreateNetwork(ctx, neutron.CreateNetworkRequest{
		Name:         "example-network",
		AdminStateUp: true,
	})
	if err != nil {
		fmt.Println("Create network failed:", err)
		return
	}

	fmt.Println("Created network:", net.ID)
}
