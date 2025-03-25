-- -- Add column statements for LCA data table

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'project' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [project] VARCHAR(255) NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'filename' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [filename] VARCHAR(255) NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'fileid' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [fileid] VARCHAR(255) NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'timestamp' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [timestamp] DATETIME2 NOT NULL DEFAULT GETDATE();
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'id' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'mat_kbob' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [mat_kbob] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'gwp_absolute' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [gwp_absolute] DECIMAL(18,2);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'gwp_relative' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [gwp_relative] DECIMAL(18,2);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'penr_absolute' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [penr_absolute] DECIMAL(18,2);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'penr_relative' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [penr_relative] DECIMAL(18,2);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'ubp_absolute' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [ubp_absolute] DECIMAL(18,2);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'ubp_relative' AND object_id = OBJECT_ID('lca_data'))
--     BEGIN
--   ALTER TABLE [dbo].[lca_data] ADD [ubp_relative] DECIMAL(18,2);
-- END

-- -- Add column statements for Cost data table
-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'project' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [project] VARCHAR(255) NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'filename' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [filename] VARCHAR(255) NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'fileid' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [fileid] VARCHAR(255) NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'timestamp' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [timestamp] DATETIME2 NOT NULL DEFAULT GETDATE();
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'id' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL DEFAULT '';
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'category' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [category] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'level' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [level] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'is_structural' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [is_structural] BIT;
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'fire_rating' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [fire_rating] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'ebkph' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [ebkph] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'ebkph_1' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [ebkph_1] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'ebkph_2' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [ebkph_2] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'ebkph_3' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [ebkph_3] VARCHAR(255);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'cost' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [cost] DECIMAL(18,2);
-- END

-- IF NOT EXISTS (SELECT *
-- FROM sys.columns
-- WHERE name = 'cost_unit' AND object_id = OBJECT_ID('cost_data'))
--     BEGIN
--   ALTER TABLE [dbo].[cost_data] ADD [cost_unit] DECIMAL(18,2);
-- END
