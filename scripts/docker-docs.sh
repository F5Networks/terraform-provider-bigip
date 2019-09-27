#!/usr/bin/env bash

set -x

: ${DOC_IMG:=quay.pdsea.f5net.com/doc-ops/containthedocs:latest}

exec docker run --rm -it \
  -v $PWD:$PWD --workdir $PWD \
  ${DOCKER_RUN_ARGS} \
  ${DOC_IMG} "$@"
