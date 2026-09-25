package commands

import (
	"context"
	"fmt"
	"os"

	"atonixcorp/cloud/core/keystone"
	"atonixcorp/cloud/core/swift"
)

func ListContainers() {
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	sw, _ := swift.NewClient(ks, swift.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	containers, err := sw.ListContainers(ctx)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	fmt.Println("Containers:")
	for _, c := range containers {
		fmt.Printf("- %s (%d objects, %d bytes)\n", c.Name, c.Count, c.Bytes)
	}
}

func CreateContainer(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: atcloud create container <name>")
		return
	}

	name := args[0]
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	sw, _ := swift.NewClient(ks, swift.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	if err := sw.CreateContainer(ctx, name); err != nil {
		fmt.Println("Failed:", err)
		return
	}

	fmt.Println("Created container:", name)
}

func DeleteContainer(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: atcloud delete container <name>")
		return
	}

	name := args[0]
	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	sw, _ := swift.NewClient(ks, swift.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	if err := sw.DeleteContainer(ctx, name); err != nil {
		fmt.Println("Failed:", err)
		return
	}

	fmt.Println("Deleted container:", name)
}

func UploadObject(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: atcloud upload object <container> <name> <file>")
		return
	}

	container := args[0]
	name := args[1]
	file := args[2]

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("File error:", err)
		return
	}
	defer f.Close()

	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	sw, _ := swift.NewClient(ks, swift.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	if err := sw.UploadObject(ctx, container, name, f); err != nil {
		fmt.Println("Upload failed:", err)
		return
	}

	fmt.Println("Uploaded object:", name)
}

func DeleteObject(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: atcloud delete object <container> <name>")
		return
	}

	container := args[0]
	name := args[1]

	ctx := context.Background()

	ks := keystone.NewClient(keystone.Credentials{
		AuthURL:   "https://keystone.example.com/v3",
		Username:  "samuel",
		Password:  "secret",
		ProjectID: "project-id",
		DomainID:  "default",
	})

	_ = ks.AuthenticatePassword(ctx)

	sw, _ := swift.NewClient(ks, swift.Config{
		Region:    "RegionOne",
		Interface: "public",
	})

	if err := sw.DeleteObject(ctx, container, name); err != nil {
		fmt.Println("Delete failed:", err)
		return
	}

	fmt.Println("Deleted object:", name)
}
