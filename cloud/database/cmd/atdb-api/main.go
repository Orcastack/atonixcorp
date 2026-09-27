package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"atonixcorp/cloud/database/api"
	"atonixcorp/cloud/database/core"
	"atonixcorp/cloud/database/openstack"
	"atonixcorp/cloud/database/store"

	"gopkg.in/yaml.v3"
)

// Config structs
type EnvConfig struct {
	Environment string `yaml:"environment"`
	Region      string `yaml:"region"`
}

type OpenStackConfig struct {
	Auth struct {
		AuthURL   string `yaml:"auth_url"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Domain    string `yaml:"domain"`
		ProjectID string `yaml:"project_id"`
	} `yaml:"auth"`

	Endpoints struct {
		Compute string `yaml:"compute"`
		Network string `yaml:"network"`
		Volume  string `yaml:"volume"`
	} `yaml:"endpoints"`
}

type PlansConfig struct {
	Plans map[string]core.Plan `yaml:"plans"`
}

// Load YAML helper
func loadYAML(path string, out interface{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", path, err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		log.Fatalf("Failed to parse %s: %v", path, err)
	}
}

func main() {
	fmt.Println("AtonixCorp DBaaS API — Starting...")

	// Load configs
	var envCfg EnvConfig
	var osCfg OpenStackConfig
	var plansCfg PlansConfig

	loadYAML("atonixcorp/cloud/database/config/env.yaml", &envCfg)
	loadYAML("atonixcorp/cloud/database/config/openstack.yaml", &osCfg)
	loadYAML("atonixcorp/cloud/database/config/plans.yaml", &plansCfg)

	fmt.Println("Environment:", envCfg.Environment)
	fmt.Println("Region:", envCfg.Region)

	// Metadata store (Postgres)
	dsn := os.Getenv("ATONIXCORP_DB_DSN")
	if dsn == "" {
		log.Fatal("ATONIXCORP_DB_DSN environment variable is required")
	}

	store, err := store.NewPostgresStore(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to metadata DB: %v", err)
	}
	defer store.Close()

	// OpenStack client
	osClient, err := openstack.NewClient(openstack.AuthConfig{
		AuthURL:   osCfg.Auth.AuthURL,
		Username:  osCfg.Auth.Username,
		Password:  osCfg.Auth.Password,
		Domain:    osCfg.Auth.Domain,
		ProjectID: osCfg.Auth.ProjectID,
	})
	if err != nil {
		log.Fatalf("Failed to authenticate with OpenStack: %v", err)
	}

	openStackService := openstack.NewServiceClient(osClient)

	// Controller
	controller := core.NewController(store, openStackService)
	controller.SetPlans(plansCfg.Plans)

	// API server
	apiServer := api.NewAPIServer(controller, store, openStackService)
	router := apiServer.Router()

	// Start HTTP server
	addr := ":8080"
	fmt.Println("AtonixCorp DBaaS API running on", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("API server failed: %v", err)
	}
}
