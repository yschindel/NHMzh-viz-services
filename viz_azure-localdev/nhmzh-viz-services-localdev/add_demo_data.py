import pandas as pd
from sqlalchemy import create_engine
import urllib.parse
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Read the CSV
df = pd.read_csv('./export_data.csv')

# Create connection string (replace with your details)
params = urllib.parse.quote_plus(
    'Driver={ODBC Driver 18 for SQL Server};'
    'Server=sql-nhmzh-vis-dev.database.windows.net;'
    'Database=sqldb-nhmzh-vis-dev;'
    f'Uid={os.environ["DB_USERNAME"]};'
    f'Pwd={os.environ["DB_PASSWORD"]};'
)

# Create engine
engine = create_engine(f"mssql+pyodbc:///?odbc_connect={params}")

# Write to SQL
df.to_sql('project_data', engine, if_exists='append', index=False, schema='dbo')