package main

import "image/color"

/*
 * Constants
 */

const numClasses = 20
const N = 5
const size = numClasses + N
const w = 12
const h = 12
const blockwd float32 = 13
const numBoxes = h*w*N

var anchors = [2*N]float32 {1.08,1.19,  3.42,4.41,  6.63,11.38,  9.42,5.11,  16.62,10.52}
//var anchors = [2*N]float32 {0.57273, 0.677385, 1.87446, 2.06253, 3.33843, 5.47434, 7.88282, 3.52778, 9.77052, 9.16828}

const thresh = 0.2
const nms_threshold = 0.4

var classNames = [20]string{"aeroplane", "bicycle", "bird", "boat", "bottle", "bus", "car", "cat","chair","cow","diningtable","dog","horse","motorbike","person","pottedplant","sheep","sofa","train","tvmonitor"}
var colors = [20]color.RGBA{
	color.RGBA{230, 25, 75, 0},
	color.RGBA{60, 180, 75, 0},
	color.RGBA{255, 225, 25, 0},
	color.RGBA{0, 130, 200, 0},
	color.RGBA{245, 130, 48, 0},
	color.RGBA{145, 30, 180, 0},
	color.RGBA{70, 240, 240, 0},
	color.RGBA{240, 50, 230, 0},
	color.RGBA{210, 245, 60, 0},
	color.RGBA{250, 190, 190, 0},
	color.RGBA{0, 128, 128, 0},
	color.RGBA{230, 190, 255, 0},
	color.RGBA{170, 110, 40, 0},
	color.RGBA{255, 250, 200, 0},
	color.RGBA{128, 0, 0, 0},
	color.RGBA{170, 255, 195, 0},
	color.RGBA{128, 128, 0, 0},
	color.RGBA{255, 215, 180, 0},
	color.RGBA{0, 0, 128, 0},
	color.RGBA{128, 128, 128, 0},
}
