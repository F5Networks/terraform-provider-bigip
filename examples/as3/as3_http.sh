#!/bin/bash

echo "\033[1m...Deploying HTTP Virtual Server and Pools ....... \033[0m "
curl -k --user admin:pass -H "Accept: application/json" -H "Content-Type:application/json" -X POST -d@example1.json https://X.X.X.X/mgmt/shared/appsvcs/declare | python -m json.tool
