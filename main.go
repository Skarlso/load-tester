package main

import (
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create kubernetes client
	k8sClient, err := NewK8sClient()
	if err != nil {
		log.Fatalf("Failed to create kubernetes client: %v", err)
	}

	// Run load test
	ctx := context.Background()
	err = k8sClient.RunLoadTest(ctx, config)
	if err != nil {
		log.Fatalf("Load test failed: %v", err)
	}

	fmt.Println("Load test completed successfully!")
}
