# viz-pbi-server

This is a server for providing files and data for a PowerBi dashboard.

## Files

*From storage: MinIO*

route: **/file**

Description:

Returns a fragments file compressed as .gz

Arguments:
  
  - **name**: The name of the file in the following format: `project2/file2_2024-10-25T16:36:04.986173Z.gz`



## Data

*From storage: MongoDB*

route: **/data**

Description:

Returns a dataset for given project and file. Includes all history of that file.

Arguments:

- **db**: the name of the project
- **collection**: the name of the file
