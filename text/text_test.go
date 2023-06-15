package text

import (
	"strings"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	. "github.com/smartystreets/goconvey/convey"
)

var filterOutputID = "filter-output-id"
var downloadServiceURL = "download-service-url"

func TestNewMetadata(t *testing.T) {
	expectedText := "Title: title\nDescription: description\nPublisher: {url name type}\nIssued: 01-02-2006 15:04:05\nNext Release: next release\nIdentifier: title\nKeywords: [keyword1 keyword2]\nLanguage: English\nContact: name, email, telephone\nTemporal: frequency\nLatest Changes: [{description name type}]\nPeriodicity: release frequency\nDistribution:\n\tExtension: csv\n\tSize: csv size\n\tURL: csv url\n\n\tExtension: csvw\n\tSize: csvw size\n\tURL: csvw url\n\n\tExtension: txt\n\tSize: txt size\n\tURL: txt url\n\n\tExtension: xlsx\n\tSize: xlsx size\n\tURL: xlsx url\n\nUnit of measure: unit of measure\nLicense: license\nMethodologies: [{methodology description methodology url methodology title}]\nPublications: [{publications description publications url publications title}]\nRelated Links: [{url title}]\nCanonical Topic: canonical\nSubtopics: [subtopic]\nSurvey: survey\nAlerts:\nDate: date\nDescription: description\nType: type\nUsage Notes:\nTitle: title\nNote: note\n\nDataset Links:\nEditions:\nHREF: \nLatest Version:\nHREF: \nID: \nSelf:\nHREF: \nVersion: 1\n\nArea Type\n\nCensus 2021 statistics are published for a number of different geographies.\nThese can be large, for example the whole of England, or small, for example\nan output area (OA), the lowest level of geography for which statistics are\nproduced.\n\nFor higher levels of geography, more detailed statistics can be produced.\nWhen a lower level of geography is used, such as output areas (which have\na minimum of 100 persons), the statistics produced have less detail. This is\nto protect the confidentiality of people and ensure that individuals or their\ncharacteristics cannot be identified.\n\n\nCoverage\n\nCensus 2021 statistics are published for the whole of England and Wales.\nHowever, you can choose to filter areas by:\n\n- country - for example, Wales\n- region - for example, London\n- local authority - for example, Cornwall\n- health area – for example, Clinical Commissioning Group\n- statistical area - for example, MSOA or LSOA\n\n\nProtecting personal data\n\nSometimes we need to make changes to data if it is possible to identify\nindividuals. This is known as statistical disclosure control.\n\nIn Census 2021, we:\n- swapped records (targeted record swapping), for example, if a household \n  was likely to be identified in datasets because it has unusual\n  characteristics, we swapped the record with a similar one from a nearby\n  small area (very unusual households could be swapped with one in a nearby\n  local authority)\n- added small changes to some counts (cell key perturbation), for example,\n  we might change a count of four to a three or a five – this might make\n  small differences between tables depending on how the data are broken down\n  when we applied perturbation\n\nRead more in Section 5 of our article Design for Census 2021.\n\nDimensions:\n\n\tID: id 1\n\n\tLabel: Label 1\n\n\tDescription: description 1\n\n\tNumber Of Options: 1\n\n\tQuality Statement: quality statement 1\n\n\n\tID: id 2\n\n\tLabel: Label 2\n\n\tDescription: description 2\n\n\tNumber Of Options: 1\n\n\tQuality Statement: quality statement 2\n\n"
	Convey("Given a metadata struct", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "01-02-2006 15:04:05",
				Dimensions: []dataset.VersionDimension{
					{
						ID:                   "id 1",
						Label:                "Label 1",
						Description:          "description 1",
						NumberOfOptions:      1,
						QualityStatementText: "quality statement 1",
					},
					{
						ID:                   "id 2",
						Label:                "Label 2",
						Description:          "description 2",
						NumberOfOptions:      1,
						QualityStatementText: "quality statement 2",
					},
				},
				Temporal: []dataset.Temporal{
					{
						StartDate: "start date",
						EndDate:   "end date",
						Frequency: "frequency",
					},
				},
				LatestChanges: []dataset.Change{
					{
						Description: "description",
						Name:        "name",
						Type:        "type",
					},
				},
				Downloads: map[string]dataset.Download{
					"xlsx": {
						URL:  "xlsx url",
						Size: "xlsx size",
					},
					"csv": {
						URL:  "csv url",
						Size: "csv size",
					},
					"csvw": {
						URL:  "csvw url",
						Size: "csvw size",
					},
					"txt": {
						URL:  "txt url",
						Size: "txt size",
					},
				},
				Alerts: &[]dataset.Alert{
					{
						Date:        "date",
						Description: "description",
						Type:        "type",
					},
				},
				Links: dataset.Links{
					Editions: dataset.Link{
						URL: "edition url",
						ID:  "edition id",
					},
					LatestVersion: dataset.Link{
						URL: "version url",
						ID:  "version id",
					},
					Self: dataset.Link{
						URL: "self url",
						ID:  "self id",
					},
				},
				Version: 1,
			},
			DatasetDetails: dataset.DatasetDetails{
				Description: "description",
				Title:       "title",
				Publisher: &dataset.Publisher{
					URL:  "url",
					Name: "name",
					Type: "type",
				},
				NextRelease: "next release",
				Keywords: &[]string{
					"keyword1",
					"keyword2",
				},
				Contacts: &[]dataset.Contact{
					{
						Name:      "name",
						Telephone: "telephone",
						Email:     "email",
					},
				},
				ReleaseFrequency: "release frequency",
				UnitOfMeasure:    "unit of measure",
				License:          "license",
				Methodologies: &[]dataset.Methodology{
					{
						Description: "methodology description",
						URL:         "methodology url",
						Title:       "methodology title",
					},
				},
				Publications: &[]dataset.Publication{
					{
						Description: "publications description",
						URL:         "publications url",
						Title:       "publications title",
					},
				},
				RelatedDatasets: &[]dataset.RelatedDataset{
					{
						URL:   "url",
						Title: "title",
					},
				},
				CanonicalTopic: "canonical",
				Subtopics: []string{
					"subtopic",
				},
				Survey: "survey",
				UsageNotes: &[]dataset.UsageNote{
					{
						Note:  "note",
						Title: "title",
					},
				},
			},
		}

		Convey("When the NewMetadata function is called", func() {
			meta := string(NewMetadata(m))
			Convey("Then the returned value should be as expected", func() {
				So(meta, ShouldResemble, expectedText)
			})
		})
	})
}

