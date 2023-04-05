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
	expectedText := "Title: \nDescription: description\nIssued: 01-02-2006 15:04:05\nNext Release: \nIdentifier: \nLanguage: English\nLatest Changes: []\nPeriodicity: \nDistribution:\nUnit of measure: \nLicense: \nCanonical Topic: \nSurvey: \nUsage Notes:\nDataset Links:\nEditions:\nHREF: \nLatest Version:\nHREF: \nID: \nSelf:\nHREF: \nVersion: 0\n\nArea Type\n\nCensus 2021 statistics are published for a number of different geographies.\nThese can be large, for example the whole of England, or small, for example\nan output area (OA), the lowest level of geography for which statistics are\nproduced.\n\nFor higher levels of geography, more detailed statistics can be produced.\nWhen a lower level of geography is used, such as output areas (which have\na minimum of 100 persons), the statistics produced have less detail. This is\nto protect the confidentiality of people and ensure that individuals or their\ncharacteristics cannot be identified.\n\n\nCoverage\n\nCensus 2021 statistics are published for the whole of England and Wales.\nHowever, you can choose to filter areas by:\n\n- country - for example, Wales\n- region - for example, London\n- local authority - for example, Cornwall\n- health area – for example, Clinical Commissioning Group\n- statistical area - for example, MSOA or LSOA\n\n\nProtecting personal data\n\nSometimes we need to make changes to data if it is possible to identify\nindividuals. This is known as statistical disclosure control.\n\nIn Census 2021, we:\n- swapped records (targeted record swapping), for example, if a household \n  was likely to be identified in datasets because it has unusual\n  characteristics, we swapped the record with a similar one from a nearby\n  small area (very unusual households could be swapped with one in a nearby\n  local authority)\n- added small changes to some counts (cell key perturbation), for example,\n  we might change a count of four to a three or a five – this might make\n  small differences between tables depending on how the data are broken down\n  when we applied perturbation\n\nRead more in Section 5 of our article Design for Census 2021.\n\nDimensions:\n\n\tID: \n\n\tLabel: Label 1\n\n\tDescription: \n\n\tNumber Of Options: 0\n\n\tQuality Statement: \n\n\n\tID: \n\n\tLabel: Label 2\n\n\tDescription: \n\n\tNumber Of Options: 0\n\n\tQuality Statement: \n\n"
	Convey("Given a metadata struct", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "01-02-2006 15:04:05",
				Dimensions: []dataset.VersionDimension{
					{
						Label: "Label 1",
					},
					{
						Label: "Label 2",
					},
				},
			},
			DatasetDetails: dataset.DatasetDetails{
				Description: "description",
			},
		}

		Convey("When the NewMetadata function is called", func() {
			meta := string(NewMetadata(m, filterOutputID, downloadServiceURL))
			Convey("Then the returned value should be as expected", func() {
				So(meta, ShouldResemble, expectedText)
			})
		})
	})
}

func TestNewMetadataCustom(t *testing.T) {
	expectedText := "Title: Label 1 and Label 2\nIssued: 1 Jan 2000\nLanguage: English\nDistribution:\n\nArea Type\n\nCensus 2021 statistics are published for a number of different geographies.\nThese can be large, for example the whole of England, or small, for example\nan output area (OA), the lowest level of geography for which statistics are\nproduced.\n\nFor higher levels of geography, more detailed statistics can be produced.\nWhen a lower level of geography is used, such as output areas (which have\na minimum of 100 persons), the statistics produced have less detail. This is\nto protect the confidentiality of people and ensure that individuals or their\ncharacteristics cannot be identified.\n\n\nCoverage\n\nCensus 2021 statistics are published for the whole of England and Wales.\nHowever, you can choose to filter areas by:\n\n- country - for example, Wales\n- region - for example, London\n- local authority - for example, Cornwall\n- health area – for example, Clinical Commissioning Group\n- statistical area - for example, MSOA or LSOA\n\n\nProtecting personal data\n\nSometimes we need to make changes to data if it is possible to identify\nindividuals. This is known as statistical disclosure control.\n\nIn Census 2021, we:\n- swapped records (targeted record swapping), for example, if a household \n  was likely to be identified in datasets because it has unusual\n  characteristics, we swapped the record with a similar one from a nearby\n  small area (very unusual households could be swapped with one in a nearby\n  local authority)\n- added small changes to some counts (cell key perturbation), for example,\n  we might change a count of four to a three or a five – this might make\n  small differences between tables depending on how the data are broken down\n  when we applied perturbation\n\nRead more in Section 5 of our article Design for Census 2021.\n\nDimensions:\n\n\tID: \n\n\tLabel: Label 1\n\n\tDescription: \n\n\tNumber Of Options: 0\n\n\tQuality Statement: \n\n\n\tID: \n\n\tLabel: Label 2\n\n\tDescription: \n\n\tNumber Of Options: 0\n\n\tQuality Statement: \n\n"
	time := time.Now()
	expectedText = strings.Replace(expectedText, "1 Jan 2000", time.Format("01-02-2006 15:04:05"), 1)
	Convey("Given a metadata struct", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "04-05-2023 10:00:003",
				Dimensions: []dataset.VersionDimension{
					{
						Label: "Label 1",
					},
					{
						Label: "Label 2",
					},
				},
			},
			DatasetDetails: dataset.DatasetDetails{
				Description: "description",
			},
		}

		Convey("When the NewMetadata function is called", func() {
			meta := string(NewMetadataCustom(m, filterOutputID, downloadServiceURL))
			Convey("Then the returned value should be as expected", func() {
				So(meta, ShouldResemble, expectedText)
			})
		})
	})
}
