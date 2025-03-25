-- EAV data table for elements (LCA and Cost)
IF NOT EXISTS (SELECT *
FROM sys.tables
WHERE name = 'data_eav_elements')
BEGIN
  CREATE TABLE [dbo].[data_eav_elements]
  (
    [project] VARCHAR(255) NOT NULL,
    [filename] VARCHAR(255) NOT NULL,
    [fileid] VARCHAR(255) NOT NULL,
    [timestamp] DATETIME2 NOT NULL,
    [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
    [param_name] VARCHAR(255) NOT NULL,
    [param_value_string] NVARCHAR(MAX),
    [param_value_number] DECIMAL(18,2),
    [param_value_boolean] BIT,
    [param_value_date] DATETIME2,
    [param_type] VARCHAR(50) NOT NULL,
    CONSTRAINT [PK_data_eav_elements] PRIMARY KEY 
    ([project], [filename], [timestamp], [id], [param_name])
  );

  -- Create indexes for common queries
  CREATE INDEX [IX_data_eav_elements_project_fileid] ON [dbo].[data_eav_elements] ([project], [fileid]);
  CREATE INDEX [IX_data_eav_elements_id] ON [dbo].[data_eav_elements] ([id]);
  CREATE INDEX [IX_data_eav_elements_timestamp] ON [dbo].[data_eav_elements] ([timestamp]);
END

-- EAV data table for materials (LCA and Cost)
IF NOT EXISTS (SELECT *
FROM sys.tables
WHERE name = 'data_eav_materials')
BEGIN
  CREATE TABLE [dbo].[data_eav_materials]
  (
    [project] VARCHAR(255) NOT NULL,
    [filename] VARCHAR(255) NOT NULL,
    [fileid] VARCHAR(255) NOT NULL,
    [timestamp] DATETIME2 NOT NULL,
    [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
    [sequence] INT NOT NULL,
    [param_name] VARCHAR(255) NOT NULL,
    [param_value_string] NVARCHAR(MAX),
    [param_value_number] DECIMAL(18,2),
    [param_value_boolean] BIT,
    [param_value_date] DATETIME2,
    [param_type] VARCHAR(50) NOT NULL,
    CONSTRAINT [PK_data_eav_materials] PRIMARY KEY 
    ([project], [filename], [timestamp], [id], [sequence], [param_name])
  );

  -- Create indexes for common queries
  CREATE INDEX [IX_data_eav_materials_project_fileid] ON [dbo].[data_eav_materials] ([project], [fileid]);
  CREATE INDEX [IX_data_eav_materials_id] ON [dbo].[data_eav_materials] ([id]);
  CREATE INDEX [IX_data_eav_materials_timestamp] ON [dbo].[data_eav_materials] ([timestamp]);
END