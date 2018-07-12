# Overview

A [Terraform](terraform.io) provider for F5 BigIP LTM.

[![Build Status](https://travis-ci.org/f5devcentral/terraform-provider-bigip.svg?branch=master)](https://travis-ci.org/f5devcentral/terraform-provider-bigip)
[![Go Report Card](https://goreportcard.com/badge/github.com/f5devcentral/terraform-provider-bigip)](https://goreportcard.com/report/github.com/f5devcentral/terraform-provider-bigip)
[![license](https://img.shields.io/badge/license-Mozilla-red.svg?style=flat)](https://github.com/f5devcentral/terraform-provider-bigip/blob/master/LICENSE)
[![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

# Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

# F5 BigIP LTM requirements

- This provider uses the iControlREST API, make sure that is installed and enabled on your F5 before proceeding.
- All the resources are validated with BigIP v12.1.1

# Dcoumentation

Provider documentation and reference can be found [here](website/docs).

# Quick Start with BIG-IP Provider

Install appropriate Go package from  https://golang.org/dl/ make sure your go version is 1.9 & above
```
go version
go version go1.9.2 darwin/amd64

mkdir workspace
export GOPATH=$HOME/workspace
mkdir -p $GOPATH/src/github.com/f5devcentral
cd $GOPATH
go get github.com/f5devcentral/terraform-provider-bigip
cd src/github.com/f5devcentral/terraform-provider-bigip/
go build
create .tf
terraform init
Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

```
# Building

Create the distributable packages like so:

```
make get-deps && make bin && make dist
```

See these pages for more information:

 * https://www.terraform.io/docs/internals/internal-plugins.html
 * https://github.com/hashicorp/terraform#developing-terraform

# Testing

Running the acceptance test suite requires an F5 to test against. Set `BIGIP_HOST`, `BIGIP_USER`
and `BIGIP_PASSWORD` to a device to run the tests against. By default tests will use the `Common`
partition for creating objects. You can change the partition by setting `BIGIP_TEST_PARTITION`.

```
BIGIP_HOST=f5.mycompany.com BIGIP_USER=foo BIGIP_PASSWORD=secret make testacc
```


Read [here](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test) for
more information about acceptance testing in Terraform.
