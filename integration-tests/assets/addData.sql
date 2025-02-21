DROP TABLE IF EXISTS data;

CREATE TABLE data (
    project VARCHAR,
    filename VARCHAR,
    fileid VARCHAR,
    timestamp VARCHAR,
    id VARCHAR,
    lca BOOLEAN,
    ebkph VARCHAR,
    ebkph_1 VARCHAR,
    ebkph_2 VARCHAR,
    ebkph_3 VARCHAR,
    cost FLOAT,
    cost_unit FLOAT,
    mat_kbob VARCHAR,
    gwp_absolute FLOAT,
    gwp_relative FLOAT,
    penr_absolute FLOAT,
    penr_relative FLOAT,
    ubp_absolute FLOAT,
    ubp_relative FLOAT
);

INSERT INTO data 
  SELECT 
      'juch-areal' AS project,
      'file1' AS filename,
      'juch-areal/file1' AS fileid,
      '2024-10-25T16:36:04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph,
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      cost, 
      cost_unit,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data1.parquet';

INSERT INTO data 
  SELECT 
      'juch-areal' AS project,
      'file1' AS filename,
      'juch-areal/file1' AS fileid,
      '2024-11-25T16:36:04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph,
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      cost, 
      cost_unit,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data2.parquet';

INSERT INTO data 
  SELECT 
      'juch-areal' AS project,
      'file2' AS filename,
      'juch-areal/file2' AS fileid,
      '2024-10-25T16:36:04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph, 
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      cost, 
      cost_unit,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data3.parquet';

INSERT INTO data 
  SELECT 
      'other-project' AS project,
      'file3' AS filename,
      'other-project/file3' AS fileid,
      '2024-10-23T16:36:04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph, 
      SUBSTRING(ebkph, 1, 1) AS ebkph_1,
      SUBSTRING(ebkph, 2, 2) AS ebkph_2,
      SUBSTRING(ebkph, 5, 2) AS ebkph_3,
      cost, 
      cost_unit,
      mat_kbob, 
      gwp_absolute, 
      gwp_relative, 
      penr_absolute, 
      penr_relative, 
      ubp_absolute, 
      ubp_relative
  FROM 'data4.parquet';

COPY data TO 'export_data.csv' (HEADER, DELIMITER ',');
COPY data TO 'export_data.parquet';