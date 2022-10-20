Installing Terraform
====================

#. To install Terraform, find and download the `appropriate package <https://www.terraform.io/downloads.html>`_ for your system. Terraform is packaged as a zip archive. Use the `SHA256 checksums for Terraform 0.12.9 <https://releases.hashicorp.com/terraform/0.12.9/terraform_0.12.9_SHA256SUMS>`_ and verify the `checksums signature file <https://releases.hashicorp.com/terraform/0.12.9/terraform_0.12.9_SHA256SUMS.sig>`_ which has been signed using HashiCorp's GPG key before opening the zip file to ensure you are not using a maliciously modified version of terraform.

#. After downloading Terraform, unzip the package. Terraform runs as a single binary named ``terraform``. Any other files in the package can be safely removed and Terraform will still function.

#. Make sure that the ``terraform`` binary is available on the ``PATH``.

   - `Set the PATH on Linux and Mac <https://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux-unix>`_ 
   - `Set the PATH on Windows <https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows>`_


Verifying the Installation
--------------------------

After installing Terraform, verify the installation worked by opening a new terminal session and run the command ``terraform``. You should see help output similar to this:


.. code-block:: javascript

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


If you get an error that ``terraform`` could not be found, your ``PATH`` environment variable was not set up properly. Please go back and ensure that your ``PATH`` variable contains the directory where Terraform was installed.

Now you should be able to build infrastructure on the F5 BIG-IP system (for example: policy configuration) using terraform resources.