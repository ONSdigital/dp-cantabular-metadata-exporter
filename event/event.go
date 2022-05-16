package event

// CSVCreated provides an avro structure for CSVCreated event
type CSVCreated struct {
	DatasetID    string   `avro:"dataset_id"`
	Edition      string   `avro:"edition"`
	Version      string   `avro:"version"`
	InstanceID   string   `avro:"instance_id"`
	RowCount     int32    `avro:"row_count"`
	DimensionIDs []string `avro:"dimension_ids"`
}

// CSVWCreated provides an avro structure for CSVWCreated event
type CSVWCreated struct {
	DatasetID  string `avro:"dataset_id"`
	Edition    string `avro:"edition"`
	Version    string `avro:"version"`
	InstanceID string `avro:"instance_id"`
	RowCount   int32  `avro:"row_count"`
}
