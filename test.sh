rm -Rf /Users/mateusz.gajewski/Desktop/Projects/src/github.com/caddyserver/buildsrv/builds/*
curl -o test/build.zip http://localhost:5050/download/build\?os\=darwin\&arch\=amd64\&features\=expires%2Cminify%2Cinmemory
curl -o test/build.zip http://localhost:5050/download/build\?os\=darwin\&arch\=amd64\&features\=expires%2Cminify%2Cinmemory

unzip -o test/build.zip -d test
./test/caddy -conf ./test/Caddyfile
#rm test/build.zip
