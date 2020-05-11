#!/bin/sh
FIREBIRD_PASSWORD="023RsdTf4UI123"
FIREBIRD_USER="fbuser"
FIREBIRD_DATABASE="NBEXAMPLE"
FB_CONTAINER=firebird

docker rm -f ${FB_CONTAINER}
docker run \
	-v "${PWD}"/backup:/backup \
	-e FIREBIRD_PASSWORD=$FIREBIRD_PASSWORD \
  -e FIREBIRD_USER=$FIREBIRD_USER\
	-e FIREBIRD_DATABASE=$FIREBIRD_DATABASE\
	--name ${FB_CONTAINER} -d  jacobalberty/firebird:3.0

ISQL="docker exec -i $FB_CONTAINER /usr/local/firebird/bin/isql -e -user ${FIREBIRD_USER} -password ${FIREBIRD_PASSWORD} -database ${FIREBIRD_DATABASE}"
$ISQL < $(dirname "$0")/create.db.sql
