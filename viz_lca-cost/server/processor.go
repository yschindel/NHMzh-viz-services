package server

import (
	"reflect"
	"strings"
	"time"
	"viz_lca-cost/logger"
)

// MessageProcessor handles the processing of incoming messages
type MessageProcessor struct {
	logger *logger.Logger
}

// NewMessageProcessor creates a new message processor
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		logger: logger.New().WithFields(logger.Fields{
			"operation": "message_processor",
		}),
	}
}

// ProcessLcaMessage processes an LCA message before it's written to the database
// It handles any transformations or validations needed
func (p *MessageProcessor) ProcessLcaMessage(message *LcaMessage) ([]EavMaterialRow, error) {
	log := p.logger.WithFields(logger.Fields{
		"operation": "process_lca_message",
		"project":   message.Project,
		"filename":  message.Filename,
	})
	log.Info("Processing LCA message")

	// Clean the filename by removing .ifc extension if present
	// message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	// Validate required fields
	if message.Project == "" {
		log.Error("Missing project field")
		return nil, ErrMissingProject
	}

	if message.Filename == "" {
		log.Error("Missing filename field")
		return nil, ErrMissingFilename
	}

	if message.Timestamp == "" {
		log.Error("Missing timestamp field")
		return nil, ErrMissingTimestamp
	}

	if len(message.Data) == 0 {
		log.Error("Empty data field")
		return nil, ErrEmptyData
	}

	// Convert the LcaDataItem to EavMaterialDataItem
	eavItems := make([]EavMaterialRow, 0)
	for _, item := range message.Data {
		// Create base item with common properties
		baseItem := EavMaterialRow{
			Project:   message.Project,
			Filename:  message.Filename,
			Timestamp: message.Timestamp,
			Id:        item.Id,
			Sequence:  item.Sequence,
		}

		// Convert all properties to EAV items using reflection
		eavItems = append(eavItems, p.structToEavItemsMaterial(baseItem, item)...)
	}

	log.Info("Successfully processed LCA items", logger.Fields{
		"count": len(eavItems),
	})
	return eavItems, nil
}

// ProcessCostMessage processes a Cost message before it's written to the database
// It handles any transformations or validations needed
func (p *MessageProcessor) ProcessCostMessage(message *CostMessage) ([]EavElementRow, error) {
	log := p.logger.WithFields(logger.Fields{
		"operation": "process_cost_message",
		"project":   message.Project,
		"filename":  message.Filename,
	})
	log.Info("Processing Cost message")

	// Clean the filename by removing .ifc extension if present
	// message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	// Validate required fields
	if message.Project == "" {
		log.Error("Missing project field")
		return nil, ErrMissingProject
	}

	if message.Filename == "" {
		log.Error("Missing filename field")
		return nil, ErrMissingFilename
	}

	if message.Timestamp == "" {
		log.Error("Missing timestamp field")
		return nil, ErrMissingTimestamp
	}

	if len(message.Data) == 0 {
		log.Error("Empty data field")
		return nil, ErrEmptyData
	}

	// Convert the CostDataItem to EavElementDataItem
	eavItems := make([]EavElementRow, 0)
	for _, item := range message.Data {
		// Create base item with common properties
		baseItem := EavElementRow{
			Project:   message.Project,
			Filename:  message.Filename,
			Timestamp: message.Timestamp,
			Id:        item.Id,
		}

		// Convert all properties to EAV items using reflection
		eavItems = append(eavItems, p.structToEavItemsElement(baseItem, item)...)
	}

	log.Info("Successfully processed Cost items", logger.Fields{
		"count": len(eavItems),
	})
	return eavItems, nil
}

// EavDataItem interface defines the common fields between EavElementDataItem and EavMaterialDataItem
type EavDataItem interface {
	// No methods needed, this is just for type grouping
}

// this is super dump but I just could not figure out how to do this with generics
func (p *MessageProcessor) structToEavItemsElement(baseItem EavElementRow, data interface{}) []EavElementRow {
	items := make([]EavElementRow, 0)
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
func (p *MessageProcessor) structToEavItemsMaterial(baseItem EavMaterialRow, data interface{}) []EavMaterialRow {
	items := make([]EavMaterialRow, 0)
	v := reflect.ValueOf(data)
	t := v.Type()

	// Iterate through all fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip the Id field as it's part of the base item
		if field.Name == "Id" || field.Name == "Sequence" {
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
