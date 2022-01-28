package csvw

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	. "github.com/smartystreets/goconvey/convey"
)

var fileURL = "ons/file.csv"
var apiURL = "api.example.com"

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
			csvw := New(m, fileURL)

			Convey("Then the values should be set to the expected fields", func() {
				So(csvw.Context, ShouldEqual, "http://www.w3.org/ns/csvw")
				So(csvw.Title, ShouldEqual, m.Title)
				So(csvw.Description, ShouldEqual, m.Description)
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
			data, err := Generate(ctx, m, fileURL, fileURL, apiURL)

			Convey("Then results should be returned with no errors", func() {
				So(data, ShouldHaveLength, 569)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given metadata that does not include a dimension", t, func() {
		m := &dataset.Metadata{
			Version: dataset.Version{
				ReleaseDate: "1 Jan 2000",
			},
			DatasetDetails: dataset.DatasetDetails{},
		}

		Convey("When the Generate csvw function is called", func() {
			data, err := Generate(ctx, m, fileURL, fileURL, apiURL)

			Convey("Then results should be returned with no errors", func() {
				So(data, ShouldHaveLength, 0)
				So(errors.Is(err, errMissingDimensions), ShouldBeTrue)
			})
		})
	})
}
