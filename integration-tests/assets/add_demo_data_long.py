import pandas as pd
from sqlalchemy import create_engine, text
import urllib.parse
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Read the EAV format parquet file
df_eav = pd.read_parquet('./denormalized_elements_materials_eav.parquet')

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

# SQL statement to create the EAV table
create_eav_table = """
DROP TABLE IF EXISTS [dbo].[data_eav];
CREATE TABLE [dbo].[data_eav]
(
  [project] VARCHAR(255) NOT NULL,
  [filename] VARCHAR(255) NOT NULL,
  [fileid] VARCHAR(255) NOT NULL,
  [timestamp] DATETIME2 NOT NULL,
  [id] VARCHAR(255) COLLATE Latin1_General_CS_AS NOT NULL,
  [group] VARCHAR(255) NOT NULL,
  [sequence] INT NOT NULL,
  [param_name] VARCHAR(255) NOT NULL,
  [param_value_string] NVARCHAR(MAX) NOT NULL,
  [param_value_number] DECIMAL(18,2) NOT NULL,
  [param_value_boolean] BIT NOT NULL,
  [param_type] VARCHAR(50) NOT NULL,
  CONSTRAINT [PK_elements_materials_eav] PRIMARY KEY 
    ([project], [filename], [timestamp], [id], [group], [sequence], [param_name])
);
"""

# Execute CREATE TABLE statement
with engine.connect() as conn:
    conn.execute(text(create_eav_table))
    conn.commit()

# Before inserting data
print(f"Total records in DataFrame: {len(df_eav)}")
print(f"Unique combinations: {len(df_eav.groupby(['project', 'filename', 'timestamp', 'id', 'group', 'sequence', 'param_name']).size())}")

# Get a list of IDs from the DataFrame for comparison
original_ids = set(df_eav['id'].unique())
print(f"Number of unique IDs in DataFrame: {len(original_ids)}")

# Before the to_sql, ensure timestamp format is maintained
if pd.api.types.is_datetime64_any_dtype(df_eav['timestamp']):
    df_eav['timestamp'] = df_eav['timestamp'].dt.strftime('%Y-%m-%dT%H:%M:%S.%fZ')

# Write to SQL table
df_eav.to_sql('data_eav', engine, if_exists='append', index=False, schema='dbo')

# After insertion, check what got inserted
with engine.connect() as conn:
    # Count total records
    result = conn.execute(text("SELECT COUNT(*) AS total_records FROM [dbo].[data_eav]"))
    total_records = result.fetchone()[0]
    print(f"Records in database: {total_records}")
    
    # Get IDs that made it into the database
    result = conn.execute(text("SELECT DISTINCT id FROM [dbo].[data_eav]"))
    db_ids = set([row[0] for row in result])
    print(f"Number of unique IDs in database: {len(db_ids)}")
    
    # Find missing IDs
    missing_ids = original_ids - db_ids
    if missing_ids:
        print(f"Found {len(missing_ids)} IDs missing from database!")
        print(f"Sample of missing IDs: {list(missing_ids)[:5]}")
        
        # Check if these IDs appear in original data
        missing_records = df_eav[df_eav['id'].isin(missing_ids)]
        print(f"Sample of missing records from original data:")
        print(missing_records.head())

    # Show sample of data in database
    result = conn.execute(text("""
        SELECT TOP 10 * 
        FROM [dbo].[data_eav]
        ORDER BY project, filename, id, group, sequence, param_name
    """))
    print("\nSample data in database:")
    for row in result:
        print(row)