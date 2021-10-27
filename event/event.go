package event

// CantabularMetadataExport provides an avro structure for CantabularMetadataExport event
type CantabularMetadataExport struct {
	DatasetID    string `avro:"dataset_id"`
	Edition      string `avro:"edition"`
	Version      int32  `avro:"version"`
	CollectionID string
}
