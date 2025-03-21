-- LCA data table
IF NOT EXISTS (SELECT *
FROM sys.tables
WHERE name = 'lca_data')
BEGIN
  CREATE TABLE [dbo].[lca_data]
  (
    [project] VARCHAR(255) NOT NULL,
    [filename] VARCHAR(255) NOT NULL,
    [fileid] VARCHAR(255) NOT NULL,
    [timestamp] DATETIME2 NOT NULL,
    [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
    [mat_kbob] VARCHAR(255),
    [gwp_absolute] DECIMAL(18,2),
    [gwp_relative] DECIMAL(18,2),
    [penr_absolute] DECIMAL(18,2),
    [penr_relative] DECIMAL(18,2),
    [ubp_absolute] DECIMAL(18,2),
    [ubp_relative] DECIMAL(18,2)
  );
END

-- Cost data table
IF NOT EXISTS (SELECT *
FROM sys.tables
WHERE name = 'cost_data')
BEGIN
  CREATE TABLE [dbo].[cost_data]
  (
    [project] VARCHAR(255) NOT NULL,
    [filename] VARCHAR(255) NOT NULL,
    [fileid] VARCHAR(255) NOT NULL,
    [timestamp] DATETIME2 NOT NULL,
    [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
    [category] VARCHAR(255),
    [level] VARCHAR(255),
    [is_structural] BIT,
    [fire_rating] VARCHAR(255),
    [ebkph] VARCHAR(255),
    [ebkph_1] VARCHAR(255),
    [ebkph_2] VARCHAR(255),
    [ebkph_3] VARCHAR(255),
    [cost] DECIMAL(18,2),
    [cost_unit] DECIMAL(18,2),
    CONSTRAINT [UQ_cost_data_id_filename_timestamp_project] UNIQUE ([id], [filename], [timestamp], [project])
  );
END