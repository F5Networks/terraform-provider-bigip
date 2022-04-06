/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"log"

	bigip "github.com/f5devcentral/go-bigip"
)

type Config struct {
	Address        string
	Port           string
	Username       string
	Password       string
	Token          string
	LoginReference string `json:"loginProviderName"`
	ConfigOptions  *bigip.ConfigOptions
}

func (c *Config) Client() (*bigip.BigIP, error) {

	log.Println("[INFO] Initializing BigIP connection")
	var client *bigip.BigIP
	var err error
	// If we have a token value, we do not want to authenticate using a
	// Token Session. The user has already authenticated with the BigIP
	// outside of the provider, so even if the BigIP is using Token Auth,
	// we don't want to do that here. We want to use bigip.NewSession
	if c.LoginReference != "" && c.Token == "" {
		client, err = bigip.NewTokenSession(c.Address, c.Port, c.Username, c.Password, c.LoginReference, c.ConfigOptions)
		if err != nil {
			log.Printf("[ERROR] Error creating New Token Session %s ", err)
			return nil, err
		}

	} else {
		client = bigip.NewSession(c.Address, c.Port, c.Username, c.Password, c.ConfigOptions)
		// The provider will use the Token value instead of the password
		if c.Token != "" {
			client.Token = c.Token
		}
	}
	if c.Address != "" && c.Username != "" && c.Password != "" {
		err = c.validateConnection(client)
		if err == nil {
			return client, nil
		}
	}
	return nil, err

}

func (c *Config) validateConnection(client *bigip.BigIP) error {
	t, err := client.SelfIPs()
	if err != nil {
		log.Printf("[ERROR] Connection to BigIP device could not have been validated: %v ", err)
		return err
	}

	if t == nil {
		log.Printf("[WARN] Could not validate connection to BigIP")
		return nil
	}
	return nil
}
