import json
import base64
from datetime import datetime
from confluent_kafka import Producer, Consumer, KafkaError

KAFKA_BROKER = 'localhost:9092'
TOPIC = 'ifc-files'

def create_message(project, filename, content):
    return {
        "project": project,
        "filename": filename,
        "content": base64.b64encode(content.encode()).decode(),
        "timestamp": datetime.now().isoformat()
    }

def produce_messages():
    producer = Producer({'bootstrap.servers': KAFKA_BROKER})

    messages = [
        create_message("project1", "file1.ifc", "Sample IFC file 1 content"),
        create_message("project2", "file2.ifc", "Sample IFC file 2 content"),
        create_message("project1", "file3.ifc", "Sample IFC file 3 content")
    ]

    for msg in messages:
        producer.produce(TOPIC, json.dumps(msg).encode('utf-8'))
        print(f"Sent message: {msg['project']}/{msg['filename']}")

    producer.flush()

def consume_messages():
    consumer = Consumer({
        'bootstrap.servers': KAFKA_BROKER,
        'group.id': 'test-group',
        'auto.offset.reset': 'earliest'
    })
    consumer.subscribe([TOPIC])

    try:
        while True:
            msg = consumer.poll(1.0)
            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    continue
                else:
                    print(f"Consumer error: {msg.error()}")
                    break

            value = json.loads(msg.value().decode('utf-8'))
            print(f"Received message: {value['project']}/{value['filename']}")
            print(f"Content: {base64.b64decode(value['content']).decode()}")
    finally:
        consumer.close()

if __name__ == "__main__":
    print("Producing messages...")
    produce_messages()
    print("\nConsuming messages...")
    consume_messages()
