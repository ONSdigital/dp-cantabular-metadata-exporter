package text

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/custom"
)

// NewMetadata returns a .txt metadata based on provided metadata
func NewMetadata(m *dataset.Metadata, filterOutputID, downloadServiceURL string) []byte {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("Title: %s\n", m.Title))
	b.WriteString(fmt.Sprintf("Description: %s\n", m.Description))
	if m.Publisher != nil {
		b.WriteString(fmt.Sprintf("Publisher: %s\n", *m.Publisher))
	}
	b.WriteString(fmt.Sprintf("Issued: %s\n", m.ReleaseDate))
	b.WriteString(fmt.Sprintf("Next Release: %s\n", m.NextRelease))
	b.WriteString(fmt.Sprintf("Identifier: %s\n", m.Title))
	if m.Keywords != nil {
		b.WriteString(fmt.Sprintf("Keywords: %s\n", *m.Keywords))
	}
	b.WriteString(fmt.Sprintf("Language: %s\n", "English"))
	if m.Contacts != nil {
		contacts := *m.Contacts
		if len(contacts) > 0 {
			b.WriteString(fmt.Sprintf("Contact: %s, %s, %s\n", contacts[0].Name, contacts[0].Email, contacts[0].Telephone))
		}
	}
	if len(m.Temporal) > 0 {
		b.WriteString(fmt.Sprintf("Temporal: %s\n", m.Temporal[0].Frequency))
	}
	b.WriteString(fmt.Sprintf("Latest Changes: %s\n", m.LatestChanges))
	b.WriteString(fmt.Sprintf("Periodicity: %s\n", m.ReleaseFrequency))
	b.WriteString("Distribution:\n")

	downloadKeys := getSortedKeys(m.Downloads)
	for _, k := range downloadKeys {
		v := m.Downloads[k]
		b.WriteString(fmt.Sprintf("\tExtension: %s\n", k))
		b.WriteString(fmt.Sprintf("\tSize: %s\n", v.Size))
		b.WriteString(fmt.Sprintf("\tURL: %s\n\n", v.URL))
	}
	b.WriteString(fmt.Sprintf("Unit of measure: %s\n", m.UnitOfMeasure))
	b.WriteString(fmt.Sprintf("License: %s\n", m.License))
	if m.Methodologies != nil {
		b.WriteString(fmt.Sprintf("Methodologies: %s\n", *m.Methodologies))
	}
	if m.Publications != nil {
		b.WriteString(fmt.Sprintf("Publications: %s\n", *m.Publications))
	}
	if m.RelatedDatasets != nil {
		b.WriteString(fmt.Sprintf("Related Links: %s\n", *m.RelatedDatasets))
	}
	b.WriteString(fmt.Sprintf("Canonical Topic: %s\n", m.CanonicalTopic))
	if len(m.Subtopics) > 0 {
		b.WriteString(fmt.Sprintf("Subtopics: %s\n", m.Subtopics))
	}
	b.WriteString(fmt.Sprintf("Survey: %s\n", m.Survey))

	// New fields
	// Combine alerts and useage notes as in csvw?
	if m.Alerts != nil && len(*m.Alerts) != 0 {
		b.WriteString(fmt.Sprintf("Alerts:\n"))

		for _, a := range *m.Alerts {
			b.WriteString(
				fmt.Sprintf(
					"Date: %s\nDescription: %s\nType: %s\n",
					a.Date,
					a.Description,
					a.Type,
				),
			)
		}
	}

	b.WriteString("Usage Notes:\n")
	if m.DatasetDetails.UsageNotes != nil {
		for _, n := range *m.DatasetDetails.UsageNotes {
			b.WriteString(fmt.Sprintf("Title: %s\nNote: %s\n\n", n.Title, n.Note))
		}
	}

	if l := m.DatasetLinks; !reflect.DeepEqual(l, dataset.Link{}) {
		edEmpty := !reflect.DeepEqual(l.Editions, dataset.Link{})
		lvEmpty := !reflect.DeepEqual(l.LatestVersion, dataset.Link{})
		sEmpty := !reflect.DeepEqual(l.Self, dataset.Link{})
		if !edEmpty || !lvEmpty {
			b.WriteString("Dataset Links:\n")
		}
		if !edEmpty {
			b.WriteString(fmt.Sprintf("Editions:\nHREF: %s\n", l.Editions.URL))
		}
		if !lvEmpty {
			b.WriteString(
				fmt.Sprintf(
					"Latest Version:\nHREF: %s\nID: %s\n",
					l.LatestVersion.URL,
					l.LatestVersion.ID,
				),
			)
		}
		if !sEmpty {
			b.WriteString(fmt.Sprintf("Self:\nHREF: %s\n", l.Self.URL))
		}
	}

	b.WriteString(fmt.Sprintf("Version: %d\n", m.Version.Version))

	if m.Version.IsBasedOn != nil {
		b.WriteString(fmt.Sprintf("Is Based On: %s\n", m.Version.IsBasedOn.ID))
	}

	b.WriteString(fmt.Sprintf("%s\n\n%s\n\n%s\n\n", variables, coverage, scc))

	b.WriteString("Dimensions:\n")
	for _, d := range m.Version.Dimensions {
		b.WriteString(fmt.Sprintf("\n\tID: %s\n", d.ID))
		b.WriteString(fmt.Sprintf("\n\tLabel: %s\n", d.Label))
		b.WriteString(fmt.Sprintf("\n\tDescription: %s\n", d.Description))
		b.WriteString(fmt.Sprintf("\n\tNumber Of Options: %d\n", d.NumberOfOptions))
		b.WriteString(fmt.Sprintf("\n\tQuality Statement: %s\n%s\n", d.QualityStatementText, d.QualityStatementURL))
	}
	return b.Bytes()
}

