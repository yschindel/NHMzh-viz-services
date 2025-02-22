DROP TABLE IF EXISTS lca_data;
DROP TABLE IF EXISTS cost_data;

CREATE TABLE lca_data (
    project VARCHAR,
    filename VARCHAR,
    fileid VARCHAR,
    timestamp VARCHAR,
    id VARCHAR,
    ebkph VARCHAR,
    ebkph_1 VARCHAR,
    ebkph_2 VARCHAR,
    ebkph_3 VARCHAR,
    mat_kbob VARCHAR,
    gwp_absolute FLOAT,
    gwp_relative FLOAT,
    penr_absolute FLOAT,
    penr_relative FLOAT,
    ubp_absolute FLOAT,
    ubp_relative FLOAT,
);

CREATE TABLE cost_data (
    project VARCHAR,
    filename VARCHAR,
    fileid VARCHAR,
    timestamp VARCHAR,
    id VARCHAR,
    ebkph VARCHAR,
    ebkph_1 VARCHAR,
    ebkph_2 VARCHAR,
    ebkph_3 VARCHAR,
    cost FLOAT,
    cost_unit FLOAT,
    CONSTRAINT UQ_cost_data_id_filename_timestamp UNIQUE (id, filename, timestamp, project)
);

INSERT INTO lca_data 
  SELECT 
      'juch-areal' AS project,
      'file1' AS filename,
      'juch-areal/file1' AS fileid,
      '2024-10-25T16:36:04.986158Z' AS timestamp,
      id, 
      ebkph,
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data1.csv';

INSERT INTO cost_data 
SELECT 
    'juch-areal' AS project,
    'file1' AS filename,
    'juch-areal/file1' AS fileid,
    '2024-10-25T16:36:04.986158Z' AS timestamp,
    id, 
    ebkph,
    SUBSTRING(ebkph, 1, 1) AS ebkph_1,
    SUBSTRING(ebkph, 2, 2) AS ebkph_2,
    SUBSTRING(ebkph, 5, 2) AS ebkph_3,
    cost, 
    cost_unit
FROM (
    SELECT DISTINCT ON (id) *
    FROM 'data1.csv'
) t;

INSERT INTO lca_data 
  SELECT 
      'juch-areal' AS project,
      'file1' AS filename,
      'juch-areal/file1' AS fileid,
      '2024-11-25T16:36:04.986158Z' AS timestamp,
      id, 
      ebkph,
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data2.csv';

INSERT INTO cost_data 
SELECT 
    'juch-areal' AS project,
    'file1' AS filename,
    'juch-areal/file1' AS fileid,
    '2024-11-25T16:36:04.986158Z' AS timestamp,
    id, 
    ebkph,
    SUBSTRING(ebkph, 1, 1) AS ebkph_1,
    SUBSTRING(ebkph, 2, 2) AS ebkph_2,
    SUBSTRING(ebkph, 5, 2) AS ebkph_3,
    cost, 
    cost_unit
FROM (
    SELECT DISTINCT ON (id) *
    FROM 'data2.csv'
) t;


INSERT INTO lca_data 
    SELECT 
      'juch-areal' AS project,
      'file2' AS filename,
      'juch-areal/file2' AS fileid,
      '2024-10-25T19:36:04.986158Z' AS timestamp,
      id, 
      ebkph,
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data3.csv';

INSERT INTO cost_data 
SELECT 
    'juch-areal' AS project,
    'file2' AS filename,
    'juch-areal/file2' AS fileid,
    '2024-10-25T19:36:04.986158Z' AS timestamp,
    id, 
    ebkph,
    SUBSTRING(ebkph, 1, 1) AS ebkph_1,
    SUBSTRING(ebkph, 2, 2) AS ebkph_2,
    SUBSTRING(ebkph, 5, 2) AS ebkph_3,
    cost, 
    cost_unit
FROM (
    SELECT DISTINCT ON (id) *
    FROM 'data3.csv'
) t;


INSERT INTO lca_data 
  SELECT 
      'other-project' AS project,
      'file3' AS filename,
      'other-project/file3' AS fileid,
      '2024-10-23T16:36:04.986158Z' AS timestamp,
      id, 
      ebkph,
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data4.csv';

INSERT INTO cost_data 
SELECT 
    'other-project' AS project,
    'file3' AS filename,
    'other-project/file3' AS fileid,
    '2024-10-23T16:36:04.986158Z' AS timestamp,
    id, 
    ebkph,
    SUBSTRING(ebkph, 1, 1) AS ebkph_1,
    SUBSTRING(ebkph, 2, 2) AS ebkph_2,
    SUBSTRING(ebkph, 5, 2) AS ebkph_3,
    cost, 
    cost_unit
FROM (
    SELECT DISTINCT ON (id) *
    FROM 'data4.csv'
) t;

COPY lca_data TO 'export_lca_data.csv' (HEADER, DELIMITER ',');
COPY lca_data TO 'export_lca_data.parquet';

COPY cost_data TO 'export_cost_data.csv' (HEADER, DELIMITER ',');
COPY cost_data TO 'export_cost_data.parquet';