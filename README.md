# GOLOv2
YOLOv2 for Golang

This projects implements the yolov2 (https://pjreddie.com/darknet/yolov2/) RegionLayer in Go. It is heavily inspired by duangenquan's C++-RegionLayer implementation (https://github.com/duangenquan/YoloV2NCS).

This projects makes use of `gocv` (https://gocv.io) and `go-ncs` (https://github.com/hybridgroup/go-ncs/), both from hybridgroup (https://github.com/hybridgroup).

It comes with a **tiny-yolo caffe model** which I derived from original weights (https://pjreddie.com/media/files/yolov2-tiny-voc.weights) with this *darknet2caffe converter*: https://github.com/marvis/pytorch-caffe-darknet-convert. It also comes with a **Movidius NCS model version of tiny-yolo** which I compiled from the converted caffe model.


# Setup

1. Install `gocv` as described on https://gocv.io/getting-started/
2. Install `go-ncs` as described on https://github.com/hybridgroup/go-ncs
3. Plug in your Movidius Neural Compute Stick
4. `$ git clone git@github.com:dymat/GOLOv2.git`
5. `$ cd GOLOv2`
6. `go run *.go`
