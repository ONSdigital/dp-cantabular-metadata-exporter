package schema

import (
	"github.com/ONSdigital/dp-kafka/v2/avro"
)

var (
  cantabularMetadataExport = `{
    "type": "record",
    "name": "cantabular-metadata-export",
    "fields": [
      {"name": "dataset_id",       "type": "string",    "default":  ""},
      {"name": "edition",          "type": "string",    "default":  ""},
      {"name": "version",          "type": "string",    "default":  ""},
      {"name": "collection_id",    "type": "string",    "default":  ""},
      {"name": "row_count",        "type": "int",       "default":   0}
    ]
  }`

  cantabularMetadataComplete = `{
    "type": "record",
    "name": "cantabular-metadata-complete",
    "fields": [
      {"name": "dataset_id",       "type": "string",    "default":  ""},
      {"name": "edition",          "type": "string",    "default":  ""},
      {"name": "version",          "type": "string",    "default":  ""},
      {"name": "collection_id",    "type": "string",    "default":  ""},
      {"name": "row_count",        "type": "int",       "default":   0}
    ]
  }`

  // CantabularMetadataExport is the Avro schema for Metadata Export messages.
  CantabularMetadataExport = &avro.Schema{
    Definition: cantabularMetadataExport,
  }

  // CantabularMetadataExport is the Avro schema for Metadata Complete messages.
  CantabularMetadataComplete = &avro.Schema{
    Definition: cantabularMetadataExport,
  }
)
