#!/bin/bash

# builds mac+linux, on mac.

for db in mssql oracle postgresql; do
    ./run-ci.sh $db
    ./wrapper/scripts/wrap-$db.sh
done