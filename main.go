package main

import (
	"gocv.io/x/gocv"
	"fmt"
	"image"
	"github.com/hybridgroup/go-ncs"
	"io/ioutil"
)

type DetectionResult struct {
	boxes []BBox
	img gocv.Mat
}

type BBox struct {
	topleft image.Point
	bottomright image.Point
	center image.Point
	label int
	confidence float32
}


func main() {

	useCaffeOrNCS := "caffe" // "ncs" or "caffe"

	// Setup Channels
	img_chan := make(chan *gocv.Mat, 1)
	res_chan := make(chan DetectionResult, 1)


	// Start Camera / Image Producer
	go producer(img_chan)


	// Start Worker (caffe or ncs)
	if useCaffeOrNCS == "caffe" {
		go caffeWorker(img_chan, res_chan)
	}


	if useCaffeOrNCS == "ncs" {
		// bootstrapping NCS devices
		const numNCSDevices= 1
		const graphFileName= "tiny-yolov2-voc.ncsmodel"


		for deviceID := 0; deviceID < numNCSDevices; deviceID++ {
			res, name := ncs.GetDeviceName(deviceID)
			if res != ncs.StatusOK {
				fmt.Printf("NCS Error: %v\n", res)
				return
			}

			fmt.Println(deviceID, "NCS: " + name)

			// open NCS device
			fmt.Println("Opening NCS device " + name + "...")
			status, s := ncs.OpenDevice(name)
			if status != ncs.StatusOK {
				fmt.Printf("NCS Error: %v\n", status)
				return
			}
			defer s.CloseDevice()

			// load precompiled graph file in NCS format
			data, err := ioutil.ReadFile(graphFileName)
			if err != nil {
				fmt.Println("Error opening graph file:", err)
				return
			}

			// allocate graph on NCS stick
			fmt.Println("Allocating graph...")
			allocateStatus, graph := s.AllocateGraph(data)
			if allocateStatus != ncs.StatusOK {
				fmt.Printf("NCS Error: %v\n", allocateStatus)
				return
			}
			defer graph.DeallocateGraph()

			// Start NCS Worker
			go ncsWorker(graph, img_chan, res_chan)
		}
	}


	// Display Detections
	showImg := gocv.NewMat()

	window := gocv.NewWindow("Hello")
	defer window.Close()

	imgCounter := 0

	for res := range(res_chan) {

		showImg = res.img

		fmt.Println(fmt.Sprintf("\nImage #%d", imgCounter))
		imgCounter++

		for _, box := range(res.boxes) {
			b := image.Rectangle{image.Point{X:box.topleft.X, Y:box.topleft.Y},	image.Point{X:box.bottomright.X, Y:box.bottomright.Y}}
			gocv.Rectangle(&showImg, b, colors[box.label % len(colors)], 2)

			textPos := box.topleft
			textPos.Y = textPos.Y-10
			gocv.PutText(&showImg, classNames[box.label], textPos, 1, 1.0, colors[box.label % len(colors)], 2)

			fmt.Println(classNames[box.label], box.confidence)
		}

		window.IMShow(showImg)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}



func producer(img_chan chan <- *gocv.Mat) {
	webcam, _ := gocv.VideoCaptureDevice(0)
	webcam.Set(3, 800)
	webcam.Set(4, 450)
	webcam.Set(5, 30)


	for {
		img := gocv.NewMat()
		webcam.Read(&img)

		// Just for Debugging: Override Camera with image from disk
		img = gocv.IMRead("person.jpg", 1)
		img_chan <- &img
	}
}

