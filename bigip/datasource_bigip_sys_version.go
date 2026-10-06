/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// sysVersionResponse mirrors the shape of a GET /mgmt/tm/sys/version
// response. Unlike most tm/ objects, sys/version is a "stats"-style
// resource: its single entry is keyed by a selfLink URL (not a fixed name)
// and each field is nested under nestedStats.entries.<Field>.description
// rather than being a flat top-level property. go-bigip has no typed
// binding for this endpoint (its Version struct is for the unrelated
// cli/version object), so this is decoded directly, matching the same
// nestedStats-parsing precedent in datasource_bigip_net_interfaces.go.
type sysVersionResponse struct {
	Entries map[string]struct {
		NestedStats struct {
			Entries map[string]struct {
				Description string `json:"description,omitempty"`
			} `json:"entries"`
		} `json:"nestedStats"`
	} `json:"entries"`
}

// firstEntry returns the nested field map of the single entry in a
// sysVersionResponse-shaped response (sys/version has exactly one entry,
// keyed by an unpredictable selfLink). Returns nil if the response had no
// entries at all.
func (r sysVersionResponse) firstEntry() map[string]struct {
	Description string `json:"description,omitempty"`
} {
	for _, entry := range r.Entries {
		return entry.NestedStats.Entries
	}
	return nil
}

// sysHardwarePlatformResponse mirrors the shape of a GET
// /mgmt/tm/sys/hardware response, as far as the top-level "platform"
// sub-object is concerned. Unlike sys/version's single level of nesting,
// sys/hardware nests one level deeper: the outer entries map is keyed by
// section (e.g. "sys/hardware/platform"), and each section's own
// nestedStats.entries is itself another map (keyed by an index like
// "sys/hardware/platform/0") whose nestedStats.entries finally holds the
// actual fields (marketingName, baseMac, etc).
type sysHardwarePlatformResponse struct {
	Entries map[string]struct {
		NestedStats struct {
			Entries map[string]struct {
				NestedStats struct {
					Entries map[string]struct {
						Description string `json:"description,omitempty"`
					} `json:"entries"`
				} `json:"nestedStats"`
			} `json:"entries"`
		} `json:"nestedStats"`
	} `json:"entries"`
}

// platformMarketingName retrieves the platform's marketing name (e.g.
// "BIG-IP Virtual Edition", "BIG-IP r10900") from GET /mgmt/tm/sys/hardware.
// sys/version itself has no platform field; sys/hardware's "platform"
// nested object is the closest equivalent (see
// scripts/inventory-tmos-version.sh, which parses the same information out
// of "tmsh show sys hardware" text output for the same reason).
//
// Returns ("", reason) rather than failing outright when sys/hardware
// doesn't return anything usable, since platform is supplemental
// information layered on top of the primary sys/version read this data
// source's acceptance criteria is centered on -- the caller surfaces
// "reason" as a warning diagnostic rather than a hard read failure.
//
// Every level of the nested map is checked explicitly with the two-value
// map-index form (ok) rather than direct field/index chaining. Note that
// this is a defensive-programming/clarity improvement, not a crash fix:
// reading a missing key from a Go map (nil or otherwise) returns the zero
// value rather than panicking, so the previous direct-chaining form could
// not have panicked even if the API response omitted a level.
func platformMarketingName(client *bigip.BigIP) (name string, reason string) {
	var hw sysHardwarePlatformResponse
	found, err := restGet(client, "sys/hardware", &hw)
	if err != nil {
		return "", fmt.Sprintf("error retrieving sys/hardware: %v", err)
	}
	if !found {
		return "", "sys/hardware not found"
	}
	// BIG-IP's iControl REST stats-style responses always report
	// "https://localhost" as the selfLink host for these per-section keys,
	// regardless of the actual address/hostname the client connected
	// through -- this is a verified BIG-IP behavior (also relied on by the
	// nestedStats key literals in the unit tests below), not a guess tied
	// to any particular device.
	platform, ok := hw.Entries["https://localhost/mgmt/tm/sys/hardware/platform"]
	if !ok {
		return "", "sys/hardware response had no \"platform\" entry"
	}
	if platform.NestedStats.Entries == nil {
		return "", "sys/hardware \"platform\" entry had no nested entries"
	}
	for _, sub := range platform.NestedStats.Entries {
		if sub.NestedStats.Entries == nil {
			return "", "sys/hardware \"platform\" sub-entry had no nested fields"
		}
		marketingName, ok := sub.NestedStats.Entries["marketingName"]
		if !ok {
			return "", "sys/hardware \"platform\" entry had no \"marketingName\" field"
		}
		return marketingName.Description, ""
	}
	return "", "sys/hardware \"platform\" entry had an empty nested-entries map"
}

// dataSourceBigipSysVersion reads the BIG-IP system's running TMOS version
// and platform information (sys/version, plus platform from sys/hardware)
// -- both singleton objects, so this data source takes no arguments. This
// is useful for migration tooling that needs to determine the currently
// running version/build in order to select a matching r-Series/VELOS
// tenant image (see scripts/inventory-tmos-version.sh, which derives the
// same recommended tenant image filename from a bigip_command-based
// tmsh parse instead, since this data source didn't previously exist).
func dataSourceBigipSysVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipSysVersionRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Running TMOS version, e.g. \"17.5.1\".",
			},
			"build": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Build number of the running TMOS version, e.g. \"0.0.7\".",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product name, e.g. \"BIG-IP\".",
			},
			"edition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Release edition, e.g. \"Final\", \"Hotfix\", \"Engineering Hotfix\".",
			},
			"platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Platform marketing name (from sys/hardware), e.g. \"BIG-IP Virtual Edition\" or a specific appliance model. Empty if sys/hardware's platform information could not be retrieved.",
			},
		},
	}
}

func dataSourceBigipSysVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	log.Println("[INFO] Reading system version")

	var version sysVersionResponse
	found, err := restGet(client, "sys/version", &version)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving system version: %v", err))
	}
	if !found {
		return diag.FromErr(fmt.Errorf("system version not found"))
	}

	fields := version.firstEntry()
	if fields == nil {
		return diag.FromErr(fmt.Errorf("system version response had no entries"))
	}

	_ = d.Set("version", fields["Version"].Description)
	_ = d.Set("build", fields["Build"].Description)
	_ = d.Set("product", fields["Product"].Description)
	_ = d.Set("edition", fields["Edition"].Description)

	var diags diag.Diagnostics
	platform, reason := platformMarketingName(client)
	_ = d.Set("platform", platform)
	if platform == "" && reason != "" {
		log.Printf("[WARN] Unable to retrieve platform info from sys/hardware: %s", reason)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unable to retrieve platform info",
			Detail:   fmt.Sprintf("platform will be empty: %s", reason),
		})
	}

	d.SetId(hashForState(fmt.Sprintf("%v", fields)))

	return diags
}
