package csvw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/log.go/v2/log"
)

// CSVW provides a structure for describing a CSV through a JSON metadata file.
// The JSON tags feature web vocabularies like Dublin Core, DCAT and Stat-DCAT
// to help further contextualize and define the metadata being provided.
// The URL field in the CSVW must reference a CSV file, and all other data
// should describe that CSVs contents.
type CSVW struct {
	Context     string    `json:"@context"`
	URL         string    `json:"url"`
	Title       string    `json:"dct:title"`
	Description string    `json:"dct:description,omitempty"`
	Issued      string    `json:"dct:issued,omitempty"`
	Publisher   Publisher `json:"dct:publisher,omitempty"`
	Contact     []Contact `json:"dcat:contactPoint,omitempty"`
	TableSchema Columns   `json:"tableSchema"`
	Theme       string    `json:"dcat:theme,omitempty"`
	License     string    `json:"dct:license,omitempty"`
	Frequency   string    `json:"dct:accrualPeriodicity,omitempty"`
	Notes       []Note    `json:"notes,omitempty"`
	// New fields
	Downloads     map[string]Download `json:"dcat:distribution,omitempty"`
	Keywords      []string            `json:"dcat:keyword,omitempty"`
	UnitOfMeasure string              `json:"dcat:conformsTo,omitempty"`
	Version       int                 `json:"dcat:version,omitempty"`
	IsBasedOn     string              `json:"prov:wasGeneratedBy,omitempty"`
	DatasetLinks  *DatasetLinks       `json:"dcat:record,omitempty"`
}

// New fields
type DatasetLinks struct {
	Editions      Link `json:"editions"`
	LatestVersion Link `json:"latestVersion"`
	Self          Link `json:"self"`
}
type Link struct {
	HREF string `json:"HREF"`
	ID   string `json:"ID"`
}

type Download struct {
	HREF string `json:"HREF"`
	Size string `json:"Size"`
}

// Old fields
// Contact represents a response model within a dataset
type Contact struct {
	Name      string `json:"vcard:fn"`
	Telephone string `json:"vcard:tel"`
	Email     string `json:"vcard:email"`
}

// Publisher defines the entity primarily responsible for the dataset
// https://www.w3.org/TR/vocab-dcat/#class-catalog
type Publisher struct {
	Name string `json:"name,omitempty"`
	Type string `json:"@type,omitempty"`
	ID   string `json:"@id"` // a URL where more info is available
}

// Columns provides the nested structure expected within the tableSchema of a CSVW
type Columns struct {
	C     []Column `json:"columns"`
	About string   `json:"aboutUrl,omitempty"`
}

// Column provides the ability to define the JSON tags required specific
// to each column within the CSVW
type Column map[string]interface{}

// Note can include alerts, corrections or usage notes which exist in the
// dataset metadata and help describe the contents of the CSV
type Note struct {
	Type       string `json:"type"` // is this an enum?
	Target     string `json:"target,omitempty"`
	Body       string `json:"body"`
	Motivation string `json:"motivation,omitempty"` // how is this different from type? do we need this? is this an enum?
}

var errMissingDimensions = errors.New("no dimensions in provided metadata")

// New CSVW returned with top level fields populated based on provided metadata
func New(m *dataset.Metadata, csvURL, externalPrefixURL, filterOutputID, downloadServiceURL string, isCustom bool) *CSVW {
	if isCustom {
		dt := time.Now()
		issuedDate := dt.Format("01-02-2006 15:04:05")
		titleDims := GenerateCustomTitle(m.Version.Dimensions)

		csvw := &CSVW{
			Context: "http://www.w3.org/ns/csvw",
			Title:   titleDims,
			Issued:  issuedDate,
			URL:     fmt.Sprintf("%s/downloads/filter-outputs/%s.csvw", downloadServiceURL, filterOutputID),
		}
		csvw.UnitOfMeasure = m.UnitOfMeasure
		return csvw
	} else {
		csvw := &CSVW{
			Context:     "http://www.w3.org/ns/csvw",
			Title:       m.Title,
			Description: m.Description,
			Issued:      m.ReleaseDate,
			Theme:       m.Theme,
			License:     m.License,
			Frequency:   m.ReleaseFrequency,
			URL:         csvURL,
		}

		if m.Contacts != nil {
			for _, c := range *m.Contacts {
				csvw.Contact = append(csvw.Contact, Contact{
					Name:      c.Name,
					Telephone: c.Telephone,
					Email:     c.Email,
				})
			}
		}

		if m.Publisher != nil {
			csvw.Publisher = Publisher{
				Name: m.Publisher.Name,
				Type: m.Publisher.Type,
				ID:   m.Publisher.URL,
			}
		}

		if m.Downloads != nil {
			csvw.Downloads = make(map[string]Download)
			for k, v := range m.Downloads {
				csvw.Downloads[k] = Download{
					HREF: v.URL,
					Size: v.Size,
				}
			}
		}

		if m.Keywords != nil {
			csvw.Keywords = *m.Keywords
		}

		csvw.UnitOfMeasure = m.UnitOfMeasure
		csvw.Version = m.Version.Version
		if m.Version.IsBasedOn != nil {
			csvw.IsBasedOn = m.Version.IsBasedOn.ID
		}

		re := regexp.MustCompile("https?://([^/]+)")
		csvw.DatasetLinks = &DatasetLinks{
			Editions: Link{
				HREF: re.ReplaceAllString(m.DatasetLinks.Editions.URL, externalPrefixURL),
				ID:   m.DatasetLinks.Editions.ID,
			},
			LatestVersion: Link{
				HREF: re.ReplaceAllString(m.DatasetLinks.LatestVersion.URL, externalPrefixURL),
				ID:   m.DatasetLinks.LatestVersion.ID,
			},
			Self: Link{
				HREF: re.ReplaceAllString(m.DatasetLinks.Self.URL, externalPrefixURL),
				ID:   m.DatasetLinks.Editions.ID,
			},
		}
		return csvw
	}
}

