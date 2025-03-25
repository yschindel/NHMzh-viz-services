package server

import (
	"log"
	"reflect"
	"strings"
	"time"
)

// MessageProcessor handles the processing of incoming messages
type MessageProcessor struct{}

// NewMessageProcessor creates a new message processor
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{}
}

// GenerateFileID creates a file ID by concatenating project and filename with a '/' separator.
// This provides a unique identifier for files within projects that will be stored in the database.
//
// The fileid format is "project/filename" where:
// - project: the project identifier
// - filename: the filename with any .ifc extension removed
//
// For example, if project="building123" and filename="model.ifc",
// the resulting fileid will be "building123/model"
func (p *MessageProcessor) GenerateFileID(project, filename string) string {
	// Remove any .ifc extension from the filename
	cleanFilename := strings.TrimSuffix(filename, ".ifc")

	// Concatenate project and filename with a '/' separator
	return project + "/" + cleanFilename
}

// SplitEbkphCode splits an eBKP-H code (like "C.02.95") into its individual components.
// The function returns the components in an array: [ebkph_1, ebkph_2, ebkph_3]
// For example:
// - "C.02.95" would return ["C", "02", "95"]
// - "D.12" would return ["D", "12", ""]
// - "E" would return ["E", "", ""]
// - "" or invalid input would return ["", "", ""]
//
// This function handles various edge cases:
// - Input without dots
// - Input with fewer than 3 components
// - Empty or nil input
func (p *MessageProcessor) SplitEbkphCode(ebkphCode string) [3]string {
	// Initialize an empty array for the results
	var components [3]string

	// Return empty components if the input is empty
	if len(strings.TrimSpace(ebkphCode)) == 0 {
		return components
	}

	// Split the eBKP-H code by dots
	parts := strings.Split(ebkphCode, ".")

	// Assign the parts to the components array
	// If a part exists, assign it, otherwise leave it as an empty string
	for i := 0; i < len(parts) && i < 3; i++ {
		components[i] = strings.TrimSpace(parts[i])
	}

	return components
}

// ProcessLcaMessage processes an LCA message before it's written to the database
// It handles any transformations or validations needed
func (p *MessageProcessor) ProcessLcaMessage(message *LcaMessage) ([]EavMaterialDataItem, error) {
	log.Printf("Processing LCA message for project: %s, filename: %s", message.Project, message.Filename)

	// Clean the filename by removing .ifc extension if present
	message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	// Validate required fields
	if message.Project == "" {
		return nil, ErrMissingProject
	}

	if message.Filename == "" {
		return nil, ErrMissingFilename
	}

	if message.Timestamp == "" {
		return nil, ErrMissingTimestamp
	}

	if len(message.Data) == 0 {
		return nil, ErrEmptyData
	}

	// Add fileid to the message for easier use by the writer
	message.FileID = p.GenerateFileID(message.Project, message.Filename)

	// Convert the LcaDataItem to EavMaterialDataItem
	eavItems := make([]EavMaterialDataItem, 0)
	for _, item := range message.Data {
		// Create base item with common properties
		baseItem := EavMaterialDataItem{
			Project:   message.Project,
			Filename:  message.Filename,
			FileID:    message.FileID,
			Timestamp: message.Timestamp,
			Id:        item.Id,
			Sequence:  item.Sequence,
		}

		// Convert all properties to EAV items using reflection
		eavItems = append(eavItems, p.structToEavItemsMaterial(baseItem, item)...)
	}

	return eavItems, nil
}

// ProcessCostMessage processes a Cost message before it's written to the database
// It handles any transformations or validations needed
func (p *MessageProcessor) ProcessCostMessage(message *CostMessage) ([]EavElementDataItem, error) {
	log.Printf("Processing Cost message for project: %s, filename: %s", message.Project, message.Filename)

	// Clean the filename by removing .ifc extension if present
	message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	// Validate required fields
	if message.Project == "" {
		return nil, ErrMissingProject
	}

	if message.Filename == "" {
		return nil, ErrMissingFilename
	}

	if message.Timestamp == "" {
		return nil, ErrMissingTimestamp
	}

	if len(message.Data) == 0 {
		return nil, ErrEmptyData
	}

	// Add fileid to the message for easier use by the writer
	message.FileID = p.GenerateFileID(message.Project, message.Filename)

	// Convert the CostDataItem to EavElementDataItem
	eavItems := make([]EavElementDataItem, 0)
	for _, item := range message.Data {
		// Create base item with common properties
		baseItem := EavElementDataItem{
			Project:   message.Project,
			Filename:  message.Filename,
			FileID:    message.FileID,
			Timestamp: message.Timestamp,
			Id:        item.Id,
		}

		// Convert all properties to EAV items using reflection
		eavItems = append(eavItems, p.structToEavItemsElement(baseItem, item)...)
	}

	return eavItems, nil
}

