package main

import (
	"image"
	"image/png"
	"math"
	"os"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
)

const (
	size       = 800
	rho        = 28
	sigma      = 10
	beta       = 8 / 3
	xInit      = 0
	yInit      = -4
	zInit      = 23
	precision  = 1000
	iterations = 2000000
)

func main() {
	vectors := lorenz(rho, sigma, beta, r3.Vector{X: xInit, Y: yInit, Z: zInit}, precision, iterations)
	points := threeToTwo(vectors, func(vector r3.Vector) r2.Point {
		return r2.Point{X: vector.X, Y: -vector.Z}
	})
	points, width, height := scaleAndCrop(points, size)

	img := image.NewGray(image.Rect(0, 0, width, height))
	for _, point := range points {
		x, y := int(point.X), int(point.Y)
		color := img.GrayAt(x, y)
		if color.Y + 8 < color.Y {
			color.Y = math.MaxUint8
		} else {
			color.Y += 8
		}
		img.SetGray(x, y, color)
	}

	f, err := os.Create("lorenz.png")
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func threeToTwo(vectors []r3.Vector, f func(vector r3.Vector) r2.Point) []r2.Point {
	points := make([]r2.Point, len(vectors))
	for i := range vectors {
		points[i] = f(vectors[i])
	}
	return points
}

func scaleAndCrop(points []r2.Point, maxDimensionSize float64) ([]r2.Point, int, int) {
	minX, minY, maxX, maxY := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for _, point := range points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	var width, height, scale float64
	if maxX-minX > maxY-minY {
		width = maxDimensionSize
		scale = width / (maxX - minX)
		height = (maxY - minY) * scale
	} else {
		height = maxDimensionSize
		scale = height / (maxY - minY)
		width = (maxX - minX) * scale
	}
	shiftX, shiftY := -minX, -minY

	result := make([]r2.Point, len(points))
	for i := range points {
		result[i].X = (points[i].X + shiftX) * scale
		result[i].Y = (points[i].Y + shiftY) * scale
	}

	return result, int(math.Ceil(width)), int(math.Ceil(height))
}

func lorenz(rho, sigma, beta float64, initialState r3.Vector, precision float64, iterations int) []r3.Vector {
	vectors := make([]r3.Vector, 0, iterations+1)

	vectors = append(vectors, initialState)
	state := initialState
	for i := 0; i < iterations; i++ {
		state = nextLorenz(precision, rho, sigma, beta, state)
		vectors = append(vectors, state)
	}

	return vectors
}

func nextLorenz(precision, rho, sigma, beta float64, state r3.Vector) r3.Vector {
	dx := sigma * (state.Y - state.X)
	dy := state.X*(rho-state.Z) - state.Y
	dz := state.X*state.Y - beta*state.Z
	return r3.Vector{
		X: state.X + dx/precision,
		Y: state.Y + dy/precision,
		Z: state.Z + dz/precision,
	}
}
