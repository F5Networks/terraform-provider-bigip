package main

import (
	"github.com/f5devcentral/terraform-provider-f5/bigip"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bigip.Provider})
}
