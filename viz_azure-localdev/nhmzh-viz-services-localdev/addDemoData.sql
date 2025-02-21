BULK INSERT [dbo].[project_data]
FROM 'export_data.csv'
WITH
(
    FORMAT = 'CSV',
    FIRSTROW = 2,  -- Skip header row
    FIELDTERMINATOR = ',',
    ROWTERMINATOR = '\n',
    CODEPAGE = '65001'  -- UTF-8 encoding
);