CREATE TABLE [dbo].[project_data]
(
  [record_id] BIGINT IDENTITY(1,1) PRIMARY KEY,
  [Id] VARCHAR(255) NOT NULL,
  [project] VARCHAR(255) NOT NULL,
  [filename] VARCHAR(255) NOT NULL,
  [timestamp] VARCHAR(255) NOT NULL,
  [category] VARCHAR(255),
  [cost] DECIMAL(18,2),
  [cost_unit] DECIMAL(18,2),
  [material_kbob] VARCHAR(255),
  [gwp_absolute] DECIMAL(18,2),
  [gwp_relative] DECIMAL(18,2),
  [penr_absolute] DECIMAL(18,2),
  [penr_relative] DECIMAL(18,2),
  [ubp_absolute] DECIMAL(18,2),
  [ubp_relative] DECIMAL(18,2),
  CONSTRAINT UQ_project_data UNIQUE (Id, project, filename)
);