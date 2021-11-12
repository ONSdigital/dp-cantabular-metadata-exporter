package event

// CantabularMetadataExport provides an avro structure for CantabularMetadataExport event
type CantabularMetadataExport struct {
	DatasetID    string `avro:"dataset_id"`
	Edition      string `avro:"edition"`
	Version      string `avro:"version"`
	CollectionID string `avro:"collection_id"`
}
