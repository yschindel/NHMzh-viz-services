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

## Routes

### Get A Fragment File

route: **/fragments**

Description:

Returns the latest version of a fragments file compressed as .gz

Arguments:

- **id**: The id of the file in the following format: `project2/file2`

Authentication:

- Required `X-API-Key` header with valid API key

### List All Fragment Files

route: **/fragments/list**

Description:

Returns a list of all fragment files in the `ifc-fragment-files` bucket.

Arguments:

- None

Authentication:

- Required `X-API-Key` header with valid API key

### Get A Data File

route: **/data**

Description:

Returns a dataset (as a .parquet file) for given project and file. Includes all history of that file.

Arguments:

- **project**: The project name
- **file**: The file name

Authentication:

- Required `X-API-Key` header with valid API key

### List All Data Files

route: **/data/list**

Description:

Returns a list of all data files in the `lca-cost-data` bucket.

Arguments:

- None

Authentication:

- Required `X-API-Key` header with valid API key
