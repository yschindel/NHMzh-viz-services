-- This file can be used to add columns to the database after the database has been created. For example, if you need to add a column to a table that already has data, you can use this file to add the column.

--drop column fileid if exists
ALTER TABLE data_eav_elements
DROP COLUMN fileid;

ALTER TABLE data_eav_materials
DROP COLUMN fileid;
