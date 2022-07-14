package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func main() {
	// if len(os.Args) < 3 {
	// 	fmt.Println("How to run:\n\tfacedetect [camera ID] [classifier XML file]")
	// 	return
	// }

	// parse args
	deviceID := 0
	faceXmlFile := "haarcascade_frontalface_default.xml"
	eyeXmlFile := "haarcascade_eye.xml"

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load faceClassifier to recognize faces
	faceClassifier := gocv.NewCascadeClassifier()
	defer faceClassifier.Close()

	if !faceClassifier.Load(faceXmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", faceXmlFile)
		return
	}
	// load faceClassifier to recognize eyes
	eyeClassifier := gocv.NewCascadeClassifier()
	defer eyeClassifier.Close()

	if !eyeClassifier.Load(eyeXmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", eyeXmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := faceClassifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)

			imgFace := img.Region(r)
			eyes := eyeClassifier.DetectMultiScale(imgFace)
			fmt.Printf("found %d faces\n", len(eyes))
			for _, e := range eyes {
				gocv.Rectangle(&imgFace, e, color.RGBA{0, 255, 0, 0}, 1)
				size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
				pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
				gocv.PutText(&imgFace, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
			}
			imgFace.Close()

		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
