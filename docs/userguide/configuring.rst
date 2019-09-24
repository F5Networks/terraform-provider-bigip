Configuring Terraform
=====================

The set of files used to describe infrastructure in Terraform is simply known as a Terraform configuration.



Below is the sample terraform resource to create a policy on the BIG-IP system.

.. code-block:: javascript
   :linenos:

    provider "bigip" {
        address = "x.x.x.x"
        username = "xxxx"
        password = "xxxx"
    }
 
    resource "bigip_ltm_policy" "test-policy" {
        name = "my_policy"
        strategy = "first-match"
        requires = ["http"]
        published_copy = "Drafts/my_policy"
        controls = ["forwarding"]
        rule {
            name = "rule6"
            action = {
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



The provider block is used to configure the named provider, in our case "bigip". A provider is responsible for creating and managing resources. Multiple provider blocks can exist if a Terraform configuration is composed of multiple providers, which is a common situation.

The resource block defines a resource that exists within the infrastructure. The resource block has two strings before opening the block: the resource type and the resource name. In our example, the resource type is "bigip_ltm_policy" and the name is "test_policy." The prefix of the type maps to the provider. In our case "bigip_ltm_policy" automatically tells Terraform that it is managed by the "bigip" provider.


+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| Parameter          | Options              | Description/Notes                                                                                                         |
+====================+======================+===========================================================================================================================+
| name               | Required             | Describes the name of the policy.                                                                                         |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| strategy           | Optional             | This value specifies the match strategy.                                                                                  |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| requires           | Optional             | This value specifies the protocol.                                                                                        |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| published_copy     | Optional             | This value determines if you want to publish the policy else it will be deployed in Drafts mode.                          |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| controls           | Optional             | This value specifies the controls.                                                                                        |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| rule               | Optional             | Use this policy to apply rules.                                                                                           |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| tm_name            | Required             | If Rule is used, then you need to provide the tm_name. It can be any value.                                               |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| forward            | Optional             | This value sets forwarding.                                                                                               |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+
| pool               | Optional             | This value will direct the stream to this pool.                                                                           |
+--------------------+----------------------+---------------------------------------------------------------------------------------------------------------------------+

