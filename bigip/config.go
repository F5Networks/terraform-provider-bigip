package bigip

import (
	"fmt"
	"log"

	"github.com/DealerDotCom/go-bigip"
)

type Config struct {
	Address  string
	Username string
	Password string
}

func (c *Config) Client() (*bigip.BigIP, error) {

	if c.Address != "" && c.Username != "" && c.Password != "" {
		log.Println("[INFO] Initializing BigIP connection")
		client := bigip.NewSession(c.Address, c.Username, c.Password)
		err := c.validateConnection(client)
		if err == nil {
			return client, nil
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("BigIP provider requires address, username and password")
	}
}

func (c *Config) validateConnection(client *bigip.BigIP) error {
	_, err := client.SelfIPs()
	if err != nil {
		return err
	}
	return nil
}
