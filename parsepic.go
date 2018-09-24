package flacpicture

import (
	"bytes"
	"image/jpeg"
)

// ParsePicture decodes the image and inflated the Width, Height, ColorDepth, IndexedColorCount fields. This is called automatically by NewFromImageData
func (c *MetadataBlockPicture) ParsePicture() error {
	switch c.MIME {
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(c.ImageData))
		if err != nil {
			return err
		}
		c.IndexedColorCount = uint32(0)
		size := img.Bounds()
		c.Width = uint32(size.Max.X)
		c.Height = uint32(size.Max.Y)
		c.ColorDepth = uint32(8)
	default:
		return ErrorUnsupportedMIME
	}
	return nil
}
