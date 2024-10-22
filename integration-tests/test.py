import json
import base64
from datetime import datetime
from confluent_kafka import Producer, Consumer, KafkaError
from readFile import read_file_to_buffer


KAFKA_BROKER = 'localhost:9092'
MINIO_IFC_URL = 'localhost:9010'
TOPIC = 'ifc-files'
MESSAGES = [
    create_message("project1", "file1.ifc", "Sample IFC file 1 content"),
    create_message("project2", "file2.ifc", "Sample IFC file 2 content"),
    create_message("project1", "file3.ifc", "Sample IFC file 3 content")
]
MINIO_IFC_BUCKET = 'ifc-files'


def create_message(project, filename, content):
    return {
        "project": project,
        "filename": filename,
        "location": create_file_name(project, filename, datetime.now().isoformat()),
        "timestamp": datetime.now().isoformat()
    }

def create_file_name(project, filename, timestamp):
    return f"{project}/{filename}_{timestamp}"


# MinIO Stuff

def add_ifc_file_to_minio(client: Minio, location: str, content: bytes):
    client.put_object(MINIO_IFC_BUCKET, location, content, -1)

def add_ifc_files_to_minio():
    client = Minio(MINIO_IFC_URL)
    content = read_file_to_buffer("assets/test.ifc")
    for msg in MESSAGES:
        add_ifc_file_to_minio(client, msg["location"], content)


# Kafka Stuff

def produce_messages():
    producer = Producer({'bootstrap.servers': KAFKA_BROKER})

    for msg in MESSAGES:
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
    print("Adding IFC files to MinIO...")
    add_ifc_files_to_minio()
    print("Waiting for the files to be processed...")
    time.sleep(10)
    print("Producing messages...")
    produce_messages()
    print("\nConsuming messages...")
    consume_messages()
