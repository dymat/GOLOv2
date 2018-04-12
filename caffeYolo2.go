package main

import (
	"gocv.io/x/gocv"
	"fmt"
	"image"
	"time"
)

func caffeWorker(img_chan chan *gocv.Mat, res_chan chan DetectionResult) {

	proto := "tiny-yolo-voc.prototxt"
	model := "tiny-yolo-voc.caffemodel"

	// open DNN classifier
	net := gocv.ReadNetFromCaffe(proto, model)
	if net.Empty() {
		fmt.Printf("Error reading network model from : %v %v\n", proto, model)
		return
	}
	defer net.Close()

	img := gocv.NewMat()
	defer img.Close()


	blob := gocv.NewMat()
	defer blob.Close()

	prob := gocv.NewMat()
	defer prob.Close()

	probMat := gocv.NewMat()
	defer probMat.Close()


	for item := range(img_chan) {
		if item.Empty(){
			continue
		}

		img = item.Clone()
		blob = gocv.BlobFromImage(img, 1.0/255.0, image.Pt(416, 416), gocv.NewScalar(0, 0, 0, 0), true, false)

		net.SetInput(blob, "data")

		prob = net.Forward("layer15-conv")
		probMat := prob.Reshape(1,1)

		boxes := regionLayer(probMat, true, float32(img.Rows()), float32(img.Cols()))
		time.Sleep(time.Millisecond*30)

		res_chan <- DetectionResult{boxes:boxes, img: img}
	}
}
