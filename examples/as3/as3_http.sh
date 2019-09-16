: '
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
'
#!/bin/bash

echo "\033[1m...Deploying HTTP Virtual Server and Pools ....... \033[0m "
curl -k --user admin:pass -H "Accept: application/json" -H "Content-Type:application/json" -X POST -d@example1.json https://X.X.X.X/mgmt/shared/appsvcs/declare | python -m json.tool
