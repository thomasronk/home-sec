package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"gocv.io/x/gocv"
)

const MinimumArea = 3000

func main() {
	deviceID := os.Args[1]
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening Video device: %v", deviceID)
	}
	//window := gocv.NewWindow("Hello")
	//img := gocv.NewMat()
	// for {
	// 	webcam.Read(&img)
	// 	window.IMShow(img)
	// 	window.WaitKey(1)
	// }
	defer webcam.Close()

	window := gocv.NewWindow("Motion Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()
	status := "Ready"

	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		status = "Ready"
		statusColor := color.RGBA{0, 255, 0, 0}

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(img, &imgDelta)

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgThresh, kernel)

		// now find contours
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i, c := range contours {
			area := gocv.ContourArea(c)
			if area < MinimumArea {
				continue
			}

			status = "Motion detected"
			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(c)
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		//window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
