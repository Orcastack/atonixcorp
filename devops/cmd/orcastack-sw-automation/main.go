package main

import (
	"context"
	"log"

	"atonixcorp/devops/internal/platform/app"
	"atonixcorp/devops/internal/platform/config"
	"atonixcorp/devops/internal/swautomation"
)

func main() {
	err := app.Run(context.Background(), app.WithServiceSecurity(app.Config{
		Name:               "orcastack-sw-automation",
		Role:               "software-automation",
		Summary:            "Executes software delivery workflows, package publication, deployment steps, and cloud-native automation jobs.",
		RegisterHTTPRoutes: swautomation.Register,
		HTTPPort:           config.String("ORCASTACK_SW_AUTOMATION_HTTP_PORT", "8088"),
		GRPCPort:           config.String("ORCASTACK_SW_AUTOMATION_GRPC_PORT", "9088"),
	}, "ORCASTACK_SW_AUTOMATION_IDENTITY", app.DefaultSWAutoIdentity))
	if err != nil {
		log.Fatal(err)
	}
}
