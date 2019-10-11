[//]: # (Original work from https://github.com/DealerDotCom/terraform-provider-bigip)
[//]: # (Modifications Copyright 2019 F5 Networks Inc.)
[//]: # (This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.)
[//]: # (If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.)

# Overview

A [Terraform](terraform.io) provider for F5 BigIP LTM.

[![Build Status](https://travis-ci.org/f5devcentral/terraform-provider-bigip.svg?branch=master)](https://travis-ci.org/f5devcentral/terraform-provider-bigip)
[![Go Report Card](https://goreportcard.com/badge/github.com/f5devcentral/terraform-provider-bigip)](https://goreportcard.com/report/github.com/f5devcentral/terraform-provider-bigip)
[![license](https://img.shields.io/badge/license-Mozilla-red.svg?style=flat)](https://github.com/f5devcentral/terraform-provider-bigip/blob/master/LICENSE)
[![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
# Support
- F5 provides support for the F5 Modules for Terraform. For more information, see this page.(https://www.f5.com/services/support/support-offerings/support-policies)
- If you’d like community support, you can open an issue on GitHub.(https://github.com/F5Networks/terraform-provider-bigip)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

# Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x / 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

# F5 BigIP LTM requirements

- This provider uses the iControlREST API, make sure that it is installed and enabled on your F5 device before proceeding.

These BIG-IP versions are supported in these Terraform versions.

| BIG-IP version	|Terraform 0.12 |	Terraform 0.11  |
|-----------------|---------------|-----------------|
| BIG-IP 14.x	    | 	   X        |       X         |
| BIG-IP 12.x	    |      X        |      	X         |
| BIG-IP 13.x	    |      X        |       X         |


# Documentation

Provider documentation and reference can be found [here](website/docs).

# Building the  Provider

Clone repository to: $GOPATH/src/github.com/terraform-providers/terraform-provider-bigip

```
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:F5Networks/terraform-provider-bigip

```
Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-bigip
$ make build

```
# Using the Provider

If you're building the provider, follow the instructions to install it as a plugin. After placing it into your plugins directory, run terraform init to initialize it.

# Developing the Provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.11+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-bigip
...

```
# Testing

Running the acceptance test suite requires an F5 to test against. Set `BIGIP_HOST`, `BIGIP_USER`
and `BIGIP_PASSWORD` to a device to run the tests against. By default tests will use the `Common`
partition for creating objects. You can change the partition by setting `BIGIP_TEST_PARTITION`.

```
BIGIP_HOST=f5.mycompany.com BIGIP_USER=foo BIGIP_PASSWORD=secret make testacc
```


Read [here](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test) for
more information about acceptance testing in Terraform.
