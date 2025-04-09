# Docker Hub Tests

This docker compose is used to test and show how the published docker images work together.

Make sure to pull the images before running them:

```bash
docker compose pull
```

Create a .env file in the `docker-tests` directory and add the following variables if you want to connect to Azure Cloud Services:

```
AZURE_DB_SERVER=abc.database.windows.net
AZURE_DB_PORT=1433
AZURE_DB_USER=<your-user>
AZURE_DB_PASSWORD=<your-password>
AZURE_DB_DATABASE=<your-database>

AZURE_STORAGE_ACCOUNT=<your-storage-account>
AZURE_STORAGE_KEY=<your-storage-account-key>
AZURE_STORAGE_URL=<your-storagae-url>
```
