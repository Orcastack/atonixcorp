package commands

import (
	"atonixcorp/cloud/core/cinder"
	"atonixcorp/cloud/core/keystone"
	"context"
	"fmt"
)

func CreateStorage(args []string) {
	name := args[0]
	size := 10 // default or parse flags

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

	cd, err := cinder.NewClient(ks, cinder.Config{
		Region:    "RegionOne",
		Interface: "public",
	})
	if err != nil {
		fmt.Println("Cinder client error:", err)
		return
	}

	vol, err := cd.CreateVolume(ctx, cinder.CreateVolumeRequest{
		Name: name,
		Size: size,
	})
	if err != nil {
		fmt.Println("Volume creation failed:", err)
		return
	}

	fmt.Println("Created volume:", vol.ID)
}

func ListStorage() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	cd, _ := cinder.NewClient(ks, cinder.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	vols, _ := cd.ListVolumes(ctx)
	fmt.Println("Volumes:", vols)
}
