/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/x509"
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"log"
	"os"
)

func Client(config *bigip.Config) (*bigip.BigIP, error) {

	log.Println("[INFO] Initializing BigIP connection")
	var client *bigip.BigIP
	var err error
	// If we have a token value, we do not want to authenticate using a
	// Token Session. The user has already authenticated with the BigIP
	// outside of the provider, so even if the BigIP is using Token Auth,
	// we don't want to do that here. We want to use bigip.NewSession
	if config.LoginReference != "" && config.Token == "" && config.Address != "" {
		client, err = bigip.NewTokenSession(config)
		// client, err = bigip.NewTokenSession(c)
		if err != nil {
			log.Printf("[ERROR] Error creating New Token Session %s ", err)
			return nil, err
		}

	} else {
		client = bigip.NewSession(config)
		// The provider will use the Token value instead of the password
		if config.Token != "" {
			client.Token = config.Token
		}
	}
	if config.Address != "" && config.Username != "" && config.Password != "" {
		log.Printf("client connection before :%+v", client.Transport.Proxy)
		client.Transport.TLSClientConfig.InsecureSkipVerify = config.CertVerifyDisable
		//httpProxy, httpProxyOK := os.LookupEnv("HTTP_PROXY")
		//httpsProxy, httpsProxyOK := os.LookupEnv("HTTPS_PROXY")
		//var proxy any
		//if httpProxyOK {
		//	proxy, _ = url.Parse(httpProxy)
		//}
		//if httpsProxyOK {
		//	proxy, _ = url.Parse(httpsProxy)
		//}
		////log.Println("HTTP PROXY from Env:", httpproxy.FromEnvironment())
		//if proxy != nil {
		//	client.Transport.Proxy = http.ProxyURL(proxy.(*url.URL))
		//}
		//log.Printf("client connection after proxy:%+v", client)
		//time.Sleep(10 * time.Second)
		if !config.CertVerifyDisable {
			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}
			certPEM, err := os.ReadFile(config.TrustedCertificate)
			if err != nil {
				return nil, fmt.Errorf("provide Valid Trusted certificate path :%+v", err)
				// log.Printf("[DEBUG]read cert PEM/crt file error:%+v", err)
			}
			// TODO: Make sure appMgr sets certificates in bigipInfo
			// certs := certPEM)

			// Append our certs to the system pool
			if ok := rootCAs.AppendCertsFromPEM(certPEM); !ok {
				log.Println("[DEBUG] No certs appended, using only system certs")
			}
			client.Transport.TLSClientConfig.RootCAs = rootCAs
		}
		err = client.ValidateConnection()
		if err == nil {
			return client, nil
		}
	}
	return client, err

}
