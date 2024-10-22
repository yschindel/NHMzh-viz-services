import json
import base64
from datetime import datetime
from kafka import KafkaProducer
from kafka import KafkaConsumer

KAFKA_BROKER = 'localhost:9092'
TOPIC = 'ifc-files'

def create_message(project, filename, content):
    return {
        "project": project,
        "filename": filename,
        "timestamp": int(datetime.now().timestamp() * 1000),
        "content": base64.b64encode(content.encode()).decode()
    }

def produce_messages():
    producer = KafkaProducer(
        bootstrap_servers=[KAFKA_BROKER],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )

    messages = [
        create_message("project1", "file1.ifc", "Sample IFC file 1 content"),
        create_message("project2", "file2.ifc", "Sample IFC file 2 content"),
        create_message("project1", "file3.ifc", "Sample IFC file 3 content")
    ]

    for msg in messages:
        producer.send(TOPIC, value=msg)
        print(f"Sent message: {msg['project']}/{msg['filename']}")

    producer.flush()
    producer.close()

def consume_messages():
    consumer = KafkaConsumer(
        TOPIC,
        bootstrap_servers=[KAFKA_BROKER],
        auto_offset_reset='earliest',
        enable_auto_commit=True,
        group_id='test-group',
        value_deserializer=lambda x: json.loads(x.decode('utf-8'))
    )

    for message in consumer:
        print(f"Received message: {message.value['project']}/{message.value['filename']}")
        print(f"Content: {base64.b64decode(message.value['content']).decode()}")

    consumer.close()

if __name__ == "__main__":
    print("Producing messages...")
    produce_messages()
    print("\nConsuming messages...")
    consume_messages()
