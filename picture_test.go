package flacpicture

import (
	"testing"

	httpclient "github.com/ddliu/go-httpclient"
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
