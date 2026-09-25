package commands

import (
	"atonixcorp/cloud/core/glance"
	"atonixcorp/cloud/core/keystone"
	"context"
	"fmt"
)

func ListImages() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	gl, _ := glance.NewClient(ks, glance.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	imgs, _ := gl.ListImages(ctx)
	fmt.Println("Images:", imgs)
}
