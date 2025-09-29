# PowerBI Server

Code available at: [https://github.com/yschindel/NHMzh-viz-services/tree/main/viz_pbi-server](https://github.com/yschindel/NHMzh-viz-services/tree/main/viz_pbi-server)

Scripts for loading data in PowerBI: [https://github.com/yschindel/NHMzh-pbi-queries](https://github.com/yschindel/NHMzh-pbi-queries)

This is a server for providing data for a PowerBi dashboard. All data is stored in MinIO as .parquet files.
The file paths are project/file. The directory structure is an important part of the information in the system as it is used to inform the PowerBi dashboard of the data that is available.

## Security Features

### API Key Authentication

All endpoints are protected with API key authentication. Clients must include a valid API key in the `X-API-Key` header with each request.

To configure the API key:

- Set the `PBI_SRV_API_KEY` environment variable to your desired API key value

If no API key is set, the server will run without authentication but will display a warning at startup.

### Deployment Behind a Reverse Proxy

This service is designed to run behind a reverse proxy (like Nginx, Traefik, or HAProxy) which handles:

- TLS termination (HTTPS)
- Certificate management
- HTTP to HTTPS redirection

The service supports the standard proxy headers:

- `X-Forwarded-Proto` - For detecting the original protocol (http/https)
- `X-Forwarded-Host` - For detecting the original host

### Configuration

The service is configured using the following environment variables:

| Environment Variable       | Description                          | Default              |
| -------------------------- | ------------------------------------ | -------------------- |
| `PBI_SERVER_PORT`          | The port on which the server listens | `3000`               |
| `PBI_SRV_API_KEY`          | API key for authentication           | (none)               |
| `MINIO_ENDPOINT`           | MinIO server endpoint                | (required)           |
| `MINIO_ACCESS_KEY`         | MinIO access key                     | (required)           |
| `MINIO_SECRET_KEY`         | MinIO secret key                     | (required)           |
| `VIZ_IFC_FRAGMENTS_BUCKET` | Bucket/Container for fragment files  | `ifc-fragment-files` |

## API Endpoints

### Get Fragment File

**Route:** `GET /blob`

**Description:**
Returns a compressed fragment file (.gz) for 3D visualization in PowerBI.

**Parameters:**

- `id` (query): File ID in format `project/filename_timestamp.gz`
- `container` (query, optional): Storage container name (defaults to configured bucket)

**Authentication:**

- Required `X-API-Key` header with valid API key

**Usage:**

```bash
curl -H "X-API-Key: your-api-key" \
     "http://localhost:3000/blob?id=project1/building1_20240115.gz"
```

### Upload Fragment File

**Route:** `POST /blob`

**Description:**
Uploads a compressed fragment file to storage.

**Parameters:**

- `container` (query, optional): Storage container name
- File data in multipart form

**Authentication:**

- Required `X-API-Key` header with valid API key

### Store Element Data

**Route:** `POST /data/elements`

**Description:**
Receives and stores processed element data from viz_lca-cost service.

**Content-Type:** `application/json`

**Authentication:**

- Required `X-API-Key` header with valid API key

### Store Material Data

**Route:** `POST /data/materials`

**Description:**
Receives and stores processed material data from viz_lca-cost service.

**Content-Type:** `application/json`

**Authentication:**

- Required `X-API-Key` header with valid API key

### CORS Preflight

**Route:** `OPTIONS /{path:.*}`

**Description:**
Handles CORS preflight requests for PowerBI integration.

**Authentication:**

- None required for OPTIONS requests
