### Manually building the ROC toolchain for this model from scratch

> This is an alternate to using the registry - all of this will be automated in future as part of the ROC toolchain.
> This requires you to have lots of dev tools on your system like `node`, `go`, `docker` etc

Clone the following repos:

* Clone https://github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.rocaas-gui
    * Run `TARGET_TYPE=store make docker-build`
    * Creates docker image `roc-ui:0.1.0`
    * tag it with `docker tag roc-ui:0.1.0 amr-registry.caas.intel.com/one-intel-edge/roc/roc-ui:0.1.0`

* Clone https://github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.rocaas-tools
    * Run `make build`
    * Run
 ```shell
./build/_output/roc create api-server \
--git-repo-path tmp/roc-api-sra-0.2.x \
--spec ../../frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x/openapi.yaml \
--package sra_0_2_x \
--output /tmp/roc-api-sra-0.2.x \
--import-module github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x:v0.2.2
```
* contd...
    * `cd  /tmp/roc-api-sra-0.2.x`
    * Run `make docker-build`
    * Creates docker image `sra_0_2_x/roc-rest-api:latest`
    * tag this as `docker tag sra_0_2_x/roc-rest-api:latest amr-registry.caas.intel.com/one-intel-edge/roc/roc-api-sra-0.2.x:latest`

> Note this code does not need to be versioned in Github and can be discarded after the docker image is built

* Create the sra-0.2.x adapter stub
    * Run
 ```shell
./build/_output/roc create stub-generator \
--git-repo-path tmp/roc-adapter-sra-0.2.x \
--package sra_0_2_x \
--output /tmp/roc-adapter-sra-0.2.x \
--import-model github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x:v0.2.2
```
* contd...
    * `cd  /tmp/roc-adapter-sra-0.2.x`
    * Run `make docker-build`
    * Creates docker image `library/sra_0_2_x:latest`
    * tag this as `docker tag library/sra_0_2_x:latest amr-registry.caas.intel.com/one-intel-edge/roc/roc-adapter-sra_0_2_x:latest`

> Note since this code is meant to be extended by a developer it makes sense to version it in Github
> For the SRA use case this has already been versioned and extended in
> [https://github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-adapters.sra](https://github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-adapters.sra)

* Clone https://github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models
    * stay in base directory
    * Run `make docker-build`
    * Creates docker image `amr-registry.caas.intel.com/one-intel-edge/roc/sra-0.2.x:0.2.0-dev-sra-0.2.x`
