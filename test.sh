curl -o build.zip http://localhost:5050/download/build\?os\=darwin\&arch\=amd64\&features\=expires%2Cminify%2Cinmemory
unzip -o build.zip
./caddy -conf Caddyfile
rm build.zip