// NewMetadataCustom returns .txt metadata based on provided metadata for custom datasets
func NewMetadataCustom(m *dataset.Metadata, filterOutputID, downloadServiceURL string) []byte {
	var b bytes.Buffer
	dt := time.Now()
	issuedDate := dt.Format("01-02-2006 15:04:05")

	titleDims := custom.GenerateCustomTitle(m.Version.Dimensions)

	b.WriteString(fmt.Sprintf("Title: %s\n", titleDims))
	b.WriteString(fmt.Sprintf("Issued: %s\n", issuedDate))
	b.WriteString(fmt.Sprintf("Language: %s\n", "English"))
	b.WriteString("Distribution:\n")
	downloadKeys := getSortedKeys(m.Downloads)
	for _, k := range downloadKeys {
		v := m.Downloads[k]
		if k == "csv" {
			b.WriteString(fmt.Sprintf("\tExtension: %s\n", k))
			b.WriteString(fmt.Sprintf("\tSize: %s\n", v.Size))
			b.WriteString(fmt.Sprintf("\tURL: %s\n\n", fmt.Sprintf("%s/downloads/filter-outputs/%s.csv", downloadServiceURL, filterOutputID)))
		} else if k == "csvw" {
			b.WriteString(fmt.Sprintf("\tExtension: %s\n", k))
			b.WriteString(fmt.Sprintf("\tSize: %s\n", v.Size))
			b.WriteString(fmt.Sprintf("\tURL: %s\n\n", fmt.Sprintf("%s/downloads/filter-outputs/%s.csv-metadata.json", downloadServiceURL, filterOutputID)))
		} else if k == "txt" {
			b.WriteString(fmt.Sprintf("\tExtension: %s\n", k))
			b.WriteString(fmt.Sprintf("\tSize: %s\n", v.Size))
			b.WriteString(fmt.Sprintf("\tURL: %s\n\n", fmt.Sprintf("%s/downloads/filter-outputs/%s.txt", downloadServiceURL, filterOutputID)))
		} else if k == "xlsx" {
			b.WriteString(fmt.Sprintf("\tExtension: %s\n", k))
			b.WriteString(fmt.Sprintf("\tSize: %s\n", v.Size))
			b.WriteString(fmt.Sprintf("\tURL: %s\n\n", fmt.Sprintf("%s/downloads/filter-outputs/%s.xlsx", downloadServiceURL, filterOutputID)))
		}

	}
	if m.Version.IsBasedOn != nil {
		b.WriteString(fmt.Sprintf("Is Based On: %s\n", m.Version.IsBasedOn.ID))
	}
	b.WriteString(fmt.Sprintf("%s\n\n%s\n\n%s\n\n", variables, coverage, scc))
	b.WriteString("Dimensions:\n")
	for _, d := range m.Version.Dimensions {
		b.WriteString(fmt.Sprintf("\n\tID: %s\n", d.ID))
		b.WriteString(fmt.Sprintf("\n\tLabel: %s\n", d.Label))
		b.WriteString(fmt.Sprintf("\n\tDescription: %s\n", d.Description))
		b.WriteString(fmt.Sprintf("\n\tNumber Of Options: %d\n", d.NumberOfOptions))
		b.WriteString(fmt.Sprintf("\n\tQuality Statement: %s\n%s\n", d.QualityStatementText, d.QualityStatementURL))
	}

	return b.Bytes()
}

func getSortedKeys[T interface{}](m map[string]T) []string {
	keySlice := make([]string, 0)
	for key := range m {
		keySlice = append(keySlice, key)
	}
	sort.Strings(keySlice)
	return keySlice
}
