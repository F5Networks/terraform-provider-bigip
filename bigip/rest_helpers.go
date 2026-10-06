/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

// Shared helpers for resources that talk to iControl REST endpoints
// go-bigip does not wrap (auth/ldap, auth/radius, auth/radius-server,
// auth/tacacs, auth/user) or wraps with a bug (sys/syslog: Syslog's
// MarshalJSON builds an empty, never-populated syslogDTO, so
// client.CreateSyslog/ModifySyslog always PUT/PATCH "{}"). These resources
// call client.APICall directly instead, matching the precedent in
// resource_bigip_ssl_key_cert.go.

import (
	"encoding/json"
	"fmt"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
)

// mangleFullPath converts a BIG-IP full path such as "/Common/system-auth"
// into the mangled form ("~Common~system-auth") iControl REST uses in URLs.
// Names with no partition (e.g. plain usernames like "admin") are returned
// unchanged.
func mangleFullPath(fullPath string) string {
	if !strings.Contains(fullPath, "/") {
		return fullPath
	}
	return strings.ReplaceAll(fullPath, "/", "~")
}

// stripPartitionPrefix converts a BIG-IP full path such as
// "/Common/test-radius-server" back into the bare object name
// ("test-radius-server"), the form some resources (e.g.
// bigip_auth_radius_server) use for both their own id and name. BIG-IP
// consistently echoes referenced-by-name fields (auth/radius's servers,
// sys/syslog's remoteServers[].name, ...) back as a full path even when the
// request body supplied a bare name, so Read functions that want a stable,
// round-trippable value need to undo that. Values that aren't a simple
// "/partition/name" full path (no leading slash, or more than two
// segments) are returned unchanged.
func stripPartitionPrefix(s string) string {
	if !strings.HasPrefix(s, "/") {
		return s
	}
	parts := strings.Split(s, "/")
	if len(parts) == 3 && parts[0] == "" {
		return parts[2]
	}
	return s
}

// restIsNotFound reports whether err (as returned by go-bigip's APICall) is
// the "object does not exist" error for an HTTP 404 response. go-bigip does
// not expose the underlying HTTP status on error -- only the BIG-IP JSON
// error body's message text -- but every observed 404 response body
// (auth/ldap, auth/radius, auth/tacacs, auth/user, sys/syslog, ...) contains
// "not found", so this is a best-effort text match. See also
// bigip/testing_helpers_test.go's IsNotFoundError, the test-only twin of
// this function.
func restIsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}

// restGet issues a raw GET against the given iControl REST path (e.g.
// "auth/ldap/~Common~system-auth") and, if the object exists, unmarshals
// the JSON response into out. It returns found=false, err=nil (rather than
// propagating the error) when the object does not exist, so callers can
// implement the standard Terraform "removed out of band -> clear ID"
// behavior in their Read function without special-casing 404 text
// themselves.
func restGet(client *bigip.BigIP, path string, out interface{}) (bool, error) {
	req := &bigip.APIRequest{
		Method:      "get",
		URL:         "mgmt/tm/" + path,
		ContentType: "application/json",
	}
	resp, err := client.APICall(req)
	if err != nil {
		if restIsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	if out != nil {
		if err := json.Unmarshal(resp, out); err != nil {
			return false, fmt.Errorf("error decoding response from %s: %w", path, err)
		}
	}
	return true, nil
}

// restSend marshals body (if non-nil) and issues an HTTP request of the
// given method against path.
func restSend(client *bigip.BigIP, method, path string, body interface{}) error {
	var bodyStr string
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyStr = string(b)
	}
	req := &bigip.APIRequest{
		Method:      method,
		URL:         "mgmt/tm/" + path,
		Body:        bodyStr,
		ContentType: "application/json",
	}
	_, err := client.APICall(req)
	return err
}

func restPost(client *bigip.BigIP, path string, body interface{}) error {
	return restSend(client, "post", path, body)
}

func restPut(client *bigip.BigIP, path string, body interface{}) error {
	return restSend(client, "put", path, body)
}

func restPatch(client *bigip.BigIP, path string, body interface{}) error {
	return restSend(client, "patch", path, body)
}

// restDelete issues a DELETE against path, treating "already gone" (404) as
// success.
func restDelete(client *bigip.BigIP, path string) error {
	req := &bigip.APIRequest{
		Method: "delete",
		URL:    "mgmt/tm/" + path,
	}
	_, err := client.APICall(req)
	if err != nil && !restIsNotFound(err) {
		return err
	}
	return nil
}
