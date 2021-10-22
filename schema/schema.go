package schema

import (
	"github.com/ONSdigital/dp-kafka/v2/avro"
)

var cantabularMetadataExport = `{
  "type": "record",
  "name": "cantabular-metadata-export",
  "fields": [
    {"name": "dataset_id",       "type": "string",    "default":  ""},
    {"name": "edition",          "type": "string",    "default":  ""},
    {"name": "version",          "type": "int",       "default":   0},
    {"name": "collection_id",    "type": "string",    "default":  ""}
}`

// CantabularMetadataExport is the Avro schema for Instance Complete messages.
var CantabularMetadataExport = &avro.Schema{
	Definition: cantabularMetadataExport,
}
