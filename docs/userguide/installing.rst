Installing Terraform
====================

1. To install Terraform, find the |appropriate_package| for your system and download it. Terraform is packaged as a zip archive. Verify the SHA256 checksums before opening the zip file to ensure you are not using a maliciously modified version of Terraform.

2. After downloading Terraform, unzip the package. Terraform runs as a single binary named ``terraform``. Any other files in the package can be safely removed and Terraform will still function.

3. Make sure that the ``terraform`` binary is available on the ``PATH``. For instructions on setting the PATH on Linux and Mac, see |path_linux|. For instructions on setting the PATH on Windows, see |path_windows|.


Verifying the Installation
--------------------------

After installing Terraform, verify the installation worked by opening a new terminal session and run the command ``terraform``. You should see help output similar to this:


.. code-block:: javascript
   :linenos:


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

Now you should be able to build infrastructure on the BIG-IP system (for example: policy configuration) using terraform resources.




.. |appropriate_package| raw:: html

   <a href="https://www.terraform.io/downloads.html" target="_blank">appropriate package</a>


.. |path_linux| raw:: html

   <a href="https://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux-unix" target="_blank">this page</a>


.. |path_windows| raw:: html

   <a href="https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows" target="_blank">this page</a>
