/*
Package main prints the current version of the package.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Manifest struct {
	Version string `json:"."`
}

func main() {
	handle, err := os.Open(".release-please-manifest.json")
	if err != nil {
		log.Printf("failed to open manifest - %v", err)
		os.Exit(1)
	}
	defer func() {
		err := handle.Close()
		if err != nil {
			log.Printf("failed to close manifest - %v", err)
		}
	}()

	contents, err := io.ReadAll(handle)
	if err != nil {
		log.Printf("failed to rad manifest - %v", err)
		os.Exit(1)
	}

	var manifest Manifest
	err = json.Unmarshal(contents, &manifest)
	if err != nil {
		log.Printf("failed to umarshal manifest - %v", err)
		os.Exit(1)
	}

	fmt.Printf("v%s", manifest.Version)
}
