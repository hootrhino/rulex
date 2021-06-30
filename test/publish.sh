#!/bin/bash

for i in {1..10000}; do
    mosquitto_pub -h  192.168.0.103 -p 1883 -t '$X_IN_END' -m "{\"temp\": $RANDOM,\"hum\":$RANDOM}"
    echo "Publish ", $i, "{\"temp\": $RANDOM,\"hum\":$RANDOM} Ok."
done
