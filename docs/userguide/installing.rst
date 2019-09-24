Installing Terraform
====================

To install Terraform, find the appropriate package for your system and download it. Terraform is packaged as a zip archive.

After downloading Terraform, unzip the package. Terraform runs as a single binary named terraform. Any other files in the package can be safely removed and Terraform will still function.

The final step is to make sure that the terraform binary is available on the PATH. See this page for instructions on setting the PATH on Linux and Mac. This page contains instructions for setting the PATH on Windows.


Verifying the Installation
--------------------------

After installing Terraform, verify the installation worked by opening a new terminal session and checking that `terraform` is available. By executing `terraform` you should see help output similar to this:

$ terraform
Usage: terraform [--version] [--help] <command> [args]

The available commands for execution are listed below.
The most common, useful commands are shown first, followed by
less common or more advanced commands. If you're just getting
started with Terraform, stick with the common commands. For the
other commands, please read the help and docs before usage.

Common commands:
    apply              Builds or changes infrastructure
    console            Interactive console for Terraform interpolations
# ...


If you get an error that `terraform` could not be found, your PATH environment variable was not set up properly. Please go back and ensure that your PATH variable contains the directory where Terraform was installed.

With Terraform installed, let's dive right into it . We'll build infrastructure on BIG-IP (example: policy configuration) using terraform resources, and for this we need to have bigip ip and its credentials.