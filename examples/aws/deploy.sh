#!/bin/sh
#### This is deploy script which helps you to deploy F5 BIG-IP in AWS then extract the Public IP (management IP) of BIG-IP and use it to do other tasks on BIG-IP like
# configuring VIPs, Pool, Monitor. This script uses two TF files the first one which is master.tf  is used to deploy f5 in AWS and the second TF file called as main.t#f is used to deploy VIPs, Pool etc. Both tf files cannot be on the same directory, they need to be in different directory. Also the temp.tt file needs to be created #where the master.tf file is located. The deploy.sh script needs to be in where the master.tf is located. deploy.sh will create variables.tf automatically.

terraform plan  # Terraform plan for BIG-IP in AWS uses master.tf
terraform apply -auto-approve  >> dump    # Dump the creation in a file called dump
cat dump
cd ../main   # change path to use a different main.tf file for deploying VIPS, Pool etc, 
bigip_management_ip=$(awk -F'f5public_ip =' '$2 ~ /[^[:blank:]]/ {print $2}' ../aws/dump) # Extract the Public IP as Management IP for BIG-IP
echo "BIG-IP is being deployed in AWS ...... hang on"
echo "BIG-IP management is $bigip_management_ip will configure Apps once deployed in AWS as a instance"
countdown()
(
  IFS=:
  set -- $*
  secs=$(( ${1#0} * 3600 + ${2#0} * 60 + ${3#0} ))
  while [ $secs -gt 0 ]
  do
    sleep 1 &
    printf "\r%02d:%02d:%02d" $((secs/3600)) $(( (secs/60)%60)) $((secs%60))
    secs=$(( $secs - 1 ))
    wait
  done
  echo
)
countdown "00:16:00" #Wait for F5 bigip to be functional 

sed "s/x.x.x.x/${bigip_management_ip# }/g" temp.tt >> test.txt   # uses dumy called temp.tt which is  file to create variable file for replacement
perl -e 'use Term::ANSIColor; print color "white"; print "ABC\n"; print color "reset";' |  perl -pe 's/\x1b\[[0-9;]*m//g' < test.txt > variables.tf # Removes special characters lik ^[[0m

terraform plan #Uses main.tf file to deploy applications VIP, Pools etc in main directory
terraform apply -auto-approve
rm test.txt  # Clean up 
rm ../aws/dump # Remove the dump file 






