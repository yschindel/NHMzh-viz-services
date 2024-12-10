DROP TABLE IF EXISTS data;

CREATE TABLE data (
    project VARCHAR,
    filename VARCHAR,
    fileId VARCHAR,
    timestamp VARCHAR,
    id VARCHAR,
    lca BOOLEAN,
    ebkph VARCHAR,
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
      'juch-areal/file1' AS fileId,
      '2024-10-25T16_36_04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph, 
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
      'juch-areal/file1' AS fileId,
      '2024-11-25T16_36_04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph, 
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
      'juch-areal/file2' AS fileId,
      '2024-10-25T16_36_04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph, 
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
      'project1' AS project,
      'file3' AS filename,
      'project1/file3' AS fileId,
      '2024-10-23T16_36_04.986158Z' AS timestamp,
      id, 
      lca, 
      ebkph, 
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

COPY data TO 'data.csv' (HEADER, DELIMITER ',');