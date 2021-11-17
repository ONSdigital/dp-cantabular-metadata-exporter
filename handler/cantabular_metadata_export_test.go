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

	Convey("Succesfully return a string representing the vault path", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		var instanceID = "instance_id_test"
		h := NewCantabularMetadataExport(*cfg, nil, nil, nil)
		vaultPath := h.generateVaultPath(instanceID)
		expectedVaultPath := "secret/shared/psk/instance_id_test.txt"
		So(vaultPath, ShouldResemble, expectedVaultPath)

	})

}
