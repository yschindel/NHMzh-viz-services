# Azure Cloud Tests

This docker compose is used to setup a container network that interacts with the Azure Cloud Services so that PowerBI can fetch data via direct query.
It is basically the 'backend' for development in the PowerBI 'frontend', since PowerBI direct query is most comfortable to run on Azure Cloud Services. What a surprise...

## PowerBI Custom Visual 3D Viewer

The Custom PBI ifc-viewer that is designed to work with these services is fetching an object/blob from the pbi-server which is acting as a proxy between the custom visual and Azure Blob Storage. Depending on your testing scenario, you can set the `AZURE_STORAGE...` variables to a real cloud instance or the local emulators included in this docker-compose.

Make sure to pull the images before running them:

```bash
docker compose pull
```

Create a .env file in the `cloud-tests` directory and add the following variables if you want to connect to Azure Cloud Services:

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

Spin up the container network with the .env file:

```bash
docker compose --env-file .env up -d
```

Now you can run the test script `integration-tests/test/main.go` using:

```bash
go run main.go
```

This will simulate a full data run starting with an IFC upload. The key thing here is that because it's starting from a single IFC file and then mocking the lca and cost calculations, all data including the 3D model will match up by globalId.
