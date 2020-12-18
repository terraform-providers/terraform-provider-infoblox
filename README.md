<!-- archived-provider -->
Please note: This Terraform provider is archived, per our [provider archiving process](https://terraform.io/docs/internals/archiving.html). What does this mean?
1. The code repository and all commit history will still be available.
1. Existing released binaries will remain available on the releases site.
1. Issues and pull requests are not being monitored.
1. New releases will not be published.

If anyone from the community or an interested third party is willing to maintain it, they can fork the repository and [publish it](https://www.terraform.io/docs/registry/providers/publishing.html) to the Terraform Registry. If you are interested in maintaining this provider, please reach out to the [Terraform Provider Development Program](https://www.terraform.io/guides/terraform-provider-development-program.html) at *terraform-provider-dev@hashicorp.com*.

# Terraform Provider for Infoblox
 <img width="171" alt="capture" src="https://user-images.githubusercontent.com/36291746/39614422-6b653088-4f8d-11e8-83fd-05b18ca974a2.PNG">

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.11.x or greater
* [Go](https://golang.org/doc/install) 1.12.x (to build the provider plugin)
* CNA License need to be installed on NIOS. If CNA is not installed then following default EA's should be added in NIOS side:
   * VM Name :: String Type
   * VM ID :: String Type
   * Tenant ID :: String Type
   * CMP Type :: String Type
   * Cloud API Owned :: List Type (Values True, False)
   * Network Name :: String Type

## Building the Provider
Clone repository to `$GOPATH/src/github.com/terraform-providers/terraform-provider-infoblox`.
```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-infoblox
```
Enter the provider directory and build the provider.
```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-infoblox
$ export GO111MODULE="on"
$ make build
```

## Using the Provider
If you're building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After the build is complete, copy the `terraform-provider-infoblox` binary into the same path as your terraform binary. After placing it into your plugins directory, run `terraform init` to initialize it.

## Developing the Provider
If you wish to work on the provider, you'll first need Go installed on your machine (version 1.9+ is required). You'll also need to correctly setup a `GOPATH`, as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.
```sh
$ export GO111MODULE="on"
$ make build
...
$ $GOPATH/bin/terraform-provider-infoblox
...
```
To test the provider, you can simply run `make test`.
```sh
$ make test
```

In order to run the full suite of acceptance tests `make testacc`.
```sh
$ make testacc
```
## Features of Provider
### Resource
* Creation of Network View in NIOS appliance
* Creation & Deletion of Network in NIOS appliance
* Allocation & Deallocation of IP from a Network
* Association & Disassociation of IP Address for a VM
* Creation and Deletion of A, CNAME, Host, Zones and Ptr records

### Data Source
* Supports Data Source for Network

## Disclaimer
To use the provider for DNS purposes, a parent (i.e. zone) must already exist. The plugin does not support the creation of zones.
while running acceptance tests create a 10.0.0.0/24 network under default network view and create a reservation for 10.0.0.2 IP
