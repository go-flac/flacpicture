package flacpicture

import (
	"bytes"
	"errors"
	"image/jpeg"

	flac "github.com/go-flac/go-flac"
)

// PictureType defines the type of this image
type PictureType uint32

const (
	PictureTypeOther PictureType = iota
	PictureTypeFileIcon
	PictureTypeOtherIcon
	PictureTypeFrontCover
	PictureTypeBackCover
	PictureTypeLeaflet
	PictureTypeMedia
	PictureTypeLeadArtist
	PictureTypeArtist
	PictureTypeConductor
	PictureTypeBand
	PictureTypeComposer
	PictureTypeLyricist
	PictureTypeRecordingLocation
	PictureTypeDuringRecording
	PictureTypeDuringPerformance
	PictureTypeScreenCapture
	PictureTypeBrightColouredFish
	PictureTypeIllustration
	PictureTypeBandArtistLogotype
	PictureTypePublisherStudioLogotype
)

type MetadataBlockPicture struct {
	PictureType       PictureType
	MIME              string
	Description       string
	Width             uint32
	Height            uint32
	ColorDepth        uint32
	IndexedColorCount uint32
	ImageData         []byte
}

// Marshal encodes the PictureBlock to a standard FLAC MetaDataBloc to be accepted by go-flac
func (c *MetadataBlockPicture) Marshal() flac.MetaDataBlock {
	res := bytes.NewBuffer([]byte{})
	res.Write(encodeUint32(uint32(c.PictureType)))
	res.Write(encodeUint32(uint32(len(c.MIME))))
	res.Write([]byte(c.MIME))
	res.Write(encodeUint32(uint32(len(c.Description))))
	res.Write([]byte(c.Description))
	res.Write(encodeUint32(c.Width))
	res.Write(encodeUint32(c.Height))
	res.Write(encodeUint32(c.ColorDepth))
	res.Write(encodeUint32(c.IndexedColorCount))
	res.Write(encodeUint32(uint32(len(c.ImageData))))
	res.Write(c.ImageData)
	return flac.MetaDataBlock{
		Type: flac.Picture,
		Data: res.Bytes(),
	}
}

// NewFromImageData generates a MetadataBlockPicture from image data
func NewFromImageData(pictype PictureType, description string, imgdata []byte, mime string) (*MetadataBlockPicture, error) {
	res := new(MetadataBlockPicture)
	res.PictureType = pictype
	res.Description = description
	res.MIME = mime
	switch mime {
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imgdata))
		if err != nil {
			return nil, err
		}
		res.IndexedColorCount = uint32(0)
		size := img.Bounds()
		res.Width = uint32(size.Max.X)
		res.Height = uint32(size.Max.Y)
		res.ColorDepth = uint32(8)
		res.ImageData = imgdata
		return res, nil
	default:
		return nil, errors.New("Unsupported MIME")
	}
}
