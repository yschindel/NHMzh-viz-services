#!/bin/bash

# Set the container name as a variable
CONTAINER_NAME="viz-test-kafka-container"
TOPIC="ifc-files"

# Sample messages to send
MESSAGES=(
    '{"id": 1, "content": "Sample IFC file 1"}'
    '{"id": 2, "content": "Sample IFC file 2"}'
    '{"id": 3, "content": "Sample IFC file 3"}'
)

# Send messages to the Kafka topic
for msg in "${MESSAGES[@]}"; do
    echo "$msg" | docker exec -i $CONTAINER_NAME /opt/kafka/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic $TOPIC
done

echo "Messages sent to topic $TOPIC"

# Read messages from the beginning of the topic and log to console
echo "Reading messages from topic $TOPIC:"
docker exec -it $CONTAINER_NAME /opt/kafka/bin/kafka-console-consumer.sh \
    --bootstrap-server localhost:9092 \
    --topic $TOPIC \
    --from-beginning \
    --max-messages ${#MESSAGES[@]}

echo "Finished reading messages from topic $TOPIC"
