#!/bin/bash

# Set the container name as a variable
CONTAINER_NAME="viz-test-kafka-container"
TOPIC="ifc-files"

# Function to create a JSON message
create_message() {
    local project=$1
    local filename=$2
    local content=$3
    local timestamp=$(date +%s%N | cut -b1-13)  # Current timestamp in milliseconds
    echo "{\"project\":\"$project\",\"filename\":\"$filename\",\"timestamp\":$timestamp,\"content\":\"$content\"}"
}

# Sample messages to send
MESSAGES=(
    "$(create_message "project1" "file1.ifc" "$(echo -n "Sample IFC file 1 content" | base64 -w 0)")"
    "$(create_message "project2" "file2.ifc" "$(echo -n "Sample IFC file 2 content" | base64 -w 0)")"
    "$(create_message "project1" "file3.ifc" "$(echo -n "Sample IFC file 3 content" | base64 -w 0)")"
)

# Send messages to the Kafka topic
for msg in "${MESSAGES[@]}"; do
    echo "$msg" | docker exec -i $CONTAINER_NAME kafka-console-producer.sh \
        --bootstrap-server localhost:9092 \
        --topic $TOPIC
done

echo "Messages sent to topic $TOPIC"

# Read messages from the beginning of the topic and log to console
echo "Reading messages from topic $TOPIC:"
docker exec -i $CONTAINER_NAME kafka-console-consumer.sh \
    --bootstrap-server localhost:9092 \
    --topic $TOPIC \
    --from-beginning \
    --max-messages ${#MESSAGES[@]}

echo "Finished reading messages from topic $TOPIC"
