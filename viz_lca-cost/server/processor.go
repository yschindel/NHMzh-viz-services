package server

import (
	"log"
	"strings"
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
func (p *MessageProcessor) ProcessLcaMessage(message *LcaMessage) error {
	log.Printf("Processing LCA message for project: %s, filename: %s", message.Project, message.Filename)

	// Clean the filename by removing .ifc extension if present
	message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	// Validate required fields
	if message.Project == "" {
		return ErrMissingProject
	}

	if message.Filename == "" {
		return ErrMissingFilename
	}

	if message.Timestamp == "" {
		return ErrMissingTimestamp
	}

	if len(message.Data) == 0 {
		return ErrEmptyData
	}

	// Add fileid to the message for easier use by the writer
	message.FileID = p.GenerateFileID(message.Project, message.Filename)

	// Additional processing logic can be added here

	return nil
}

// ProcessCostMessage processes a Cost message before it's written to the database
// It handles any transformations or validations needed
func (p *MessageProcessor) ProcessCostMessage(message *CostMessage) error {
	log.Printf("Processing Cost message for project: %s, filename: %s", message.Project, message.Filename)

	// Clean the filename by removing .ifc extension if present
	message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	// Validate required fields
	if message.Project == "" {
		return ErrMissingProject
	}

	if message.Filename == "" {
		return ErrMissingFilename
	}

	if message.Timestamp == "" {
		return ErrMissingTimestamp
	}

	if len(message.Data) == 0 {
		return ErrEmptyData
	}

	// Add fileid to the message for easier use by the writer
	message.FileID = p.GenerateFileID(message.Project, message.Filename)

	// Additional processing logic can be added here

	return nil
}
