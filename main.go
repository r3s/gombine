package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// imageData struct fold holding each input image and related data
type imageData struct {
	img    image.Image
	width  int
	height int
	path   string
}

// usage function to print usage on -help
func usage() {
	fmt.Println("gombine [options] <file1> <file2> ...")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println("Ex: gombine -format=png -side=bottom -out=go.png 1.png 2.png")
	os.Exit(0)
}

// getImageData function to populate a imageData object with input image details
// Takes the image, and filename as arguments
// Returns the filled imageData object and an error if any
func getImageData(img *image.Image, filename string) (imageData, error) {
	imd := new(imageData)
	imd.img = *img
	imd.path = filename
	h, w, err := getDim(imd)
	imd.height, imd.width = h, w
	if err != nil {
		return *imd, err
	}

	return *imd, nil

}

// getDim function to get the dimensions of an input image
// Takes imageData as argument
// Return height, width and error if any
func getDim(imd *imageData) (int, int, error) {
	f, err := os.Open(imd.path)
	if err != nil {
		return -1, -1, err
	}
	defer f.Close()
	// Decode config of image to get height and width
	config, _, err := image.DecodeConfig(f)
	if err != nil {
		return -1, -1, err
	}
	return config.Height, config.Width, nil
}

// getTotalDim function to get the total height and width
// i.e, sum of widths and heights of all input images
// Takes the array of imageData as argument
// Returns total height, width and error if any
func getTotalDim(images []*imageData) (int, int, error) {

	height, width := 0, 0
	// Loop through images and add the height and width
	for _, imd := range images {
		height = height + imd.height
		width = width + imd.width
	}

	if height == 0 && width == 0 {
		return height, width, errors.New("Total Height and Width cannot be 0")
	}

	return height, width, nil
}

// getMaxDim function to get the maximum width and height from
// all the input images. Takes imageData array as argument
// Returns max height, width and error if any
func getMaxDim(images []*imageData) (int, int, error) {
	maxh, maxw := 0, 0
	// Loop through images to find the largest height and width
	for _, imd := range images {
		if imd.height > maxh {
			maxh = imd.height
		}
		if imd.width > maxw {
			maxw = imd.width
		}
	}
	return maxh, maxw, nil
}

// processImages function to loop through all images in the imageData array
// and calculate the total height, width and max height, width.
// Finally calls makeImage to create the image
// Takes the array of imageData, format and side as arguments
func processImages(images []*imageData, format, side, outfile string) {
	th, tw, err := getTotalDim(images)
	if err != nil {
		log.Fatal(err)
	}
	maxh, maxw, err := getMaxDim(images)
	if err != nil {
		log.Fatal(err)
	}
	// Create the output image
	err = makeImage(th, tw, maxh, maxw, images, format, side, outfile)
	if err != nil {
		log.Fatal(err)
	}
}

// makeImage function to create the combined image from all the input images
// Takes total height, width, max height, width, input images, format to
// encode and the side to which the images are to be combined as arguments
// Returns error if any
func makeImage(th, tw, maxh, maxw int, images []*imageData, format, side, outfile string) error {
	var img *image.RGBA
	posx, posy := 0, 0
	switch side {
	case "bottom":
		img = image.NewRGBA(image.Rect(0, 0, maxw, th))
		for _, imd := range images {
			r := image.Rect(posx, posy, posx+imd.width, posy+imd.height)
			draw.Draw(img, r, imd.img, image.Point{0, 0}, draw.Over)
			posy = posy + imd.height
		}

	case "right":
		img = image.NewRGBA(image.Rect(0, 0, tw, maxh))
		for _, imd := range images {
			r := image.Rect(posx, posy, posx+imd.width, posy+imd.height)
			draw.Draw(img, r, imd.img, image.Point{0, 0}, draw.Over)
			posx = posx + imd.width
		}
	default:
		return errors.New("Please choose bottom or right for side")
	}

	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer out.Close()
	// Encode the image in the format given as input
	if format == "png" {
		err = png.Encode(out, img)
	} else if format == "jpg" {
		// jpeg quality is set as 90.
		option := jpeg.Options{Quality: 90}
		err = jpeg.Encode(out, img, &option)
	} else {
		return errors.New("Please choose jpg or png for format")
	}
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetPrefix("GOMBINE")
	// format, side and out commandline arguments
	format := flag.String("format", "png", "Output file format, supports png/jpg")
	side := flag.String("side", "bottom", "bottom or right")
	outfile := flag.String("out", "combined.png", "Output file name")
	// Usage message
	flag.Usage = usage

	flag.Parse()
	images := []*imageData{}

	// flag.Args() contain the list of input files
	// The maximum number of input files is set as 10
	if len(flag.Args()) > 0 && len(flag.Args()) < 10 {
		for _, filename := range flag.Args() {
			res, err := os.Stat(filename)

			if err != nil {
				log.Fatal(err)
			}

			if !res.IsDir() {

				fimg, err := os.Open(filename)
				if err != nil {
					log.Fatal("Unable to open file", filename)
				}
				defer fimg.Close()
				// Decode the file to get the image data
				img, _, err := image.Decode(fimg)
				if err != nil {
					log.Fatal("Unable to decode ", filename)
				}
				// Fill image data object
				imd, err := getImageData(&img, filename)
				if err != nil {
					log.Fatal(err)
				}
				// Append to imadeData array
				images = append(images, &imd)

			}
		}
		// Start processing images
		processImages(images, *format, *side, *outfile)

	} else {
		log.Fatal("No files given")
	}
}
