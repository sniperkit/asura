#! /bin/bash
set -e

# test the counter using a go test script
bash tests/test_app/test.sh

# test the cli against the examples in the tutorial at teragrid.com
# TODO: make these less fragile
# bash tests/test_cli/test.sh


