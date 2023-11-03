---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_web_acceleration"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_web_acceleration resource
---

# bigip\_ltm\_profile_web_acceleration

`bigip_ltm_profile_web_acceleration` Configures a custom web-acceleration profile for use.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/sample-resource.

## Example Usage


```hcl
resource "bigip_ltm_profile_web_acceleration" "sample-resource" {
  name              = "/Common/sample-resource"
  defaults_from     = "/Common/test2"
  cache_size        = 101
  cache_max_entries = 201
}
```      

## Argument Reference

* `name` (Required,type `string`) Specifies the name of the web acceleration profile service ,name of Profile should be full path. Full path is the combination of the `partition + web acceleration profile name`,For example `/Common/sample-resource`.

* `defaults_from` - (optional,type `string`) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `cache_size` - (optional,type `int`) 	Specifies the maximum size for the cache. When the cache reaches the maximum size, the system starts removing the oldest entries. The default value is `100 megabytes`.

* `cache_max_entries` - (Optional, type `int`) Specifies the maximum number of entries that can be in the cache. The default value is `0` (zero), which means that the system does not limit the maximum entries.

* `cache_max_age` - (Optional, type `int`) Specifies how long the system considers the cached content to be valid. The default value is `3600 seconds`.

* `cache_object_min_size` - (Optional,type `int`) Specifies the smallest object that the system considers eligible for caching. The default value is `500 bytes`.

* `cache_object_max_size` - (Optional, type `int`) Specifies the smallest object that the system considers eligible for caching. The default value is `500 bytes`.

* `cache_uri_exclude` - (Optional,type `list`) Configures a list of URIs to exclude from the cache. The default value of `none` specifies no URIs are excluded.

* `cache_uri_include` - (Optional,type `list`) Configures a list of URIs to include in the cache. The default value of `.*` specifies that all URIs are cacheable.

* `cache_uri_include_override` - (Optional,type `list`) Configures a list of URIs to include in the cache even if they would normally be excluded due to factors like object size or HTTP request type. The default value of none specifies no URIs are to be forced into the cache.

* `cache_uri_pinned` - (Optional,type `list`) Configures a list of URIs to keep in the cache. The pinning process keeps URIs in cache when they would normally be evicted to make room for more active URIs.

* `cache_client_cache_control_mode` - (Optional, type `string`) Specifies which cache disabling headers sent by clients the system ignores. The default value is `all`.

* `cache_insert_age_header` - (Optional, type `string`) Inserts Age and Date headers in the response. The default value is `enabled`.

* `cache_aging_rate` - (Optional,type `int`) Specifies how quickly the system ages a cache entry. The aging rate ranges from 0 (slowest aging) to 10 (fastest aging). The default value is `9`.
