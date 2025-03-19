import ifcopenshell
import ifcopenshell.geom
import ifcopenshell.util
import multiprocessing
import gc
import time
import os
import sys

filename = "test.ifc"
gltfFileName = "test.glb"

# Check if the IFC file exists
if not os.path.exists(filename):
    print(f"Error: IFC file '{filename}' not found!")
    sys.exit(1)

# Print ifcopenshell version for debugging
print(f"Using IfcOpenShell version: {ifcopenshell.version}")

try:
    f = ifcopenshell.open(filename)
    print(f"Successfully opened IFC file with {len(f.by_type('IfcProduct'))} products")
    
    ifcopenshell.ifcopenshell_wrapper.turn_off_detailed_logging()
    ifcopenshell.ifcopenshell_wrapper.set_log_format_json()

    # Create settings object with more detailed configuration
    settings = ifcopenshell.geom.settings()
    settings.set(settings.APPLY_DEFAULT_MATERIALS, True)
    settings.set(settings.USE_WORLD_COORDS, True)
    settings.set(settings.INCLUDE_CURVES, False)
    
    # Print available settings for debugging
    print("Available settings:")
    for setting_name in dir(settings):
        if setting_name.isupper() and not setting_name.startswith('_'):
            print(f"  {setting_name}")
    
    # Create serializer settings
    serializer_settings = ifcopenshell.geom.serializer_settings()
    
    # Check if there are existing temp files and remove them
    for temp_file in [f"{gltfFileName}.indices.tmp", f"{gltfFileName}.vertices.tmp", gltfFileName]:
        if os.path.exists(temp_file):
            os.remove(temp_file)
            print(f"Removed existing file: {temp_file}")
    
    # Create GLTF serializer
    print("Creating GLTF serializer...")
    serializer = ifcopenshell.geom.serializers.gltf(gltfFileName, settings, serializer_settings)
    
    print("Writing header...")
    serializer.writeHeader()
    
    # Count elements for progress tracking
    product_count = len(f.by_type('IfcProduct'))
    print(f"Processing {product_count} products...")
    
    # Process elements with progress reporting
    processed = 0
    for progress, elem in ifcopenshell.geom.iterate(
        settings, f, with_progress=True,
        exclude=("IfcSpace", "IfcOpeningElement"),
        num_threads=multiprocessing.cpu_count()
    ):
        processed += 1
        if processed % 10 == 0 or processed == 1:
            print(f"Processed {processed}/{product_count} elements ({progress*100:.1f}%)")
        
        result = serializer.write(elem)
        if not result:
            print(f"Warning: Failed to write element {elem.id()}")
    
    print("Finalizing serializer...")
    serializer.finalize()
    
    print("Cleaning up serializer...")
    del serializer
    gc.collect()
    time.sleep(2)  # Increased wait time
    
    # Verify output files
    for check_file in [gltfFileName, f"{gltfFileName}.indices.tmp", f"{gltfFileName}.vertices.tmp"]:
        if os.path.exists(check_file):
            size = os.path.getsize(check_file)
            print(f"File: {check_file}, Size: {size} bytes")
        else:
            print(f"File not found: {check_file}")
    
    # Try to manually move tmp files into the final glb if needed
    # This is a workaround based on the issue described in the forum
    if os.path.getsize(gltfFileName) == 0 and os.path.exists(f"{gltfFileName}.indices.tmp") and os.path.exists(f"{gltfFileName}.vertices.tmp"):
        print("Attempting manual file processing...")
        
        # Try a different approach by keeping the serializer alive longer
        print("Creating a new serializer and trying again...")
        new_serializer = ifcopenshell.geom.serializers.gltf(gltfFileName + ".new", settings, serializer_settings)
        new_serializer.writeHeader()
        
        for progress, elem in ifcopenshell.geom.iterate(
            settings, f, with_progress=True,
            exclude=("IfcSpace", "IfcOpeningElement"),
            num_threads=1  # Try with a single thread
        ):
            new_serializer.write(elem)
        
        new_serializer.finalize()
        
        # Keep the serializer alive for a bit longer
        time.sleep(3)
        print(f"New file size: {os.path.getsize(gltfFileName + '.new')} bytes")
        
        # Only now delete the serializer
        del new_serializer
        gc.collect()
        time.sleep(2)

except Exception as e:
    print(f"Error occurred: {e}")
    import traceback
    traceback.print_exc()