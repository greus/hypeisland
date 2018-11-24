#!/usr/bin/env bash

minify --type js index.js > index.min.js

sed -e "s^//TODO^$(<index.min.js)^g" index.html > index.bundled.html

sed -i -e 's,<script src=index.js></script>,,' index.bundled.html
