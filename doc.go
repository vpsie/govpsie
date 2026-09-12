// Package govpsie is the Go client library for the VPSie cloud API.
//
// It provides typed access to VPSie resources — servers, storage, networking,
// DNS, load balancers, Kubernetes, container registries, managed databases and
// more — through a single Client.
//
// # Getting started
//
// Create a Client from an *http.Client and set the authentication header:
//
//	client := govpsie.NewClient(oauth2.NewClient(context.Background(), nil))
//	client.SetRequestHeaders(map[string]string{
//		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
//	})
//
//	servers, err := client.Server.List(context.Background(), &govpsie.ListOptions{})
//
// # Services
//
// Each API area is exposed as a service field on Client (for example
// client.Server, client.Storage, client.Kubernetes, client.ManagedDB). Every
// request flows through Client.NewRequest and Client.Do, which unwraps the
// standard response envelope {error, data, total} and converts a non-2xx status
// into an error carrying the API message.
//
// # Authentication
//
// The API authenticates with a Vpsie-Auth header. Obtain a token from the VPSie
// console or via POST /auth/from/api with your clientId and clientSecret.
package govpsie
