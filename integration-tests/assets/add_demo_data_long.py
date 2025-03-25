import pandas as pd
from sqlalchemy import create_engine, text
import urllib.parse
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Read the EAV format parquet files
elements_df = pd.read_parquet('./data_eav_elements.parquet')
materials_df = pd.read_parquet('./data_eav_materials.parquet')

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

# SQL statement to create the EAV elements table
create_eav_elements_table = """
DROP TABLE IF EXISTS [dbo].[data_eav_elements];
CREATE TABLE [dbo].[data_eav_elements]
(
  [project] VARCHAR(255) NOT NULL,
  [filename] VARCHAR(255) NOT NULL,
  [fileid] VARCHAR(255) NOT NULL,
  [timestamp] DATETIME2 NOT NULL,
  [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
  [param_name] VARCHAR(255) NOT NULL,
  [param_value_string] NVARCHAR(MAX),
  [param_value_number] DECIMAL(18,2),
  [param_value_boolean] BIT,
  [param_value_date] DATETIME2,
  [param_type] VARCHAR(50) NOT NULL,
  CONSTRAINT [PK_data_eav_elements] PRIMARY KEY 
    ([project], [filename], [timestamp], [id], [param_name])
);
"""

# SQL statement to create the EAV materials table
create_eav_materials_table = """
DROP TABLE IF EXISTS [dbo].[data_eav_materials];
CREATE TABLE [dbo].[data_eav_materials]
(
  [project] VARCHAR(255) NOT NULL,
  [filename] VARCHAR(255) NOT NULL,
  [fileid] VARCHAR(255) NOT NULL,
  [timestamp] DATETIME2 NOT NULL,
  [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
  [sequence] INT NOT NULL,
  [param_name] VARCHAR(255) NOT NULL,
  [param_value_string] NVARCHAR(MAX),
  [param_value_number] DECIMAL(18,2),
  [param_value_boolean] BIT,
  [param_value_date] DATETIME2,
  [param_type] VARCHAR(50) NOT NULL,
  CONSTRAINT [PK_data_eav_materials] PRIMARY KEY 
    ([project], [filename], [timestamp], [id], [sequence], [param_name])
);
"""

# Execute CREATE TABLE statements
with engine.connect() as conn:
    conn.execute(text(create_eav_elements_table))
    conn.execute(text(create_eav_materials_table))
    conn.commit()

# Before inserting data
print(f"Total records in elements DataFrame: {len(elements_df)}")
print(f"Total records in materials DataFrame: {len(materials_df)}")

# Get a list of IDs from both DataFrames for comparison
element_ids = set(elements_df['id'].unique())
material_ids = set(materials_df['id'].unique())
print(f"Number of unique IDs in elements DataFrame: {len(element_ids)}")
print(f"Number of unique IDs in materials DataFrame: {len(material_ids)}")

# Before the to_sql, ensure timestamp format is maintained
if pd.api.types.is_datetime64_any_dtype(elements_df['timestamp']):
    elements_df['timestamp'] = elements_df['timestamp'].dt.strftime('%Y-%m-%dT%H:%M:%S.%fZ')
if pd.api.types.is_datetime64_any_dtype(materials_df['timestamp']):
    materials_df['timestamp'] = materials_df['timestamp'].dt.strftime('%Y-%m-%dT%H:%M:%S.%fZ')

# Write to SQL tables
elements_df.to_sql('data_eav_elements', engine, if_exists='append', index=False, schema='dbo')
materials_df.to_sql('data_eav_materials', engine, if_exists='append', index=False, schema='dbo')

# After insertion, check what got inserted
with engine.connect() as conn:
    # Count total records in each table
    result = conn.execute(text("SELECT COUNT(*) AS total_records FROM [dbo].[data_eav_elements]"))
    elements_records = result.fetchone()[0]
    result = conn.execute(text("SELECT COUNT(*) AS total_records FROM [dbo].[data_eav_materials]"))
    materials_records = result.fetchone()[0]
    print(f"\nRecords in elements table: {elements_records}")
    print(f"Records in materials table: {materials_records}")
    
    # Get IDs that made it into the database
    result = conn.execute(text("SELECT DISTINCT id FROM [dbo].[data_eav_elements]"))
    db_element_ids = set([row[0] for row in result])
    result = conn.execute(text("SELECT DISTINCT id FROM [dbo].[data_eav_materials]"))
    db_material_ids = set([row[0] for row in result])
    print(f"\nNumber of unique IDs in elements table: {len(db_element_ids)}")
    print(f"Number of unique IDs in materials table: {len(db_material_ids)}")
    
    # Find missing IDs
    missing_element_ids = element_ids - db_element_ids
    missing_material_ids = material_ids - db_material_ids
    
    if missing_element_ids:
        print(f"\nFound {len(missing_element_ids)} element IDs missing from database!")
        print(f"Sample of missing element IDs: {list(missing_element_ids)[:5]}")
    
    if missing_material_ids:
        print(f"\nFound {len(missing_material_ids)} material IDs missing from database!")
        print(f"Sample of missing material IDs: {list(missing_material_ids)[:5]}")

    # Show sample of data in both tables
    print("\nSample data in elements table:")
    result = conn.execute(text("""
        SELECT TOP 5 * 
        FROM [dbo].[data_eav_elements]
        ORDER BY project, filename, id, param_name
    """))
    for row in result:
        print(row)
        
    print("\nSample data in materials table:")
    result = conn.execute(text("""
        SELECT TOP 5 * 
        FROM [dbo].[data_eav_materials]
        ORDER BY project, filename, id, sequence, param_name
    """))
    for row in result:
        print(row)