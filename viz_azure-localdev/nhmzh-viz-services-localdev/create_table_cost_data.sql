DROP TABLE IF EXISTS [dbo].[cost_data];
CREATE TABLE [dbo].[cost_data]
(
  [record_id] BIGINT IDENTITY(1,1) PRIMARY KEY,
  [project] VARCHAR(255) NOT NULL,
  [filename] VARCHAR(255) NOT NULL,
  [fileid] VARCHAR(255) NOT NULL,
  [timestamp] VARCHAR(255) NOT NULL,
  [id] VARCHAR(255) NOT NULL,
  [ebkph] VARCHAR(255),
  [ebkph_1] VARCHAR(255),
  [ebkph_2] VARCHAR(255),
  [ebkph_3] VARCHAR(255),
  [mat_kbob] VARCHAR(255),
  [cost] DECIMAL(18,2),
  [cost_unit] DECIMAL(18,2),
);