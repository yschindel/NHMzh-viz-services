DROP TABLE IF EXISTS [dbo].[cost_data];
CREATE TABLE [dbo].[cost_data]
(
  [project] VARCHAR(255) NOT NULL,
  [filename] VARCHAR(255) NOT NULL,
  [fileid] VARCHAR(255) NOT NULL,
  [timestamp] DATETIME2 NOT NULL,
  [id] VARCHAR(255) NOT NULL,
  [ebkph] VARCHAR(255),
  [ebkph_1] VARCHAR(255),
  [ebkph_2] VARCHAR(255),
  [ebkph_3] VARCHAR(255),
  [cost] DECIMAL(18,2),
  [cost_unit] DECIMAL(18,2),
  CONSTRAINT [UQ_cost_data_id_filename_timestamp_project] UNIQUE ([id], [filename], [timestamp], [project])
);