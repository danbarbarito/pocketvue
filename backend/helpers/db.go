package helpers

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// FindCollectionByName finds a collection by name and returns a descriptive error
func FindCollectionByName(app core.App, collectionName string) (*core.Collection, error) {
	collection, err := app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to find %s collection: %w", collectionName, err)
	}
	return collection, nil
}

// FindRecordByID finds a record by ID and returns a descriptive error
func FindRecordByID(app core.App, collectionName string, recordID string) (*core.Record, error) {
	collection, err := FindCollectionByName(app, collectionName)
	if err != nil {
		return nil, err
	}

	record, err := app.FindRecordById(collection, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to find record %s in %s: %w", recordID, collectionName, err)
	}
	return record, nil
}

// SaveRecord saves a record and returns a descriptive error
func SaveRecord(app core.App, record *core.Record) error {
	if err := app.Save(record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}
	return nil
}

// CreateRecord creates a new record for the given collection
func CreateRecord(app core.App, collectionName string) (*core.Record, error) {
	collection, err := FindCollectionByName(app, collectionName)
	if err != nil {
		return nil, err
	}
	return core.NewRecord(collection), nil
}

// FindAllRecords finds all records in a collection
func FindAllRecords(app core.App, collectionName string) ([]*core.Record, error) {
	collection, err := FindCollectionByName(app, collectionName)
	if err != nil {
		return nil, err
	}

	records, err := app.FindAllRecords(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records from %s: %w", collectionName, err)
	}
	return records, nil
}
