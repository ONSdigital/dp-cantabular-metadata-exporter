package schema

import (
	"github.com/ONSdigital/dp-kafka/v2/avro"
)

var cantabularMetadataExport = `{
  "type": "record",
  "name": "cantabular-metadata-export",
  "fields": [
    {"name": "id",     "type": "string", "default": ""}
  ]
}`

// CantabularMetadataExport is the Avro schema for Instance Complete messages.
var CantabularMetadataExport = &avro.Schema{
	Definition: cantabularMetadataExport,
}