func TestNewMetadataCustom(t *testing.T) {
	expectedText := "Title: Label 1 and Label 2\nIssued: 1 Jan 2000\nLanguage: English\nDistribution:\n\tExtension: csv\n\tSize: csv size\n\tURL: download-service-url/downloads/filter-outputs/filter-output-id.csv\n\n\tExtension: csvw\n\tSize: csvw size\n\tURL: download-service-url/downloads/filter-outputs/filter-output-id.csv-metadata.json\n\n\tExtension: txt\n\tSize: txt size\n\tURL: download-service-url/downloads/filter-outputs/filter-output-id.txt\n\n\tExtension: xlsx\n\tSize: xlsx size\n\tURL: download-service-url/downloads/filter-outputs/filter-output-id.xlsx\n\nIs Based On: id\n\nArea Type\n\nCensus 2021 statistics are published for a number of different geographies.\nThese can be large, for example the whole of England, or small, for example\nan output area (OA), the lowest level of geography for which statistics are\nproduced.\n\nFor higher levels of geography, more detailed statistics can be produced.\nWhen a lower level of geography is used, such as output areas (which have\na minimum of 100 persons), the statistics produced have less detail. This is\nto protect the confidentiality of people and ensure that individuals or their\ncharacteristics cannot be identified.\n\n\nCoverage\n\nCensus 2021 statistics are published for the whole of England and Wales.\nHowever, you can choose to filter areas by:\n\n- country - for example, Wales\n- region - for example, London\n- local authority - for example, Cornwall\n- health area – for example, Clinical Commissioning Group\n- statistical area - for example, MSOA or LSOA\n\n\nProtecting personal data\n\nSometimes we need to make changes to data if it is possible to identify\nindividuals. This is known as statistical disclosure control.\n\nIn Census 2021, we:\n- swapped records (targeted record swapping), for example, if a household \n  was likely to be identified in datasets because it has unusual\n  characteristics, we swapped the record with a similar one from a nearby\n  small area (very unusual households could be swapped with one in a nearby\n  local authority)\n- added small changes to some counts (cell key perturbation), for example,\n  we might change a count of four to a three or a five – this might make\n  small differences between tables depending on how the data are broken down\n  when we applied perturbation\n\nRead more in Section 5 of our article Design for Census 2021.\n\nDimensions:\n\n\tID: id 1\n\n\tLabel: Label 1\n\n\tDescription: description 1\n\n\tNumber Of Options: 1\n\n\tQuality Statement: quality statement 1\n\n\n\tID: id 2\n\n\tLabel: Label 2\n\n\tDescription: description 2\n\n\tNumber Of Options: 1\n\n\tQuality Statement: quality statement 2\n\n"

	Convey("Given a metadata struct", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "04-05-2023 10:00:003",
				Dimensions: []dataset.VersionDimension{
					{
						ID:                   "id 1",
						Label:                "Label 1",
						Description:          "description 1",
						NumberOfOptions:      1,
						QualityStatementText: "quality statement 1",
					},
					{
						ID:                   "id 2",
						Label:                "Label 2",
						Description:          "description 2",
						NumberOfOptions:      1,
						QualityStatementText: "quality statement 2",
					},
				},
				IsBasedOn: &dataset.IsBasedOn{
					Type: "type",
					ID:   "id",
				},
				Downloads: map[string]dataset.Download{
					"xlsx": {
						URL:  "xlsx url",
						Size: "xlsx size",
					},
					"csv": {
						URL:  "csv url",
						Size: "csv size",
					},
					"csvw": {
						URL:  "csvw url",
						Size: "csvw size",
					},
					"txt": {
						URL:  "txt url",
						Size: "txt size",
					},
				},
			},
			DatasetDetails: dataset.DatasetDetails{
				Title:       "title",
				Description: "description",
			},
		}

		Convey("When the NewMetadataCustom function is called", func() {
			timeNow := time.Now()
			expectedText = strings.Replace(expectedText, "1 Jan 2000", timeNow.Format("01-02-2006 15:04:05"), 1)
			meta := string(NewMetadataCustom(m, filterOutputID, downloadServiceURL))
			Convey("Then the returned value should be as expected", func() {
				So(meta, ShouldResemble, expectedText)
			})
		})
	})
}
