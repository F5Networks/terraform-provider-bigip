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

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

# Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x / 0.12.x /0.13.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

# F5 BigIP LTM requirements

- This provider uses the iControlREST API, make sure that it is installed and enabled on your F5 device before proceeding.

These BIG-IP versions are supported in these Terraform versions.

| BIG-IP version	|Terraform 1.x  |	Terraform 0.13  |	Terraform 0.12  | Terraform 0.11  |
|-----------------|---------------|-----------------|-----------------|-----------------|
| BIG-IP 16.x	    |      X        |       X         |       X         |      X          |
| BIG-IP 15.x	    |      X        |       X         |       X         |      X          |
| BIG-IP 14.x	    | 	   X        |       X         |       X         |      X          |
| BIG-IP 12.x	    |      X        |      	X         |       X         |      X          | 
| BIG-IP 13.x	    |      X        |       X         |       X         |      X          |


# Documentation

Documentation for the F5 BIG-IP terraform integration is available at https://clouddocs.f5.com/products/orchestration/terraform/latest/

Terraform provider documentation is available at  https://www.terraform.io/docs/providers/bigip/index.html

# Building the  Provider

Clone repository to: $GOPATH/src/github.com/F5Networks/terraform-provider-bigip

```
$ mkdir -p $GOPATH/src/github.com/F5Networks; cd $GOPATH/src/github.com/F5Networks
$ git clone https://github.com/F5Networks/terraform-provider-bigip.git

```
Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/F5Networks/terraform-provider-bigip
$ make build

```
# Using the Provider

If you're building the provider, follow the instructions to install it as a plugin. After placing it into your plugins directory, run terraform init to initialize it.

# Developing the Provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.11+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ GOFLAGS=-mod=vendor GO111MODULE=on make build
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

# Community Help

We encourage you to use our [Slack channel](https://f5cloudsolutions.herokuapp.com) for discussion and assistance on Terraform Resources (click the **terraform** channel). There are F5 employees who are members of this community who typically monitor the channel Monday-Friday 9-5 PST and will offer best-effort assistance. This slack channel community support should **not** be considered a substitute for F5 Technical Support. See the [Slack Channel Statement](slack-channel-statement.md) for guidelines on using this channel.

## Copyright

Copyright 2014-2020 F5 Networks Inc.

### F5 Networks Contributor License Agreement

[Terraform F5 Contributor License Agreement.pdf](F5_Contributor_License_Agreement.pdf)

Before you start contributing to any project sponsored by F5 Networks, Inc. (F5) on GitHub, you will need to sign a Contributor License Agreement (CLA).

If you are signing as an individual, we recommend that you talk to your employer (if applicable) before signing the CLA since some employment agreements may have restrictions on your contributions to other projects. Otherwise by submitting a CLA you represent that you are legally entitled to grant the licenses recited therein.

If your employer has rights to intellectual property that you create, such as your contributions, you represent that you have received permission to make contributions on behalf of that employer, that your employer has waived such rights for your contributions, or that your employer has executed a separate CLA with F5.

If you are signing on behalf of a company, you represent that you are legally entitled to grant the license recited therein. You represent further that each employee of the entity that submits contributions is authorized to submit such contributions on behalf of the entity pursuant to the CLA.


