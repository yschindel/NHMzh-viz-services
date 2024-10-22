# Integration Tests

## Prerequisites

- Docker
- Docker Compose
- Python 3.11

### Setting up the Python environment

Open a terminal in the `integration-tests` directory.

1. Create a virtual environment:

```bash
python -m venv nhmzh-viz-services
```

2. Activate the virtual environment:

On Unix or MacOS:

```bash
source nhmzh-viz-services/bin/activate
```

On Windows:

```bash
nhmzh-viz-services\Scripts\activate
```

3. Install the dependencies:

```bash
pip install -r requirements.txt
```

## Running the tests

Open a new terminal window in the `integration-tests` directory:

```bash
cd integration-tests
```

### Start the services

Run the following command to start the services:

```bash
docker compose up --build
```

### Open the MinIO console

In your browser, go to `localhost:9001` to see the MinIO console.

- Delete all items in the `ifc-files` bucket.
- Delete the `ifc-files` bucket.

### Add content to the Kafka topic

Run the following command to add content to the Kafka topic, use the terminal window with the virtual python environment activated:

```bash
python test.py
```

### Verify the content was added correctly

Verify that the bucket was created and that the content was added correctly:

1. Refresh the Object Browser in the MinIO console.
2. You should see a new bucket called `ifc-fragment-files`.
3. Open the `ifc-fragment-files` bucket. If no folders show up, try refreshing the bucket.
4. Open the folders and check that the files are present.

If you see the files, the integration tests passed.
