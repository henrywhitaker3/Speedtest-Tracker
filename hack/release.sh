#!/bin/bash

output=$(./hack/version.sh)

if [[ $(echo "$output" | jq -r '.skip') != "true" ]]; then
  version=$(echo "$output" | jq -r '.version')
  tag=$(echo "$output" | jq -r '.tag')
  changelog=$(echo "$output" | jq -r '.changelog')

  echo "Releasing $version"

  echo "$changelog" >CHANGELOG.md
  git add CHANGELOG.md
  git commit -m "chore(docs): update changelog"
  git push
  git tag "$tag"
  git push origin "$tag"
else
  echo "No new version to release"
fi
