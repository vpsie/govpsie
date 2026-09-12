// Command basic is a runnable example of the govpsie SDK.
//
// It reads a VPSie API token from the VPSIE_ACCESS_TOKEN environment variable,
// lists the account's datacenters and servers, and prints them.
//
// Usage:
//
//	export VPSIE_ACCESS_TOKEN="your-api-token"
//	go run ./examples/basic
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vpsie/govpsie"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("VPSIE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("set VPSIE_ACCESS_TOKEN to run this example")
	}

	client := govpsie.NewClient(oauth2.NewClient(context.Background(), nil))
	client.SetRequestHeaders(map[string]string{"Vpsie-Auth": token})

	ctx := context.Background()

	datacenters, err := client.DataCenter.List(ctx, &govpsie.ListOptions{})
	if err != nil {
		log.Fatalf("list datacenters: %v", err)
	}

	fmt.Printf("Datacenters (%d):\n", len(datacenters))
	for _, dc := range datacenters {
		fmt.Printf("  - %s, %s (%s)\n", dc.DcName, dc.Country, dc.Identifier)
	}

	servers, err := client.Server.List(ctx, &govpsie.ListOptions{})
	if err != nil {
		log.Fatalf("list servers: %v", err)
	}

	fmt.Printf("\nServers (%d):\n", len(servers))
	for _, s := range servers {
		fmt.Printf("  - %s (%s)\n", s.Hostname, s.Identifier)
	}
}
