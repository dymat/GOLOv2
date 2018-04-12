package main

import (
	ncs "github.com/hybridgroup/go-ncs"
	"fmt"
	"gocv.io/x/gocv"
	"image"
)

func ncsWorker(graph *ncs.Graph, image_chan chan *gocv.Mat, result_chan chan <- DetectionResult) {

	dims := image.Pt(416, 416)
	scaleMat := makeScaleMat(1.0/255.0, dims)

	errors := 0

	// Worker-Loop
	for img := range image_chan {

		if errors > 9 {
			return
		}

		// convert image to format needed by NCS
		fp32Image := prepareImage(img, &scaleMat, dims)
		fp16Blob := fp32Image.ConvertFp16()

		// load image tensor into graph on NCS stick
		loadStatus := graph.LoadTensor(fp16Blob.ToBytes())
		if loadStatus != ncs.StatusOK {
			fmt.Println("Error loading tensor data:", loadStatus)
			errors++
			continue
		}

		// get result from NCS stick in fp16 format
		resultStatus, data := graph.GetResult()
		if resultStatus != ncs.StatusOK {
			fmt.Println("Error getting results:", resultStatus)
			errors++
			continue
		}


		errors = 0

		// convert results from fp16 back to float32
		fp16Results := gocv.NewMatFromBytes(1, len(data)/2, gocv.MatTypeCV16S, data)
		fp32Results := fp16Results.ConvertFp16()

		boxes := regionLayer(fp32Results, false, float32(img.Rows()), float32(img.Cols()))
		result_chan <- DetectionResult{boxes: boxes, img: *img}

	}
}




func makeScaleMat(scalar float32, dims image.Point) gocv.Mat {
	scaleMat := gocv.NewMatWithSize(1, dims.Y*dims.X*3, gocv.MatTypeCV32F)
	for i :=0; i <dims.Y*dims.X*3; i++ {
		scaleMat.SetFloatAt(0, i, scalar)
	}

	return scaleMat
}

func prepareImage(img *gocv.Mat, scaleMat *gocv.Mat, dims image.Point) gocv.Mat {
	resized := gocv.NewMat()
	gocv.Resize(*img, &resized, dims, 0, 0, gocv.InterpolationLinear)

	fp32Image := gocv.NewMat()
	resized.ConvertTo(&fp32Image, gocv.MatTypeCV32F)

	fp32Image = fp32Image.Reshape(1,1)

	scaled := gocv.NewMat()
	gocv.Multiply(fp32Image, *scaleMat, &scaled)
	fp32Image = scaled.Reshape(3, dims.Y)

	return fp32Image
}
