
-- Add column statements for LCA data table
IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'fileid' AND object_id = OBJECT_ID('lca_data'))
    BEGIN
  ALTER TABLE [dbo].[lca_data] ADD [fileid] VARCHAR(255) NOT NULL DEFAULT '';
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph' AND object_id = OBJECT_ID('lca_data'))
    BEGIN
  ALTER TABLE [dbo].[lca_data] ADD [ebkph] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph_1' AND object_id = OBJECT_ID('lca_data'))
    BEGIN
  ALTER TABLE [dbo].[lca_data] ADD [ebkph_1] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph_2' AND object_id = OBJECT_ID('lca_data'))
    BEGIN
  ALTER TABLE [dbo].[lca_data] ADD [ebkph_2] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph_3' AND object_id = OBJECT_ID('lca_data'))
    BEGIN
  ALTER TABLE [dbo].[lca_data] ADD [ebkph_3] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'mat_kbob' AND object_id = OBJECT_ID('lca_data'))
    BEGIN
  ALTER TABLE [dbo].[lca_data] ADD [mat_kbob] VARCHAR(255);
END

-- Add column statements for Cost data table
IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'fileid' AND object_id = OBJECT_ID('cost_data'))
    BEGIN
  ALTER TABLE [dbo].[cost_data] ADD [fileid] VARCHAR(255) NOT NULL DEFAULT '';
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph' AND object_id = OBJECT_ID('cost_data'))
    BEGIN
  ALTER TABLE [dbo].[cost_data] ADD [ebkph] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph_1' AND object_id = OBJECT_ID('cost_data'))
    BEGIN
  ALTER TABLE [dbo].[cost_data] ADD [ebkph_1] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph_2' AND object_id = OBJECT_ID('cost_data'))
    BEGIN
  ALTER TABLE [dbo].[cost_data] ADD [ebkph_2] VARCHAR(255);
END

IF NOT EXISTS (SELECT *
FROM sys.columns
WHERE name = 'ebkph_3' AND object_id = OBJECT_ID('cost_data'))
    BEGIN
  ALTER TABLE [dbo].[cost_data] ADD [ebkph_3] VARCHAR(255);
END