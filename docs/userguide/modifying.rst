Modifying Infrastructure
========================

Next you can modify the resource to see how Terraform handles change.

Terraform was built to help manage and enact change in environments where infrastructure is continuously evolving. As you change Terraform configurations, Terraform builds an execution plan that only modifies what is necessary to reach the desired state.

By using Terraform to change infrastructure, you can version control not only your configurations but also your state so you can see how the infrastructure evolves over time.

1. Modify the policy name of the resource. Edit the ``bigip_ltm_policy.test-policy`` resource in your configuration and change it to the following:

.. code-block:: javascript

   provider "bigip" {
    address = "x.x.x.x"
    username = "x.x.x.x"
    password = "x.x.x.x"
    }

    resource "bigip_ltm_policy" "test-policy" {
    name = "test_policy"
    strategy = "first-match"
    requires = ["http"]
    published_copy = "Drafts/test_policy"
    controls = ["forwarding"]
    rule {
    name = "rule6"
    action {
    tm_name = "20"
    forward = true
    pool = "/Common/mypool"
    }
    }
    depends_on = ["bigip_ltm_pool.mypool"]
    }

    resource "bigip_ltm_pool" "mypool" {
    name = "/Common/mypool"
    monitors = ["/Common/http"]
    allow_nat = "yes"
    allow_snat = "yes"
    load_balancing_mode = "round-robin"
    }


Terraform configurations are meant to be changed like this. You can also completely remove resources and Terraform will know to destroy the old one.


2. After changing the configuration, run ``terraform apply`` again to see how Terraform will apply this change to the existing resources. The prefix -/+ means that Terraform will destroy and recreate the resource, rather than update it in-place. While some attributes can be updated in-place (which are shown with the ~ prefix), Terraform handles these details for you, and the execution plan makes it clear what Terraform will do. 

