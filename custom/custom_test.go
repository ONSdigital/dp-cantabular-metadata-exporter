package custom_test

import (
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/custom"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateCustomTitle(t *testing.T) {
	Convey("Given a handler", t, func() {
		Convey("given a set of 3 filter dimensions", func() {
			dims := []dataset.VersionDimension{
				{
					Label: "Label 1",
				},
				{
					Label: "Label 2",
				},
				{
					Label: "Label 3",
				},
			}

			Convey("when GenerateCustomTitle is called", func() {
				title := custom.GenerateCustomTitle(dims)
				expected := "Label 1, Label 2 and Label 3"
				So(title, ShouldResemble, expected)
			})
		})

		Convey("given a set of 2 filter dimensions", func() {
			dims := []dataset.VersionDimension{
				{
					Label: "Label 1",
				},
				{
					Label: "Label 2",
				},
			}

			Convey("when GenerateCustomTitle is called", func() {
				title := custom.GenerateCustomTitle(dims)
				expected := "Label 1 and Label 2"
				So(title, ShouldResemble, expected)
			})
		})

		Convey("given a set of 1 filter dimensions", func() {
			dims := []dataset.VersionDimension{
				{
					Label: "Label 1",
				},
			}

			Convey("when GenerateCustomTitle is called", func() {
				title := custom.GenerateCustomTitle(dims)
				expected := " and Label 1"
				So(title, ShouldResemble, expected)
			})
		})

		Convey("given 0 filter dimensions", func() {
			dims := []dataset.VersionDimension{}

			Convey("when GenerateCustomTitle is called", func() {
				title := custom.GenerateCustomTitle(dims)
				expected := ""
				So(title, ShouldResemble, expected)
			})
		})
	})
}
