package main

import (
	"gocv.io/x/gocv"
	"math"
	"image"
	"sort"
)


type Box struct {
	x float32
	y float32
	w float32
	h float32
	classProbs [numClasses]float32
	confidence float32
	currentClassIdx int
}

/*
 * Region Layer
 */


func regionLayer(predictions gocv.Mat, transposePredictions bool, img_height, img_width float32) []BBox {

	var data [w*h*5*(numClasses+5)]float32

	if transposePredictions {
		predictions = predictions.Reshape(1, 125)
		data = transpose(&predictions)
	} else {
		data = matToArray(&predictions)
	}


	var boxes []Box
	for i := 0; i < numBoxes; i++ {
		index := i * size
		var n = i % N
		var row = float32((i/N) / h)
		var col = float32((i/N) % w)

		box := Box{}

		box.x = (col + logisticActivate(data[index + 0])) / blockwd
		box.y = (row + logisticActivate(data[index + 1])) / blockwd
		box.w = float32(math.Exp(float64(data[index + 2]))) * anchors[2*n] / blockwd
		box.h = float32(math.Exp(float64(data[index + 3]))) * anchors[2*n+1] / blockwd

		box.confidence = logisticActivate(data[index + 4])

		if box.confidence < thresh {
			continue
		}


		box.classProbs = softmax(data[index+5 : index+5+numClasses])
		for j := 0; j < numClasses; j++ {
			box.classProbs[j] *= box.confidence
			if box.classProbs[j] < thresh {
				box.classProbs[j] = 0
			}
		}

		boxes = append(boxes, box)
	}


	// Non-Maximum-Suppression
	for k := 0; k < numClasses; k++ {
		for i := 0; i < len(boxes); i++ {
			boxes[i].currentClassIdx = k
		}

		sort.Sort(sort.Reverse(IndexSortList(boxes)))

		for i := 0; i < len(boxes); i++ {
			if boxes[i].classProbs[k] == 0 {
				continue
			}

			for j := i+1; j < len(boxes); j++ {
				if box_iou(boxes[i], boxes[j]) > nms_threshold {
					boxes[j].classProbs[k] = 0
				}
			}
		}
	}


	detectionBBoxes := []BBox{}

	for i := 0; i < len(boxes); i++ {
		max_i := max_index(boxes[i].classProbs[:])

		if max_i == -1 || boxes[i].classProbs[max_i] < thresh {
			continue
		}

		left := (boxes[i].x - boxes[i].w/2.) * img_width
		right := (boxes[i].x + boxes[i].w/2.) * img_width
		top := (boxes[i].y - boxes[i].h/2.) * img_height
		bottom := (boxes[i].y + boxes[i].h/2.) * img_height

		if left < 0 { left = 0 }
		if right > img_width { right = img_width }
		if top < 0 { top = 0 }
		if bottom > img_height { bottom = img_height }


		if left > right || top > bottom {
			continue
		}

		if int(right - left) == 0 || int(bottom - top) == 0 {
			continue
		}

		bbox := BBox{
			topleft: image.Point{ int(left), int(top)},
			bottomright: image.Point{int(right), int(bottom)},
			center: image.Point{int(boxes[i].x*img_width), int(boxes[i].y*img_height)},
			confidence: boxes[i].classProbs[max_i],
			label: max_i,
		}

		detectionBBoxes = append(detectionBBoxes, bbox)
	}

	return detectionBBoxes
}

func matToArray(m *gocv.Mat) [w*h*5*(numClasses+5)]float32 {

	result := [w*h*5*(numClasses+5)]float32{}
	i := 0
	for r := 0; r < m.Rows(); r++ {
		for c := 0; c < m.Cols(); c++ {
			result[i] = m.GetFloatAt(r, c)
			i++
		}
	}

	return result
}


func transpose(gocvMat *gocv.Mat) [w*h*5*(numClasses+5)]float32 {

	result := [w*h*5*(numClasses+5)]float32{}
	i := 0
	for c := 0; c < gocvMat.Cols(); c++ {
		for r := 0; r < gocvMat.Rows(); r++ {
			result[i] = gocvMat.GetFloatAt(r, c)
			i++
		}
	}

	return result
}


/*
 * Sorting intermediate results
 */


type IndexSortList []Box

func (i IndexSortList) Len() int {
	return len(i)
}

func (i IndexSortList) Swap(j,k int)  {
	i[j], i[k] = i[k], i[j]
}

func (i IndexSortList) Less(j,k int) bool  {
	classIdx := i[j].currentClassIdx
	return i[j].classProbs[classIdx] - i[k].classProbs[classIdx] < 0
}

func logisticActivate(x float32) float32 {
	return 1.0/(1.0 + float32(math.Exp(float64(-x))))
}


func softmax(x []float32) [numClasses]float32 {
	var sum float32 = 0.0
	var largest float32 = 0.0
	var e float32

	var output [numClasses]float32

	for i:=0; i<numClasses; i++ {
		if x[i] > largest {
			largest = x[i]
		}
	}

	for i:=0; i<numClasses; i++ {
		e = float32(math.Exp(float64(x[i] - largest)))
		sum += e
		output[i] = e
	}

	if sum > 1 {
		for i:=0; i<numClasses; i++ {
			output[i] /= sum
		}
	}

	return output
}


func overlap(x1, w1, x2, w2 float32) float32 {
	l1 := x1 - w1/2
	l2 := x2 - w2/2
	left := math.Max(float64(l1), float64(l2))

	r1 := x1 + w1/2
	r2 := x2 + w2/2
	right := math.Min(float64(r1), float64(r2))

	return float32(right - left)
}


func box_intersection(a, b Box) float32 {
	w := overlap(a.x, a.w, b.x, b.w);
	h := overlap(a.y, a.h, b.y, b.h);
	if w < 0 || h < 0 {
		return 0
	}

	area := w*h
	return area
}

func box_union(a,b Box) float32 {
	i := box_intersection(a, b)
	u := a.w*a.h + b.w*b.h - i
	return u
}


func box_iou(a, b Box) float32 {
	return box_intersection(a,b) / box_union(a,b)
}

func max_index(a []float32) int {
	if len(a) == 0 {
		return -1
	}

	max_i := 0
	max_val := math.Inf(-1)
	min_val := math.Inf(1)

	for i, val := range (a) {
		if float64(val) > max_val {
			max_i = i
			max_val = float64(val)
		}
		if float64(val) < min_val {
			min_val = float64(val)
		}
	}

	if max_val == min_val {
		return -1
	}

	return max_i
}