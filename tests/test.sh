#!/bin/bash

set -x
set -e

DIR=$(dirname "$0")
curl -X POST 'http://localhost:9009/solve' -H 'Content-Type: application/json' -d@$DIR/1.json
