import pandas as pd
from sqlalchemy import create_engine, text
import urllib.parse
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Read both Parquet files
df_lca = pd.read_parquet('./export_lca_data.parquet')
df_cost = pd.read_parquet('./export_cost_data.parquet')

# Create connection string
params = urllib.parse.quote_plus(
    'Driver={ODBC Driver 18 for SQL Server};'
    'Server=sql-nhmzh-vis-dev.database.windows.net;'
    'Database=sqldb-nhmzh-vis-dev;'
    f'Uid={os.environ["DB_USERNAME"]};'
    f'Pwd={os.environ["DB_PASSWORD"]};'
    'Connection Timeout=30;'
)

# Create engine
engine = create_engine(f"mssql+pyodbc:///?odbc_connect={params}")

# SQL statements to create tables
create_lca_table = """
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
"""

create_cost_table = """
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
"""

# Execute CREATE TABLE statements
with engine.connect() as conn:
    conn.execute(text(create_lca_table))
    conn.execute(text(create_cost_table))
    conn.commit()

# Ensure column types match SQL table definitions
lca_columns = [
    'project', 'filename', 'fileid', 'timestamp', 'id',
    'ebkph', 'ebkph_1', 'ebkph_2', 'ebkph_3', 'mat_kbob',
    'gwp_absolute', 'gwp_relative', 'penr_absolute', 'penr_relative',
    'ubp_absolute', 'ubp_relative'
]

cost_columns = [
    'project', 'filename', 'fileid', 'timestamp', 'id',
    'ebkph', 'ebkph_1', 'ebkph_2', 'ebkph_3',
    'cost', 'cost_unit'
]

# Verify and select columns in correct order
df_lca = df_lca[lca_columns]
df_cost = df_cost[cost_columns]

# Write to SQL tables
df_lca.to_sql('lca_data', engine, if_exists='append', index=False, schema='dbo')
df_cost.to_sql('cost_data', engine, if_exists='append', index=False, schema='dbo')