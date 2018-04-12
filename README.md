# GOLOv2
YOLOv2 for Golang

This projects implements the yolov2 (https://pjreddie.com/darknet/yolov2/) RegionLayer in Go. It is heavily inspired by duangenquan's C++-RegionLayer implementation (https://github.com/duangenquan/YoloV2NCS).

This projects makes use of `gocv` (https://gocv.io) and `go-ncs` (https://github.com/hybridgroup/go-ncs/), both from hybridgroup (https://github.com/hybridgroup).

It comes with a **tiny-yolo caffe model** which I derived from original weights (https://pjreddie.com/media/files/yolov2-tiny-voc.weights) with this *darknet2caffe converter*: https://github.com/marvis/pytorch-caffe-darknet-convert. It also comes with a **Movidius NCS model version of tiny-yolo** which I compiled from the converted caffe model.


# !Important!
The caffe version works well -- **the NCS version doesn't**. This is one of the main reasons I published my code: **I hope for some inpirations from others to make the NCS version work.**
