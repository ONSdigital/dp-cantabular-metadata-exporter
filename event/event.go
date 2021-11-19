package event

// CantabularMetadataExport provides an avro structure for CantabularMetadataExport event
type CSVCreated struct {
	DatasetID    string `avro:"dataset_id"`
	Edition      string `avro:"edition"`
	Version      string `avro:"version"`
	InstanceID   string `avro:"instance_id"`
	RowCount     int32  `avro:"row_count"`
}

// CantabularMetadataExport provides an avro structure for CantabularMetadataExport event
type CSVWCreated struct {
	DatasetID    string `avro:"dataset_id"`
	Edition      string `avro:"edition"`
	Version      string `avro:"version"`
	InstanceID   string `avro:"instance_id"`
	RowCount     int32  `avro:"row_count"`
}
