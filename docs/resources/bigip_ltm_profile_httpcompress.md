---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_httpcompress"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_httpcompress resource
---

# bigip\_ltm\_profile_httpcompress

`bigip_ltm_profile_httpcompress`  Virtual server HTTP compression profile configuration

Resources should be named with their `full path`.The full path is the combination of the `partition + name` (example: `/Common/my-httpcompresprofile` ) or  `partition + directory + name` of the resource  (example: `/Common/test/my-httpcompresprofile`)

## Example Usage


```hcl
resource "bigip_ltm_profile_httpcompress" "sjhttpcompression" {
  name                 = "/Common/sjhttpcompression2"
  defaults_from        = "/Common/httpcompression"
  uri_exclude          = ["www.abc.f5.com", "www.abc2.f5.com"]
  uri_include          = ["www.xyzbc.cisco.com"]
  content_type_include = ["nicecontent.com"]
  content_type_exclude = ["nicecontentexclude.com"]
}

```      

## Argument Reference

* `name` (Required,type `string`) Name of the LTM http compress profile,named with their `full path`.The full path is the combination of the `partition + name` (example: `/Common/my-httpcompresprofile` ) or  `partition + directory + name` of the resource  (example: `my-httpcompresprofile`)

* `defaults_from` - (Optional,type `string`) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `uri_exclude`  - (Optional,type `set`) Disables compression on a specified list of HTTP Request-URI responses. Use a regular expression to specify a list of URIs you do not want to compress.

* `uri_include`  - (Optional,type `set`) Enables compression on a specified list of HTTP Request-URI responses. Use a regular expression to specify a list of URIs you want to compress.

* `content_type_include` - (Optional,type `set`) Specifies a list of content types for compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to compress.

* `content_type_exclude` - (Optional,type `set`) Excludes a specified list of content types from compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to compress.

* `compression_buffersize` - (Optional,type `int`) Specifies the maximum number of compressed bytes that the system buffers before inserting a Content-Length header (which specifies the compressed size) into the response. The default is `4096` bytes.

* `gzip_compression_level` - (Optional,type `int`) Specifies the degree to which the system compresses the content. Higher compression levels cause the compression process to be slower. The default is 1 - Least Compression (Fastest)

* `gzip_memory_level` - (Optional,type `int`) Specifies the number of bytes of memory that the system uses for internal compression buffers when compressing a server response. The default is `8 kilobytes/8192 bytes`.

* `gzip_window_size` - (Optional,type `int`)  Specifies the number of kilobytes in the window size that the system uses when compressing a server response. The default is `16` kilobytes

* `keep_accept_encoding` - (Optional,type `string`) Specifies, when checked (enabled), that the system does not remove the Accept-Encoding: header from an HTTP request. The default is `disabled`.

* `vary_header` - (Optional,type `string`) Specifies, when checked (enabled), that the system inserts a Vary header into cacheable server responses. The default is `enabled`.

* `cpu_saver` - (Optional,type `string`) Specifies, when checked (enabled), that the system monitors the percent CPU usage and adjusts compression rates automatically when the CPU usage reaches either the CPU Saver High Threshold or the CPU Saver Low Threshold. The default is `enabled`.


## Import

BIG-IP LTM HTTP Compress profiles can be imported using the `name`, e.g.

```
$ terraform import bigip_ltm_profile_httpcompress.test-httpcomprs_import /Common/test-httpcomprs
```
