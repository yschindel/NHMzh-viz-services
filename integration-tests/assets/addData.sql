DROP TABLE IF EXISTS data;

CREATE TABLE data (
    project VARCHAR,
    filename VARCHAR,
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
      'project1' AS project,
      'file1' AS filename,
      '2024-11-21T15:53:30.81691Z' AS timestamp,
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
      'project1' AS project,
      'file1' AS filename,
      '2024-11-21T15:54:26.347154Z' AS timestamp,
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
      'project1' AS project,
      'file3' AS filename,
      '2024-11-21T15:53:30.816913Z' AS timestamp,
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

COPY data TO 'data.csv' (HEADER, DELIMITER ',');