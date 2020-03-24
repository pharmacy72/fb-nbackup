#!/bin/sh
FIREBIRD_PASSWORD="023RsdTf4UI123"
FIREBIRD_USER="fbuser"
FIREBIRD_DATABASE="NBEXAMPLE"

docker rm -f firebird
docker run \
	-v backup:/backup \
	-e FIREBIRD_PASSWORD=$FIREBIRD_PASSWORD \
       	-e FIREBIRD_USER=$FIREBIRD_USER\
	-e FIREBIRD_DATABASE=$FIREBIRD_DATABASE\
	--name firebird -d  jacobalberty/firebird:3.0

