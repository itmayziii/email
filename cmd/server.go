package main

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/testdata/conformance/function"
	"log"
	"os"
)

import (
	// Blank-import the function package so the init() runs
	_ "github.com/itmayziii/email-api"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterCloudEventFunctionContext(ctx, "/", function.CloudEvent); err != nil {
		log.Fatalf("funcframework.RegisterCloudEventFunctionContext: %v\n", err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
