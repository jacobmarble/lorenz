package main

import (
	"fmt"
	"testing"

	"github.com/golang/geo/r2"
	"github.com/stretchr/testify/assert"
)

func TestScaleAndCrop(t *testing.T) {
	testCases := []struct {
		points           []r2.Point
		maxDimensionSize float64
		expectPoints     []r2.Point
		expectWidth      int
		expectHeight     int
	}{{
		points: []r2.Point{
			{X: 1, Y: 1},
			{X: 2, Y: 3},
		},
		maxDimensionSize: 10,
		expectPoints:     []r2.Point{{X: 0, Y: 0}, {X: 5, Y: 10}},
		expectWidth:      5,
		expectHeight:     10,
	}, {
		points: []r2.Point{
			{X: 1, Y: 1},
			{X: 3, Y: 2},
		},
		maxDimensionSize: 10,
		expectPoints:     []r2.Point{{X: 0, Y: 0}, {X: 10, Y: 5}},
		expectWidth:      10,
		expectHeight:     5,
	}, {
		points: []r2.Point{
			{X: 18, Y: 20},
			{X: 28, Y: 10},
		},
		maxDimensionSize: 10,
		expectPoints:     []r2.Point{{X: 0, Y: 10}, {X: 10, Y: 0}},
		expectWidth:      10,
		expectHeight:     10,
	}, {
		points: []r2.Point{
			{X: 18, Y: 20},
			{X: 28, Y: 10},
		},
		maxDimensionSize: 100,
		expectPoints:     []r2.Point{{X: 0, Y: 100}, {X: 100, Y: 0}},
		expectWidth:      100,
		expectHeight:     100,
	}}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			points, width, height := scaleAndCrop(testCase.points, testCase.maxDimensionSize)
			assert.Equal(t, testCase.expectPoints, points)
			assert.Equal(t, testCase.expectWidth, width)
			assert.Equal(t, testCase.expectHeight, height)
		})
	}
}
