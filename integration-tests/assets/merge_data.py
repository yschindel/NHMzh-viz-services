import duckdb
import os
import time
from datetime import datetime, timedelta

start_time = time.time()

# Get the directory where the script is located
script_dir = os.path.dirname(os.path.abspath(__file__))
# Change the working directory to the script's directory
os.chdir(script_dir)

# get all csv files in the current directory
csv_files = [f for f in os.listdir('.') if f.endswith('.csv')]

# remove all non utf-8 characters from the csv files
for file in csv_files:
    try:
        # First try to read with utf-8
        with open(file, 'r', encoding='utf-8') as f:
            content = f.read()
    except UnicodeDecodeError:
        # If that fails, try with ISO-8859-1 (common for German text)
        with open(file, 'r', encoding='ISO-8859-1') as f:
            content = f.read()
    
    # Write back as UTF-8
    with open(file, 'w', encoding='utf-8') as f:
        f.write(content)


# timestamps
# create one timestamp for each file, use the last couple of days
timestamps = [(datetime.now() - timedelta(days=i)).strftime('%Y-%m-%dT%H:%M:%S.%fZ') for i in range(len(csv_files))]

duckdb.sql("""CREATE TABLE data (
id VARCHAR,
lca BOOL,
ebkph VARCHAR, 
cost DOUBLE, 
mat_kbob VARCHAR,
gwp_absolute DOUBLE,
gwp_relative DOUBLE,
penr_absolute DOUBLE,
penr_relative DOUBLE,
ubp_absolute DOUBLE,
ubp_relative DOUBLE,
cost_unit DOUBLE, 
timestamp VARCHAR, 
)""")

for file, timestamp in zip(csv_files, timestamps):
    duckdb.sql(f"INSERT INTO data SELECT *, VARCHAR '{timestamp}' as timestamp FROM read_csv('{file}')")


# log all table names
print(duckdb.sql("SHOW TABLES"))

# save the database to a parquet file
duckdb.sql("COPY (SELECT * FROM data) TO 'data.parquet' (FORMAT 'parquet');")

end_time = time.time()
print(f"Time taken: {end_time - start_time} seconds")