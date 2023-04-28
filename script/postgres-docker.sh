#!/bin/bash

# It is important that snapshotscopy{to|from} path is accessible at the same location in
# the container and outside of it. If you're using a custom zeta home, you must call this script
# with ZETA_HOME set to your custom zeta home when starting the database.

if [ -n "$ZETA_HOME" ]; then
        ZETA_STATE=${ZETA_HOME}/state
else
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
                ZETA_STATE=~/.local/state/zeta
        elif [[ "$OSTYPE" == "darwin"* ]]; then
                ZETA_STATE="${HOME}/Library/Application Support/zeta"
        else
                 echo "$OSTYPE" not supported
        fi
fi

SNAPSHOTS_COPY_TO_PATH=${ZETA_STATE}/data-node/networkhistory/snapshotscopyto
SNAPSHOTS_COPY_FROM_PATH=${ZETA_STATE}/data-node/networkhistory/snapshotscopyfrom

mkdir -p "$SNAPSHOTS_COPY_TO_PATH"
chmod 777 "$SNAPSHOTS_COPY_TO_PATH"

mkdir -p "$SNAPSHOTS_COPY_FROM_PATH"
chmod 777 "$SNAPSHOTS_COPY_FROM_PATH"

docker run --rm \
           -e POSTGRES_USER=zeta \
           -e POSTGRES_PASSWORD=zeta \
           -e POSTGRES_DB=zeta \
           -p 5432:5432 \
           -v "$SNAPSHOTS_COPY_TO_PATH":"$SNAPSHOTS_COPY_TO_PATH":z \
           -v "$SNAPSHOTS_COPY_FROM_PATH":"$SNAPSHOTS_COPY_FROM_PATH":z \
           timescale/timescaledb:2.8.0-pg14
