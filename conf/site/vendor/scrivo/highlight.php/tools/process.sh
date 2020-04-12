#!/bin/bash

HLJS_V="9.18.1"
DOJO_V="1.13.0"

HLJS_DL="https://api.github.com/repos/highlightjs/highlight.js/tarball/$HLJS_V"
DOJO_DL="http://download.dojotoolkit.org/release-$DOJO_V/dojo-release-$DOJO_V-src.tar.gz"

curl -L $HLJS_DL --output "lib_highlight.tar.gz"
curl -L $DOJO_DL --output "lib_dojo.tar.gz"

rm -rf lib_dojo lib_highlight 2> /dev/null

mkdir lib_dojo lib_highlight

tar xzf lib_dojo.tar.gz -C lib_dojo --strip-components 1
tar xzf lib_highlight.tar.gz -C lib_highlight --strip-components 1

cd lib_highlight
npm install
node tools/build.js -t node

# Translate language definitions to something Highlight.php can handle
cd ..
node launcher.js > languages.dat
php get_language_definitions.php

# Copy styles from highlight.js to our own styles directory
rm -r ../styles/ 2> /dev/null
mkdir -p ../styles/
cp -a lib_highlight/src/styles/ ../styles/
php get_styles_colors.php

# Copy unit tests
rm -r ../test/detect/ 2> /dev/null
rm -r ../test/markup/ 2> /dev/null

mkdir -p ../test/{detect,markup}/
cp -a lib_highlight/test/detect/ ../test/detect/
cp -a lib_highlight/test/markup/ ../test/markup/

rm ../test/{detect,markup}/index.js 2> /dev/null
