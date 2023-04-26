/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

resource "bigip_ltm_pool" "k8s_prod" {
  name                = "/Common/test_pool_pa_tc1"
  load_balancing_mode = "round-robin"
  allow_snat          = "yes"
  allow_nat           = "yes"
}

resource "bigip_ltm_policy" "policy-issue-591" {
  name     = "/Common/policy-issue-591"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule-issue591"
    condition {
      http_host = true
      contains  = true
      values = [
        "domain1.net",
        "domain2.nl"
      ]
      request = true
    }
    condition {
      http_uri    = true
      path        = true
      not         = true
      starts_with = true
      values      = ["/role-service"]
      request     = true
    }
    action {
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.k8s_prod.name
    }
    action {
      forward    = false
      replace    = true
      connection = false
      http_uri   = true
      path       = "tcl:[string map {/role-service/ /} [HTTP::uri]]"
      request    = true
    }
  }
}

resource "bigip_ltm_policy" "policy_issue794_tc1" {
  name     = "/Common/policy_issue794_tc1"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule1"
    condition {
      case_insensitive = true
      http_uri         = true
      path             = true
      request          = true
      starts_with      = true
      values           = ["/wallet"]
    }
    action {
      http_host  = true
      replace    = true
      request    = true
      connection = false
      value      = "example.com"
    }
    action {
      forward    = false
      http_uri   = true
      replace    = true
      request    = true
      connection = false
      value      = "tcl:[string map {/wallet/ /wallet-dynamic} [HTTP::uri]]"
    }
    action {
      forward    = true
      request    = true
      connection = false
      pool       = bigip_ltm_pool.k8s_prod.name
    }
  }
}

resource "bigip_ltm_policy" "policy_issue794_tc2" {
  name     = "/Common/policy_issue794_tc2"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule1"
    condition {
      case_insensitive = true
      http_uri         = true
      path             = true
      request          = true
      starts_with      = true
      values           = ["/wallet"]
    }
    action {
      http_host  = true
      replace    = true
      request    = true
      connection = false
      value      = "example3.com"
    }
    action {
      forward    = false
      http_uri   = true
      replace    = true
      request    = true
      connection = false
      path       = "tcl:[string map {/wallet/ /wallet-dynamic} [HTTP::uri]]"
    }
    action {
      forward    = true
      request    = true
      connection = false
      pool       = bigip_ltm_pool.k8s_prod.name
    }
  }
}