# PowerBI Server

This is a server for providing data for a PowerBi dashboard. All data is stored in MinIO as .parquet files.
The file paths are project/file. The directory structure is an important part of the information in the system as it is used to inform the PowerBi dashboard of the data that is available.

## Routes

### Get A Fragment File

route: **/fragments**

Description:

Returns a fragments file compressed as .gz

Arguments:

- **name**: The name of the file in the following format: `project2/file2_2024-10-25T16:36:04.986173Z.gz`

### List All Fragment Files

route: **/fragments/list**

Description:

Returns a list of all fragment files in the `ifc-fragment-files` bucket.

Arguments:

- None

### Get A Data File

route: **/data**

Description:

Returns a dataset (as a .parquet file) for given project and file. Includes all history of that file.

Arguments:

- **project**: The project name
- **file**: The file name

### List All Data Files

route: **/data/list**

Description:

Returns a list of all data files in the `lca-cost-data` bucket.

Arguments:

- None
