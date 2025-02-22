DROP TABLE IF EXISTS [dbo].[lca_data];
CREATE TABLE [dbo].[lca_data]
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
  [mat_kbob] VARCHAR(255),
  [gwp_absolute] DECIMAL(18,2),
  [gwp_relative] DECIMAL(18,2),
  [penr_absolute] DECIMAL(18,2),
  [penr_relative] DECIMAL(18,2),
  [ubp_absolute] DECIMAL(18,2),
  [ubp_relative] DECIMAL(18,2)
);