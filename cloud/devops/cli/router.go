package router

import (
	"fmt"

	"atonixcorp/cloud/devops/cli/commands"
)

func Dispatch(args []string) {
	if len(args) < 3 {
		commands.Help()
		return
	}

	action := args[1]
	resource := args[2]

	switch action {

	// ------------------------------------------------------------
	// CREATE
	// ------------------------------------------------------------
	case "create":
		switch resource {
		case "compute":
			commands.CreateCompute(args[3:])
		case "network":
			commands.CreateNetwork(args[3:])
		case "storage":
			commands.CreateStorage(args[3:])
		case "container":
			commands.CreateContainer(args[3:])
		default:
			fmt.Println("Unknown resource:", resource)
		}

	// ------------------------------------------------------------
	// DELETE
	// ------------------------------------------------------------
	case "delete":
		switch resource {
		case "compute":
			commands.DeleteCompute(args[3:])
		case "network":
			commands.DeleteNetwork(args[3:])
		case "storage":
			commands.DeleteStorage(args[3:])
		case "container":
			commands.DeleteContainer(args[3:])
		case "object":
			commands.DeleteObject(args[3:])
		default:
			fmt.Println("Unknown resource:", resource)
		}

	// ------------------------------------------------------------
	// LIST
	// ------------------------------------------------------------
	case "list":
		switch resource {
		case "compute":
			commands.ListCompute()
		case "network":
			commands.ListNetwork()
		case "storage":
			commands.ListStorage()
		case "projects":
			commands.ListProjects()
		case "containers":
			commands.ListContainers()
		case "metrics":
			commands.ListMetrics()
		case "events":
			commands.ListEvents()
		default:
			fmt.Println("Unknown resource:", resource)
		}

	// ------------------------------------------------------------
	// UPLOAD
	// ------------------------------------------------------------
	case "upload":
		if resource == "object" {
			commands.UploadObject(args[3:])
		} else {
			fmt.Println("Unknown upload target:", resource)
		}

	// ------------------------------------------------------------
	// IDENTITY
	// ------------------------------------------------------------
	case "identity":
		commands.IdentityInfo()

	// ------------------------------------------------------------
	// HELP (fallback)
	// ------------------------------------------------------------
	default:
		commands.Help()
	}
}
