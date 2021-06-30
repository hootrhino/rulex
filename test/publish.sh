#!/bin/bash

for i in {1..1000}; do
    mosquitto_pub -h  127.0.0.1 -p 1883 -t '$X_IN_END' -q 2 -m "{\"temp\": $RANDOM,\"hum\":$RANDOM}"
    echo "Publish ", $i, "{\"temp\": $RANDOM,\"hum\":$RANDOM} Ok."
done
