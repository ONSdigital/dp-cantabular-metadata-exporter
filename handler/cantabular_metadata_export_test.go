package handler

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

var eventTest = event.CSVCreated{
	DatasetID: "test_id",
	Edition:   "test-edition",
	Version:   "1",
}

func TestGenerateTextFilename(t *testing.T) {

	Convey("Succesfully return a string representing the .txt filename", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		h := NewCantabularMetadataExport(*cfg, nil, nil, nil)
		filename := h.generateTextFilename(&eventTest)
		expectedFilename := "datasets/test_id-test-edition-1.txt"
		So(filename, ShouldResemble, expectedFilename)

	})

}

func TestGenerateVaultPath(t *testing.T) {

	Convey("Succesfully return a string representing the vault path with given file extension", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		h := NewCantabularMetadataExport(*cfg, nil, nil, nil)

		vaultPath := h.generateVaultPath(&eventTest, "foo")
		expectedVaultPath := "secret/shared/psk/test_id-test-edition-1.foo"
		So(vaultPath, ShouldResemble, expectedVaultPath)

		vaultPath = h.generateVaultPath(&eventTest, "bar")
		expectedVaultPath = "secret/shared/psk/test_id-test-edition-1.bar"
		So(vaultPath, ShouldResemble, expectedVaultPath)

	})

}
