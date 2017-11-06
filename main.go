package main

import (
	"github.com/f5devcentral/terraform-provider-bigip/bigip"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bigip.Provider})
}
