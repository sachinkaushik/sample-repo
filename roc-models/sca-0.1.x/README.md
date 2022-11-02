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
--git-repo-path tmp/roc-api-sca-0.1.x \
--spec ../../frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sca-0.1.x/openapi.yaml \
--package sca_0_1_x \
--output /tmp/roc-api-sca-0.1.x \
--import-module github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sca-0.1.x:v0.1.0
```
* contd...
    * `cd  /tmp/roc-api-sca-0.1.x`
    * Run `make docker-build`
    * Creates docker image `sca_0_1_x/roc-rest-api:latest`
    * tag this as `docker tag sca_0_1_x/roc-rest-api:latest amr-registry.caas.intel.com/one-intel-edge/roc/roc-api-sca-0.1.x:latest`

> Note this code does not need to be versioned in Github and can be discarded after the docker image is built

* Create the sca-0.1.x adapter stub
    * Run
 ```shell
./build/_output/roc create stub-generator \
--git-repo-path tmp/roc-adapter-sca-0.1.x \
--package sca_0_1_x \
--output /tmp/roc-adapter-sca-0.1.x \
--import-model github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sca-0.1.x:v0.1.0
```
* contd...
    * `cd  /tmp/roc-adapter-sca-0.1.x`
    * Run `make docker-build`
    * Creates docker image `library/sca_0_1_x:latest`
    * tag this as `docker tag library/sca_0_1_x:latest amr-registry.caas.intel.com/one-intel-edge/roc/roc-adapter-sca_0_1_x:latest`

> Note since this code is meant to be extended by a developer it makes sense to version it in Github
> For the SCA use case this has already been versioned and extended in
> [https://github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-adapters.sca](https://github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-adapters.sca)

* Clone https://github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models
    * stay in base directory
    * Run `make docker-build`
    * Creates docker image `amr-registry.caas.intel.com/one-intel-edge/roc/sca-0.1.x:0.2.0-dev-sca-0.1.x`
