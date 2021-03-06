package handler

import (
	"testing"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	. "github.com/smartystreets/goconvey/convey"
)

var eventTest = event.CSVCreated{
	DatasetID: "test_id",
	Edition:   "test-edition",
	Version:   "1",
}

func TestGenerateTextFilename(t *testing.T) {

	Convey("Succesfully return a string representing the .txt filename", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		h := NewCantabularMetadataExport(*cfg, nil, nil, nil, nil, nil)
		filename := h.generateTextFilename(&eventTest)
		expectedFilename := "datasets/test_id-test-edition-1.txt"
		So(filename, ShouldResemble, expectedFilename)

	})

}

func TestGenerateVaultPath(t *testing.T) {

	Convey("Succesfully return a string representing the vault path with given file extension", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		h := NewCantabularMetadataExport(*cfg, nil, nil, nil, nil, nil)

		Convey("Vault path ends in 'foo'", func() {
			vaultPath := h.generateVaultPath(&eventTest, "foo")
			expectedVaultPath := "secret/shared/psk/test_id-test-edition-1.foo"
			So(vaultPath, ShouldResemble, expectedVaultPath)
		})

		Convey("Vault path ends in 'bar'", func() {
			vaultPath := h.generateVaultPath(&eventTest, "bar")
			expectedVaultPath := "secret/shared/psk/test_id-test-edition-1.bar"
			So(vaultPath, ShouldResemble, expectedVaultPath)
		})
	})

}

func TestGenerateDownloadURL(t *testing.T) {

	Convey("Succesfully return a string representing the download URL with given file extension", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		h := NewCantabularMetadataExport(*cfg, nil, nil, nil, nil, nil)

		Convey("Download URL ends in 'foo'", func() {
			vaultPath := h.generateDownloadURL(&eventTest, "foo")
			expectedVaultPath := "http://localhost:23600/downloads/datasets/test_id/editions/test-edition/versions/1.foo"
			So(vaultPath, ShouldResemble, expectedVaultPath)
		})

		Convey("Download URL ends in 'bar'", func() {
			vaultPath := h.generateDownloadURL(&eventTest, "bar")
			expectedVaultPath := "http://localhost:23600/downloads/datasets/test_id/editions/test-edition/versions/1.bar"
			So(vaultPath, ShouldResemble, expectedVaultPath)
		})
	})

}
