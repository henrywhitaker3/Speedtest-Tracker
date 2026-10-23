#!/bin/bash

output=$(git cliff --tag-pattern "v2.*" --bump --unreleased --context)
skip="false"
changelog=$(git cliff --tag-pattern "v2.*" --bump)

if [[ $(echo "$output" | jq -r '.[0].bump_type') == "null" ]] && [[ $(echo "$output" | jq -r '.[0].previous') != "null" ]]; then
  skip="true"
fi

tag=$(echo "$output" | jq -r '.[0].version')
version=${tag#"$sanitised-"}

jq -n \
  --arg tag "$tag" \
  --arg version "$version" \
  --arg changelog "$changelog" \
  --arg skip "$skip" \
  '{"tag": $tag, "version": $version, "changelog": $changelog, "skip": $skip | test("true")}'
