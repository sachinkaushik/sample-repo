<!--
SPDX-FileCopyrightText: 2021 Open Networking Foundation

SPDX-License-Identifier: LicenseRef-Intel
-->

# SRA Model Mega Patch

To apply this to rocaas-umbrella use 
```shell
curl --location --request PATCH 'http://localhost:8181/roc-api' -d @mega-patch-portland-north.json
curl --location --request PATCH 'http://localhost:8181/roc-api' -d @mega-patch-portland-southeast.json
```

> It should be possible to load config from many different stores in to one Mega Patch file, using
> `additionalProperties` -> `store-id`
> This is currently not working in roc-api - an example will be provided here once it is