// EavDataItem interface defines the common fields between EavElementDataItem and EavMaterialDataItem
type EavDataItem interface {
	// No methods needed, this is just for type grouping
}

// this is super dump but I just could not figure out how to do this with generics
func (p *MessageProcessor) structToEavItemsElement(baseItem EavElementDataItem, data interface{}) []EavElementDataItem {
	items := make([]EavElementDataItem, 0)
	v := reflect.ValueOf(data)
	t := v.Type()

	// Iterate through all fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip the Id field as it's part of the base item
		if field.Name == "Id" {
			continue
		}

		// Get the json tag name, fallback to field name if not present
		jsonTag := field.Tag.Get("json")
		paramName := strings.Split(jsonTag, ",")[0]
		if paramName == "" {
			paramName = strings.ToLower(field.Name)
		}

		item := baseItem // Create a copy of the base item
		item.ParamName = paramName

		// Reset all value fields
		item.ParamValueString = nil
		item.ParamValueNumber = nil
		item.ParamValueBoolean = nil
		item.ParamValueDate = nil

		// Set the appropriate value based on field type
		switch value.Kind() {
		case reflect.String:
			str := value.String()
			item.ParamType = "string"
			item.ParamValueString = &str
		case reflect.Float32:
			num := float32(value.Float())
			item.ParamType = "number"
			item.ParamValueNumber = &num
		case reflect.Bool:
			b := value.Bool()
			item.ParamType = "boolean"
			item.ParamValueBoolean = &b
		case reflect.Struct:
			// Check if it's a time.Time
			if value.Type() == reflect.TypeOf(time.Time{}) {
				timeStr := value.Interface().(time.Time).Format(time.RFC3339)
				item.ParamType = "date"
				item.ParamValueDate = &timeStr
			}
		}

		items = append(items, item)
	}

	return items
}

// this is super dump but I just could not figure out how to do this with generics
func (p *MessageProcessor) structToEavItemsMaterial(baseItem EavMaterialDataItem, data interface{}) []EavMaterialDataItem {
	items := make([]EavMaterialDataItem, 0)
	v := reflect.ValueOf(data)
	t := v.Type()

	// Iterate through all fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip the Id field as it's part of the base item
		if field.Name == "Id" {
			continue
		}

		// Get the json tag name, fallback to field name if not present
		jsonTag := field.Tag.Get("json")
		paramName := strings.Split(jsonTag, ",")[0]
		if paramName == "" {
			paramName = strings.ToLower(field.Name)
		}

		item := baseItem // Create a copy of the base item
		item.ParamName = paramName

		// Reset all value fields
		item.ParamValueString = nil
		item.ParamValueNumber = nil
		item.ParamValueBoolean = nil
		item.ParamValueDate = nil

		// Set the appropriate value based on field type
		switch value.Kind() {
		case reflect.String:
			str := value.String()
			item.ParamType = "string"
			item.ParamValueString = &str
		case reflect.Float32:
			num := float32(value.Float())
			item.ParamType = "number"
			item.ParamValueNumber = &num
		case reflect.Bool:
			b := value.Bool()
			item.ParamType = "boolean"
			item.ParamValueBoolean = &b
		case reflect.Struct:
			// Check if it's a time.Time
			if value.Type() == reflect.TypeOf(time.Time{}) {
				timeStr := value.Interface().(time.Time).Format(time.RFC3339)
				item.ParamType = "date"
				item.ParamValueDate = &timeStr
			}
		}

		items = append(items, item)
	}

	return items
}