.. code-block:: javascript

   root@terraforn-ubuntu3:~/go/src/github.com/terraform-providers/terraform-provider-bigip# terraform apply
    bigip_ltm_pool.mypool: Refreshing state... [id=/Common/mypool]
    bigip_ltm_policy.test-policy: Refreshing state... [id=new_policy]

    An execution plan has been generated and is shown below.
    Resource actions are indicated with the following symbols:
    -/+ destroy and then create replacement

    Terraform will perform the following actions:

    # bigip_ltm_policy.test-policy must be replaced
    -/+ resource "bigip_ltm_policy" "test-policy" {
    controls = [
    "forwarding",
    ]
    ~ id = "new_policy" -> (known after apply)
    ~ name = "new_policy" -> "test_policy" # forces replacement
    ~ published_copy = "Drafts/new_policy" -> "Drafts/test_policy" # forces replacement
    requires = [
    "http",
    ]
    ~ strategy = "/Common/first-match" -> "first-match"

    ~ rule {
    name = "rule6"

    ~ action {
    + app_service = (known after apply)
    + application = (known after apply)
    ~ asm = false -> (known after apply)
    ~ avr = false -> (known after apply)
    ~ cache = false -> (known after apply)
    ~ carp = false -> (known after apply)
    + category = (known after apply)
    ~ classify = false -> (known after apply)
    + clone_pool = (known after apply)
    ~ code = 0 -> (known after apply)
    ~ compress = false -> (known after apply)
    + content = (known after apply)
    ~ cookie_hash = false -> (known after apply)
    ~ cookie_insert = false -> (known after apply)
    ~ cookie_passive = false -> (known after apply)
    ~ cookie_rewrite = false -> (known after apply)
    ~ decompress = false -> (known after apply)
    ~ defer = false -> (known after apply)
    ~ destination_address = false -> (known after apply)
    ~ disable = false -> (known after apply)
    + domain = (known after apply)
    ~ enable = false -> (known after apply)
    + expiry = (known after apply)
    ~ expiry_secs = 0 -> (known after apply)
    + expression = (known after apply)
    + extension = (known after apply)
    + facility = (known after apply)
    forward = true
    + from_profile = (known after apply)
    ~ hash = false -> (known after apply)
    + host = (known after apply)
    ~ http = false -> (known after apply)
    ~ http_basic_auth = false -> (known after apply)
    ~ http_cookie = false -> (known after apply)
    ~ http_header = false -> (known after apply)
    - http_host = false -> null
    ~ http_referer = false -> (known after apply)
    ~ http_reply = false -> (known after apply)
    ~ http_set_cookie = false -> (known after apply)
    ~ http_uri = false -> (known after apply)
    + ifile = (known after apply)
    ~ insert = false -> (known after apply)
    + internal_virtual = (known after apply)
    + ip_address = (known after apply)
    + key = (known after apply)
    ~ l7dos = false -> (known after apply)
    ~ length = 0 -> (known after apply)
    + location = (known after apply)
    ~ log = false -> (known after apply)
    ~ ltm_policy = false -> (known after apply)
    + member = (known after apply)
    + message = (known after apply)
    + netmask = (known after apply)
    + nexthop = (known after apply)
    + node = (known after apply)
    ~ offset = 0 -> (known after apply)
    + path = (known after apply)
    ~ pem = false -> (known after apply)
    ~ persist = false -> (known after apply)
    ~ pin = false -> (known after apply)
    + policy = (known after apply)
    pool = "/Common/mypool"
    ~ port = 0 -> (known after apply)
    + priority = (known after apply)
    + profile = (known after apply)
    + protocol = (known after apply)
    + query_string = (known after apply)
    + rateclass = (known after apply)
    ~ redirect = false -> (known after apply)
    ~ remove = false -> (known after apply)
    ~ replace = false -> (known after apply)
    ~ request = false -> (known after apply)
    ~ request_adapt = false -> (known after apply)
    ~ reset = false -> (known after apply)
    ~ response = false -> (known after apply)
    ~ response_adapt = false -> (known after apply)
    + scheme = (known after apply)
    + script = (known after apply)
    ~ select = false -> (known after apply)
    ~ server_ssl = false -> (known after apply)
    ~ set_variable = false -> (known after apply)
    + snat = (known after apply)
    + snatpool = (known after apply)
    ~ source_address = false -> (known after apply)
    ~ ssl_client_hello = false -> (known after apply)
    ~ ssl_server_handshake = false -> (known after apply)
    ~ ssl_server_hello = false -> (known after apply)
    ~ ssl_session_id = false -> (known after apply)
    ~ status = 0 -> (known after apply)
    ~ tcl = false -> (known after apply)
    ~ tcp_nagle = false -> (known after apply)
    + text = (known after apply)
    ~ timeout = 0 -> (known after apply)
    tm_name = "20"
    ~ uie = false -> (known after apply)
    ~ universal = false -> (known after apply)
    + value = (known after apply)
    + virtual = (known after apply)
    + vlan = (known after apply)
    ~ vlan_id = 0 -> (known after apply)
    ~ wam = false -> (known after apply)
    ~ write = false -> (known after apply)
    }
    }
    }

    Plan: 1 to add, 0 to change, 1 to destroy.

    Do you want to perform these actions?
    Terraform will perform the actions described above.
    Only 'yes' will be accepted to approve.

    Enter a value: yes

    bigip_ltm_policy.test-policy: Destroying... [id=new_policy]
    bigip_ltm_policy.test-policy: Destruction complete after 0s
    bigip_ltm_policy.test-policy: Creating...
    bigip_ltm_policy.test-policy: Creation complete after 0s [id=test_policy]

    Apply complete! Resources: 1 added, 0 changed, 1 destroyed.


Once again, Terraform prompts for approval of the execution plan before proceeding. As indicated by the execution plan, Terraform first destroyed the existing instance and then created a new one in its place. You can use terraform show again to see the new values associated with this instance.


Destroying Infrastructure
-------------------------

We've now seen how to build and change infrastructure. Before we move on to creating multiple resources and showing resource dependencies, we're going to go over how to completely destroy the Terraform-managed infrastructure.

Destroying your infrastructure is a rare event in production environments. But if you are using Terraform to spin up multiple environments such as development, test, or QA environments, then destroying is a useful action.

Resources can be destroyed using the ``terraform destroy`` command, which is similar to ``terraform apply`` but it behaves as if all of the resources have been removed from the configuration.

The ``-`` prefix indicates that the instance will be destroyed. As with ``apply``, Terraform shows its execution plan and waits for approval before making any changes. Just like with ``apply``, Terraform determines the order in which things must be destroyed. 


