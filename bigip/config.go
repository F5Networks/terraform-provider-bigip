package bigip

import (
	"fmt"
	"log"

	"github.com/scottdware/go-bigip"
)

type Config struct {
	Address        string
	Username       string
	Password       string
	LoginReference string
}

func (c *Config) Client() (*bigip.BigIP, error) {

	if c.Address != "" && c.Username != "" && c.Password != "" {
		log.Println("[INFO] Initializing BigIP connection")
		var client *bigip.BigIP
		var err error
		if c.LoginReference != "" {
			client, err = bigip.NewTokenSession(c.Address, c.Username, c.Password, c.LoginReference)
			if err != nil {
				return nil, err
			}
		} else {
			client = bigip.NewSession(c.Address, c.Username, c.Password)
		}
		err = c.validateConnection(client)
		if err == nil {
			return client, nil
		}
		return nil, err
	}
	return nil, fmt.Errorf("BigIP provider requires address, username and password")
}

func (c *Config) validateConnection(client *bigip.BigIP) error {
	_, err := client.SelfIPs()
	if err != nil {
		return err
	}
	return nil
}