// Generate the CSVW structured metadata file to describe a CSV
func Generate(ctx context.Context, metadata *dataset.Metadata, downloadURL, aboutURL, apiDomain, externalPrefixURL, filterOutputID, downloadServiceURL string, isCustom bool) ([]byte, error) {
	logData := log.Data{
		"dataset_id": metadata.DatasetDetails.ID,
		"csv_header": metadata.CSVHeader,
	}

	log.Info(ctx, "generating csvw file", logData)

	// if len(metadata.Dimensions) == 0 {
	// 	return nil, Error{
	// 		err:     errMissingDimensions,
	// 		logData: logData,
	// 	}
	// }

	if len(metadata.CSVHeader) < 1 {
		return nil, Error{
			err:     errors.New("CSV header empty or missing dimensions"),
			logData: logData,
		}
	}

	h := metadata.CSVHeader

	csvw := New(metadata, downloadURL, externalPrefixURL, filterOutputID, downloadServiceURL, isCustom)

	var list []Column
	obs := newObservationColumn(h[0], metadata.UnitOfMeasure)
	list = append(list, obs)

	// add dimension columns
	if len(metadata.Dimensions) > 0 {
		for i, d := range metadata.Dimensions {
			l, err := newLabelColumn(i, apiDomain, h, d)
			if err != nil {
				return nil, Error{
					err:     fmt.Errorf("failed to create label column: %w", err),
					logData: logData,
				}
			}
			list = append(list, l)
		}
	}

	aboutURL, err := formatAboutURL(aboutURL, apiDomain)
	if err != nil {
		return nil, Error{
			err:     fmt.Errorf("failed to format AboutURL: %w", err),
			logData: logData,
		}
	}
	if !isCustom {
		csvw.TableSchema = Columns{
			About: aboutURL,
			C:     list,
		}
	} else {
		csvw.TableSchema = Columns{
			About: "",
			C:     list,
		}
	}

	csvw.AddNotes(metadata, downloadURL)

	b, err := json.Marshal(csvw)
	if err != nil {
		return nil, Error{
			err:     fmt.Errorf("failed to marshal csvw: %w", err),
			logData: logData,
		}
	}

	return b, nil
}

func formatAboutURL(aboutURL, domain string) (string, error) {
	about, err := url.Parse(aboutURL)
	if err != nil {
		return "", Error{
			err: fmt.Errorf("failed to parse aboutURL: %w", err),
			logData: log.Data{
				"about_url": aboutURL,
			},
		}
	}

	d, err := url.Parse(domain)
	if err != nil {
		return "", Error{
			err: fmt.Errorf("failed to parse domain: %w", err),
			logData: log.Data{
				"domain": domain,
			},
		}
	}

	d.Path += about.Path

	return d.String(), nil
}

// AddNotes to CSVW from alerts or usage notes in provided metadata
func (csvw *CSVW) AddNotes(metadata *dataset.Metadata, url string) {
	if metadata.Alerts != nil {
		for _, a := range *metadata.Alerts {
			csvw.Notes = append(csvw.Notes, Note{
				Type:   a.Type,
				Body:   a.Description,
				Target: url,
			})
		}
	}

	if metadata.DatasetDetails.UsageNotes != nil {
		for _, u := range *metadata.DatasetDetails.UsageNotes {
			csvw.Notes = append(csvw.Notes, Note{
				Type: u.Title,
				Body: u.Note,
				// TODO: store column number against usage notes in Dataset API
				//Target: url + "#col=need-to-store",
			})
		}
	}
}

func newObservationColumn(title, name string) Column {
	c := newColumn(title, name)

	c["datatype"] = "string"

	return c
}

// TBC what the final content of these columns should be. For now is a rough port and amalgamation
// of the code and label columns from the CMD implementation. Will be updated when we have full
// metadata spec.
func newLabelColumn(i int, apiDomain string, header []string, dim dataset.VersionDimension) (Column, error) {
	dimURL := dim.URL
	if len(apiDomain) > 0 {
		uri, err := url.Parse(dimURL)
		if err != nil {
			return nil, Error{
				err: fmt.Errorf("failed to parse dimension url: %w", err),
				logData: log.Data{
					"dimension_url": dimURL,
				},
			}
		}

		dimURL = fmt.Sprintf("%s%s", apiDomain, uri.Path)
	}

	labelCol := newColumn(dim.Name, dim.Label)
	labelCol["description"] = dim.Description
	labelCol["valueURL"] = dimURL + "/codes/{" + dim.Label + "}"
	labelCol["required"] = true
	labelCol["optionCount"] = dim.NumberOfOptions
	// TODO: determine what could go in c["datatype"] and c["required"]

	return labelCol, nil
}

func newColumn(title, label string) Column {
	c := make(Column)
	if len(title) > 0 {
		c["titles"] = title
	}

	if len(label) > 0 {
		c["label"] = label
	}

	return c
}

func GenerateCustomTitle(dims []dataset.VersionDimension) string {
	var title string
	l := len(dims)
	for i, d := range dims {
		if i == 0 {
			title += d.Label
		} else if i == (l - 1) {
			title += " and " + d.Label
		} else {
			title += ", " + d.Label
		}
	}
	return title
}
