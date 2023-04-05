package csvw

import (
	"context"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	. "github.com/smartystreets/goconvey/convey"
)

var fileURL = "ons/file.csv"
var apiURL = "api.example.com"
var externalPrefixURL = "external.prefixurl.com"
var filterOutputID = "filter-output-id"
var downloadServiceURL = "download-service-url"
var isCustom = false

var ctx = context.Background()

func TestNew(t *testing.T) {
	Convey("Given a complete metadata struct", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "1 Jan 2000",
				CSVHeader:   []string{"cantabular_table", "sex", "age"},
			},
			DatasetDetails: dataset.DatasetDetails{
				Title:       "title",
				Description: "description",
			},
		}

		Convey("When the New csvw function is called", func() {
			csvw := New(m, fileURL, externalPrefixURL, filterOutputID, downloadServiceURL)

			Convey("Then the values should be set to the expected fields", func() {
				So(csvw.Context, ShouldEqual, "http://www.w3.org/ns/csvw")
				So(csvw.Title, ShouldEqual, m.Title)
				So(csvw.Description, ShouldEqual, m.Description)
			})
		})
	})
}

func TestNewCustom(t *testing.T) {
	Convey("Given a complete metadata struct", t, func() {
		time := time.Now()
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: time.Format("01-02-2006 15:04:05"),
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
				Title: "Label 1 and Label 2",
			},
		}

		Convey("When the NewCustom csvw function is called", func() {
			csvw := NewCustom(m, fileURL, externalPrefixURL, filterOutputID, downloadServiceURL)

			Convey("Then the values should be set to the expected fields", func() {
				So(csvw.Context, ShouldEqual, "http://www.w3.org/ns/csvw")
				So(csvw.Title, ShouldEqual, m.Title)
				So(csvw.Issued, ShouldEqual, m.ReleaseDate)
			})
		})
	})
}

func TestFormatAboutURL(t *testing.T) {
	Convey("Given a valid domain config and url", t, func() {
		domain := "http://api.example.com/v1"
		url := "http://localhost:22000/datasets/1/editions/2/version/3/metadata"

		Convey("When the formatAboutURL function is called", func() {
			url, err := formatAboutURL(url, domain)

			Convey("Then the returned values should be as expected", func() {
				So(url, ShouldEqual, "http://api.example.com/v1/datasets/1/editions/2/version/3/metadata")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGenerate(t *testing.T) {

	Convey("Given metadata that includes a dimension", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "1 Jan 2000",
				Dimensions: []dataset.VersionDimension{
					{
						Name: "geography",
						Links: dataset.Links{
							Self: dataset.Link{
								URL: "api/versions/self",
							},
						},
						Description: "areas included in dataset",
						Label:       "Geographic areas",
					},
				},
				CSVHeader: []string{"cantabular_table", "a", "b", "c", "d"},
			},
			DatasetDetails: dataset.DatasetDetails{},
		}

		Convey("When the Generate csvw function is called", func() {
			data, err := Generate(ctx, m, fileURL, fileURL, apiURL, externalPrefixURL, filterOutputID, downloadServiceURL, isCustom)

			Convey("Then results should be returned with no errors", func() {
				So(data, ShouldHaveLength, 536)
				So(err, ShouldBeNil)
			})
		})
	})

	// Convey("Given metadata that does not include a dimension", t, func() {
	// 	m := &dataset.Metadata{
	// 		Version: dataset.Version{
	// 			ReleaseDate: "1 Jan 2000",
	// 		},
	// 		DatasetDetails: dataset.DatasetDetails{},
	// 	}

	// 	Convey("When the Generate csvw function is called", func() {
	// 		data, err := Generate(ctx, m, fileURL, fileURL, apiURL)

	// 		Convey("Then results should be returned with no errors", func() {
	// 			So(data, ShouldHaveLength, 0)
	// 			So(errors.Is(err, errMissingDimensions), ShouldBeTrue)
	// 		})
	// 	})
	// })
}
