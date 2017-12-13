echo " This script you need to run once you download the repo from the github, this will fix the dependencies issues also this will build the binary for you which will be used for BIG-IP configuration using terraform" 


#!/bin/bash
go get github.com/aws/aws-sdk-go/aws
go get github.com/aws/aws-sdk-go/aws/ec2metadata
go get github.com/aws/aws-sdk-go/aws/session
go get github.com/aws/aws-sdk-go/service/s3
go get github.com/bgentry/speakeasy
go get github.com/mattn/go-isatty
go get github.com/posener/complete
go get github.com/posener/complete/cmd/install
go get github.com/armon/go-radix
go get golang.org/x/crypto/openpgp/armor
go get golang.org/x/crypto/openpgp/errors
go get golang.org/x/crypto/openpgp/packet
go get golang.org/x/crypto/openpgp/s2k


ls
cd terraform-provider-bigip/
 go build -o terraform-provider-bigip
#
