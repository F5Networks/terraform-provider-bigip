# AS3 Deployment using null resource
- This shows details about how you can deploy AS3 RPM using null resource calling shell scripts to deploy AS3 RPM and AS3 Json payload. 
# How you can use AS3 with null resource ?
- Look at the ``as3.tf`` file it uses null resources twice 
- First null resource ``install_as3" uses script ``install_as3.sh`` script to load the as3 rpm on the BIG-IP
- Second null resource uses shell script ``as3_http.sh`` to deploy example1.json.