.. code-block:: javascript

   root@terraforn-ubuntu3:~/go/src/github.com/terraform-providers/terraform-provider-bigip# terraform destroy
    bigip_ltm_pool.mypool: Refreshing state... [id=/Common/mypool]
    bigip_ltm_policy.test-policy: Refreshing state... [id=test_policy]

    An execution plan has been generated and is shown below.
    Resource actions are indicated with the following symbols:
    - destroy

    Terraform will perform the following actions:

    # bigip_ltm_policy.test-policy will be destroyed
    - resource "bigip_ltm_policy" "test-policy" {
    - controls = [
    - "forwarding",
    ] -> null
    - id = "test_policy" -> null
    - name = "test_policy" -> null
    - published_copy = "Drafts/test_policy" -> null
    - requires = [
    - "http",
    ] -> null
    - strategy = "/Common/first-match" -> null

    - rule {
    - name = "rule6" -> null

    - action {
    - asm = false -> null
    - avr = false -> null
    - cache = false -> null
    - carp = false -> null
    - classify = false -> null
    - code = 0 -> null
    - compress = false -> null
    - cookie_hash = false -> null
    - cookie_insert = false -> null
    - cookie_passive = false -> null
    - cookie_rewrite = false -> null
    - decompress = false -> null
    - defer = false -> null
    - destination_address = false -> null
    - disable = false -> null
    - enable = false -> null
    - expiry_secs = 0 -> null
    - forward = true -> null
    - hash = false -> null
    - http = false -> null
    - http_basic_auth = false -> null
    - http_cookie = false -> null
    - http_header = false -> null
    - http_host = false -> null
    - http_referer = false -> null
    - http_reply = false -> null
    - http_set_cookie = false -> null
    - http_uri = false -> null
    - insert = false -> null
    - l7dos = false -> null
    - length = 0 -> null
    - log = false -> null
    - ltm_policy = false -> null
    - offset = 0 -> null
    - pem = false -> null
    - persist = false -> null
    - pin = false -> null
    - pool = "/Common/mypool" -> null
    - port = 0 -> null
    - redirect = false -> null
    - remove = false -> null
    - replace = false -> null
    - request = false -> null
    - request_adapt = false -> null
    - reset = false -> null
    - response = false -> null
    - response_adapt = false -> null
    - select = false -> null
    - server_ssl = false -> null
    - set_variable = false -> null
    - source_address = false -> null
    - ssl_client_hello = false -> null
    - ssl_server_handshake = false -> null
    - ssl_server_hello = false -> null
    - ssl_session_id = false -> null
    - status = 0 -> null
    - tcl = false -> null
    - tcp_nagle = false -> null
    - timeout = 0 -> null
    - tm_name = "20" -> null
    - uie = false -> null
    - universal = false -> null
    - vlan_id = 0 -> null
    - wam = false -> null
    - write = false -> null
    }
    }
    }

    # bigip_ltm_pool.mypool will be destroyed
    - resource "bigip_ltm_pool" "mypool" {
    - allow_nat = "yes" -> null
    - allow_snat = "yes" -> null
    - id = "/Common/mypool" -> null
    - load_balancing_mode = "round-robin" -> null
    - monitors = [
    - "/Common/http",
    ] -> null
    - name = "/Common/mypool" -> null
    - reselect_tries = 0 -> null
    - service_down_action = "none" -> null
    - slow_ramp_time = 0 -> null
    }

    Plan: 0 to add, 0 to change, 2 to destroy.

    Do you really want to destroy all resources?
    Terraform will destroy all your managed infrastructure, as shown above.
    There is no undo. Only 'yes' will be accepted to confirm.

    Enter a value: yes

    bigip_ltm_policy.test-policy: Destroying... [id=test_policy]
    bigip_ltm_policy.test-policy: Destruction complete after 0s
    bigip_ltm_pool.mypool: Destroying... [id=/Common/mypool]
    bigip_ltm_pool.mypool: Destruction complete after 0s

    Destroy complete! Resources: 2 destroyed.
    root@terraforn-ubuntu3:~/go/src/github.com/terraform-providers/terraform-provider-bigip#


To read more on BIG-IP Terraform resources and how to use them, see |terraform_doc|.



.. |terraform_doc| raw:: html

   <a href="https://www.terraform.io/docs/providers/bigip/index.html" target="_blank">Terraform documentation</a>