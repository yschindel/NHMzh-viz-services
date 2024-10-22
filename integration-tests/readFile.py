def read_file_to_buffer(file_path: str) -> bytes:
    """
    Reads a file from disk into a bytes buffer.
    
    Args:
    file_path (str): The path to the file to be read.
    
    Returns:
    bytes: The contents of the file as a bytes object.
    
    Raises:
    FileNotFoundError: If the specified file does not exist.
    IOError: If there's an error reading the file.
    """
    try:
        with open(file_path, 'rb') as file:
            return file.read()
    except FileNotFoundError:
        print(f"Error: The file '{file_path}' was not found.")
        raise
    except IOError as e:
        print(f"Error reading the file '{file_path}': {str(e)}")
        raise

