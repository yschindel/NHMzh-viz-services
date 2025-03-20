#!/bin/bash

# Check if input file is provided
if [ -z "$1" ]; then
    echo "Usage: ./postprocess-gltf.sh <input-file.glb>"
    exit 1
fi

# Get input file and create output filename
INPUT_FILE=$1
FILENAME=$(basename "$INPUT_FILE" .glb)
TEMP_FILE="${FILENAME}_temp.glb"
OUTPUT_FILE="out.glb"

# Apply transformations one at a time
gltf-transform dedup "$INPUT_FILE" "$TEMP_FILE" 
gltf-transform instance "$TEMP_FILE" "$OUTPUT_FILE" 

