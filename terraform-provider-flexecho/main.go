package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/silk-us/terraform-provider-flexecho/flexecho"
)

// provider entrypoint. terraform launches this binary and talks grpc to it
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return flexecho.Provider()
		},
	})
}
