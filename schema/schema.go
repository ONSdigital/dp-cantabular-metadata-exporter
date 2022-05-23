package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var (
	csvCreated = `{
    "type": "record",
    "name": "cantabular-metadata-export",
    "fields": [
      {"name": "instance_id", "type": "string", "default": ""},
      {"name": "dataset_id",  "type": "string", "default": ""},
      {"name": "edition",     "type": "string", "default": ""},
      {"name": "version",     "type": "string", "default": ""},
      {"name": "row_count",   "type": "int",    "default": 0 },
      {"name": "file_name",   "type": "string", "default": ""},
      {"name": "dimensions", "type": { "type": "array", "items": "string"}, "default": [] }
    ]
  }`

	csvwCreated = `{
    "type": "record",
    "name": "cantabular-metadata-complete",
    "fields": [
      {"name": "instance_id", "type": "string", "default": ""},
      {"name": "dataset_id",  "type": "string", "default": ""},
      {"name": "edition",     "type": "string", "default": ""},
      {"name": "version",     "type": "string", "default": ""},
      {"name": "row_count",   "type": "int", "default": 0}
    ]
  }`

	// CSVCreated is the Avro schema for Metadata Export messages.
	CSVCreated = &avro.Schema{
		Definition: csvCreated,
	}

	// CSVWCreated is the Avro schema for Metadata Complete messages.
	CSVWCreated = &avro.Schema{
		Definition: csvwCreated,
	}
)
