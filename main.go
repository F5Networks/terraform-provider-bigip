package main

import (
	"github.com/DealerDotCom/terraform-provider-bigip/bigip"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bigip.Provider,
	})
}
