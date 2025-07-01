# Data Model Documentation

The data model for 'NHM-Nachhaltigkeitsmonitoring' consists of 3 tables.

## Tables

**EAV style data modelling:**
(Element-Attribute-Value or 'Long-format data')
This is a useful method when you must/want to use a SQL database but you need to keep the properties per element flexible. It means that the properties/attributes of an element are spread across multiple rows. A pivot operation has to be performed to query all properties of an element.

Pro:
New properties can be added without altering the database schema
Con:
Less intuitive to query

### 1. Per Element Data

EAV style

This Table contains all properties of an element that have a simple 1:1 relationship with an element, for example:

- A wall is on level 1
- A column is structural

An element is uniquely identified by the combination of project, filename, id and timestamp.
This is because the same id could exist in different files and the same filename could be used in multiple projects.

A parameter/attribute is uniquely identified by the combination of project, filename, id, timestamp and param_name

### 2. Per Material-Layer Data

EAV style

This Table contains all information about a 'material-layer'. Material-uses have a many-to-1 relationship with an element. The 'layer' is represented by the 'sequence' column.
The term 'layer' is not exclusive to multi-layer elements like walls and floors, but can also mean a the use of a material in a window or door.
A 'material-layer' is uniquely identified by the combination of project, filename, id, sequence and timestamp and param_name.
The id is the id of the element the layer belongs to, param_name is usually the material name.

This table could be utilized to hold other entities that have a many-to-1 style relation with an element.

### 3. Metadata for Data Updates

This table contains information about the individual update events. An entry in this table is created when a new IFC model is uploaded to the platform (Which triggers a corresponding kafka message).

## Querying

See [Link-to-query-repo](this link) for an example query that can be used in powerbi to create a single wide-format table from the elements and materials table.

Such a query is the place to define what columns are desired in PowerBI. It basically appends the element and material rows - after the pivot - and, if possible, copies element attibutes onto material rows so that PowerBI can filter materials based on element properties.
