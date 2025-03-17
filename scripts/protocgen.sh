#!/usr/bin/env bash

set -e

echo "--> Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # this regex checks if a proto file has its go_package set to github.com/troykessler/mailbox/api/...
    # gogo proto files SHOULD ONLY be generated if this is false
    # we don't want gogo proto to run for proto files which are natively built for google.golang.org/protobuf
    if grep -q "option go_package" "$file" && grep -H -o -c 'option go_package.*github.com/troykessler/hyperlane-cosmos/api' "$file" | grep -q ':0$'; then
      buf generate --template buf.gen.gogo.yaml "$file"
    fi
  done
done

echo "--> Generating pulsar proto code"
module_list=$(find . -name "*module.proto" -not -path "./hyperlane/core/_*" | tr '\n' ','  | sed 's/,$//')
buf generate --template buf.gen.pulsar.yaml --path "$module_list"

cd ..

cp -r github.com/troykessler/hyperlane-cosmos/* ./
rm -rf api && mkdir api
mv hyperlane/* ./api
rm -rf github.com troykessler hyperlane