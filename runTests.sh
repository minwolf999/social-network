#!/bin/bash

find "." -type f -name '*test.go' | while read -r file; do
    dir=$(dirname "$file")
    echo "Exécution des tests dans le répertoire : $dir"
    (cd "$dir" && go test)
done
