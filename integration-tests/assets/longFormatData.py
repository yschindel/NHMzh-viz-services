import pandas as pd
import numpy as np
from datetime import datetime, timedelta
import random
import os
import hashlib

# Read the real IDs from export_cost_data.parquet without modifying them
try:
    # Read the export_cost_data.parquet file
    cost_data_df = pd.read_parquet('export_cost_data.parquet')
    
    # Extract the unique IDs - keep them exactly as they are
    real_ids = cost_data_df['id'].unique().tolist()
    
    print(f"Using {len(real_ids)} real IDs from export_cost_data.parquet")
    
except Exception as e:
    print(f"Could not read real IDs from export_cost_data.parquet: {e}")
    print("Falling back to synthetic IDs")
    # Create synthetic IDs if we can't read real ones
    real_ids = [f"elem{i:03d}" for i in range(1, 101)]

# Sample data
projects = ['project-a', 'project-b', 'project-c']
files = ['file1', 'file2']

# Replace EBKPH codes with IFC categories
ifc_categories = [
    'IfcWall', 'IfcSlab', 'IfcBeam', 'IfcColumn', 'IfcRoof', 
    'IfcStair', 'IfcDoor', 'IfcWindow', 'IfcCurtainWall', 'IfcRailing'
]

materials = [
    'Concrete', 'Steel', 'Timber', 'Glass', 'Aluminum', 
    'Brick', 'Gypsum', 'Stone', 'Insulation', 'Carpet'
]
material_properties = {
    'Concrete': {'gwp': 120.5, 'penr': 450.75, 'ubp': 1200.25},
    'Steel': {'gwp': 210.4, 'penr': 525.6, 'ubp': 1580.3},
    'Timber': {'gwp': 40.2, 'penr': 210.5, 'ubp': 560.8},
    'Glass': {'gwp': 90.7, 'penr': 380.3, 'ubp': 980.5},
    'Aluminum': {'gwp': 180.3, 'penr': 490.2, 'ubp': 1480.6},
    'Brick': {'gwp': 85.6, 'penr': 310.4, 'ubp': 920.7},
    'Gypsum': {'gwp': 35.8, 'penr': 160.3, 'ubp': 450.2},
    'Stone': {'gwp': 65.3, 'penr': 240.8, 'ubp': 780.4},
    'Insulation': {'gwp': 45.3, 'penr': 89.2, 'ubp': 320.5},
    'Carpet': {'gwp': 30.1, 'penr': 120.4, 'ubp': 280.6}
}

# Generate data
element_rows = []
material_rows = []
element_count = len(real_ids)

# Generate base timestamp for each project/file
base_timestamps = {}
start_date = datetime(2024, 1, 1)
for project in projects:
    for file in files:
        if random.random() < 0.7:  # Not all project/file combinations exist
            base_timestamps[f"{project}/{file}"] = start_date + timedelta(
                days=random.randint(0, 180),
                hours=random.randint(0, 23),
                minutes=random.randint(0, 59)
            )

# Generate elements with properties
for i, element_id in enumerate(real_ids):
    # Choose project/file randomly
    project_file_key = random.choice(list(base_timestamps.keys()))
    project, file = project_file_key.split('/')
    fileid = project_file_key
    timestamp = base_timestamps[project_file_key]
    
    # Element properties
    ifc_category = random.choice(ifc_categories)
    ebkph_1 = "E"  # First part is always "E"
    ebkph_2 = str(random.randint(1, 9)).zfill(2)  # Format as 01, 02, etc.
    ebkph_3 = str(random.randint(1, 9)).zfill(2)  # Format as 01, 02, etc.
    full_ebkph = f"{ebkph_1}.{ebkph_2}.{ebkph_3}"  # Format as E.01.01
    cost = round(random.uniform(1000, 15000), 2)
    cost_unit = round(random.uniform(50, 500), 2)  # Cost per unit for the element
    
    # Add element properties
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'ifcType',
        'param_value_string': ifc_category,
        'param_type': 'string'
    })
    
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'level',
        'param_value_string': ebkph_2,
        'param_type': 'string'
    })

    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'is_structural',
        'param_value_boolean': random.choice([True, False]),
        'param_type': 'boolean'
    })
    
    # Add full EBKPH
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'ebkph',
        'param_value_string': full_ebkph,
        'param_type': 'string'
    })
    
    # Add split EBKPH fields
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'ebkph_1',
        'param_value_string': ebkph_1,
        'param_type': 'string'
    })
    
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'ebkph_2',
        'param_value_string': ebkph_2,
        'param_type': 'string'
    })
    
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'ebkph_3',
        'param_value_string': ebkph_3,
        'param_type': 'string'
    })
    
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'cost',
        'param_value_number': cost,
        'param_type': 'number'
    })
    
    element_rows.append({
        'project': project,
        'filename': file,
        'fileid': fileid,
        'timestamp': timestamp,
        'id': element_id,
        'param_name': 'cost_unit',
        'param_value_number': cost_unit,
        'param_type': 'number'
    })
    
    # Add material layers (1-5 random materials)
    material_count = random.randint(1, 5)
    selected_materials = random.sample(materials, material_count)
    
    for layer_idx, material in enumerate(selected_materials):
        # Add material name
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'mat_kbob',
            'param_value_string': material,
            'param_type': 'string'
        })
                
        # Add GWP absolute (fully random, not tied to material)
        gwp_abs = round(random.uniform(30, 300), 2)
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'gwp_absolute',
            'param_value_number': gwp_abs,
            'param_type': 'number'
        })
        
        # Add GWP relative (completely independent random value)
        gwp_rel = round(random.uniform(1, 50), 2)
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'gwp_relative',
            'param_value_number': gwp_rel,
            'param_type': 'number'
        })
        
        # Add PENR absolute
        penr_abs = round(random.uniform(100, 600), 2)
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'penr_absolute',
            'param_value_number': penr_abs,
            'param_type': 'number'
        })
        
        # Add PENR relative
        penr_rel = round(random.uniform(5, 100), 2)
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'penr_relative',
            'param_value_number': penr_rel,
            'param_type': 'number'
        })
        
        # Add UBP absolute
        ubp_abs = round(random.uniform(300, 2000), 2)
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'ubp_absolute',
            'param_value_number': ubp_abs,
            'param_type': 'number'
        })
        
        # Add UBP relative
        ubp_rel = round(random.uniform(10, 300), 2)
        material_rows.append({
            'project': project,
            'filename': file,
            'fileid': fileid,
            'timestamp': timestamp,
            'id': element_id,
            'sequence': layer_idx,
            'param_name': 'ubp_relative',
            'param_value_number': ubp_rel,
            'param_type': 'number'
        })

# Create DataFrames
elements_df = pd.DataFrame(element_rows)
materials_df = pd.DataFrame(material_rows)

# Sort the data
elements_df = elements_df.sort_values(['project', 'fileid', 'timestamp', 'id', 'param_name'])
materials_df = materials_df.sort_values(['project', 'fileid', 'timestamp', 'id', 'sequence', 'param_name'])

# Save to parquet files
elements_df.to_parquet('data_eav_elements.parquet', index=False)
materials_df.to_parquet('data_eav_materials.parquet', index=False)

print(f"Generated {len(elements_df)} rows in EAV format for elements")
print(f"Generated {len(materials_df)} rows in EAV format for materials")
print(f"Files saved: data_eav_elements.parquet and data_eav_materials.parquet")