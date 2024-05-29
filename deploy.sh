#!/bin/sh

podman pull ghcr.io/tcurdt/test-project:live
podman image inspect ghcr.io/tcurdt/test-project:live | jq -r '.[0].Config.Labels.SHA'





IMAGES=
  ghcr.io/tcurdt/test-project:live

# cleanup
git checkout .
git clean -fdx

# get updates
git pull --rebase

# get the current SHA
regctl image inspect ghcr.io/tcurdt/test-project:live
regctl tag ls ghcr.io/tcurdt/test-project
regctl image inspect ghcr.io/tcurdt/test-project:live | jq '.config.Labels.SHA'


IMAGE=ghcr.io/tcurdt/test-project
TAG=live

digest="$(regctl image digest $IMAGE:$TAG)"

for tag in $(regctl tag ls $IMAGE --include 'commit-.*'); do
  echo "checking $tag"
  if [ "$digest" = "$(regctl image digest $IMAGE:${tag})" ]; then
    echo "match $tag"
    break
  fi
done

ociresolve -i in -o out

kubectl apply



# prepare the manifests
ruplacer 'from' 'to' --no-regex .

# apply

kubectl apply -f
podman play kube ./test.yaml

# cleanup
git checkout .
