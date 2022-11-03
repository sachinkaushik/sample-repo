# ROC files for existing RIs

For `sra-0.2x` model example to generate file is (run from `$HOME/git/intel-iotg-decloud/sample-repo`):

```shell
docker run -v $HOME/git/intel-iotg-decloud/sample-repo/roc-models:/models -v $HOME/git/intel-iotg-decloud/sample-repo/roc-models/sca-0.1.x/out:/tmp/out -w /models --env DOCKER_REGISTRY=default-route-openshift-image-registry.apps.nex.one-edge.intel.com --env DOCKER_REPOSITORY_BASE=springboard-dev-common/ --entrypoint="/usr/local/bin/generate-code.sh" amr-registry.caas.intel.com/one-intel-edge/roc/roc-cli:0.0.4 -m sca -r 0.1.x -a -d
```

> Adjust accordingly for your checked out folder and the RI model you are generating.