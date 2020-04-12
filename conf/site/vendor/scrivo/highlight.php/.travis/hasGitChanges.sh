#!/usr/bin/env bash

curl -L https://cs.sensiolabs.org/download/php-cs-fixer-v2.phar -o php-cs-fixer

chmod +x php-cs-fixer

./php-cs-fixer fix

if ! git diff-index --quiet HEAD --; then
    git --no-pager diff
    exit 1
fi
