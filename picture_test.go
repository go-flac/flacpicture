package flacpicture

import (
	"archive/zip"
	"bytes"
	"testing"

	httpclient "github.com/ddliu/go-httpclient"
	flac "github.com/go-flac/go-flac"
)

func TestPNGPictureDecode(t *testing.T) {
	imgres, err := httpclient.Begin().Get("https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	imgdata, err := imgres.ReadAll()
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	pic, err := NewFromImageData(PictureTypeArtist, "test description", imgdata, "image/png")
	if err != nil {
		t.Errorf("Error while constructing image data: %s", err.Error())
		t.Fail()
	}

	if pic.MIME != "image/png" {
		t.Errorf("MIME decode error: got %s", pic.MIME)
		t.Fail()
	}

	if pic.Height != 600 || pic.Width != 800 {
		t.Errorf("JPEG size error: got %dx%d", pic.Width, pic.Height)
		t.Fail()
	}
}
func TestJPEGPictureDecode(t *testing.T) {
	imgres, err := httpclient.Begin().Get("https://jpeg.org/images/jpeg-home.jpg")
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	imgdata, err := imgres.ReadAll()
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	pic, err := NewFromImageData(PictureTypeArtist, "test description", imgdata, "image/jpeg")
	if err != nil {
		t.Errorf("Error while constructing image data: %s", err.Error())
		t.Fail()
	}

	if pic.MIME != "image/jpeg" {
		t.Errorf("MIME decode error: got %s", pic.MIME)
		t.Fail()
	}

	if pic.Height != 298 || pic.Width != 690 {
		t.Errorf("JPEG size error: got %dx%d", pic.Width, pic.Height)
		t.Fail()
	}
}

func TestPictureModify(t *testing.T) {
	imgres, err := httpclient.Begin().Get("https://jpeg.org/images/jpeg-home.jpg")
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	imgdata, err := imgres.ReadAll()
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	pic, err := NewFromImageData(PictureTypeArtist, "test description", imgdata, "image/jpeg")
	if err != nil {
		t.Errorf("Error while constructing image data: %s", err.Error())
		t.Fail()
	}

	pic.Description = "another description"

	pic, err = ParseFromMetaDataBlock(pic.Marshal())
	if err != nil {
		t.Errorf("Error while parsing modified image data: %s", err.Error())
		t.Fail()
	}

	if pic.Description != "another description" {
		t.Errorf("description unepected: %s", pic.Description)
		t.Fail()
	}
}

func TestJPEGFromExistingFLAC(t *testing.T) {
	zipres, err := httpclient.Begin().Get("http://helpguide.sony.net/high-res/sample1/v1/data/Sample_BeeMoved_96kHz24bit.flac.zip")
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	zipdata, err := zipres.ReadAll()
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	zipfile, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		t.Errorf("Error while decompressing test file: %s", err.Error())
		t.FailNow()
	}
	if zipfile.File[0].Name != "Sample_BeeMoved_96kHz24bit.flac" {
		t.Errorf("Unexpected test file content: %s", zipfile.File[0].Name)
		t.FailNow()
	}

	flachandle, err := zipfile.File[0].Open()
	if err != nil {
		t.Errorf("Failed to decompress test file: %s", err)
		t.FailNow()
	}

	f, err := flac.ParseBytes(flachandle)
	if err != nil {
		t.Errorf("Failed to parse flac file: %s", err)
		t.FailNow()
	}

	var pic *MetadataBlockPicture
	for _, meta := range f.Meta {
		if meta.Type == flac.Picture {
			pic, err = ParseFromMetaDataBlock(*meta)
			if err != nil {
				t.Errorf("Error while parsing metadata image: %s", err.Error())
				t.Fail()
			}
		}
	}
	if pic.PictureType != PictureTypeFrontCover {
		t.Error("Picture type does not match")
		t.Fail()
	}
	if pic.MIME != "image/jpeg" {
		t.Errorf("Picture MIME does not match: %s", pic.MIME)
		t.Fail()
	}
}
