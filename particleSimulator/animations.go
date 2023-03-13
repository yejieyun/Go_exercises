package main

import (
	"canvas"
	"image"
)

func AnimateSystem(boards []*Board, canvasWidth int, frequency int) []image.Image {
	images := make([]image.Image, 0, len(boards))
	for i, b := range boards {
		if i%frequency == 0 {
			images = append(images, b.DrawToCanvas(canvasWidth))
		}
	}
	return images
}

func (b *Board) DrawToCanvas(canvasWidth int) image.Image {
	aspectRatio := b.height / b.width
	canvasHeight := int(float64(canvasWidth) * aspectRatio)
	c := canvas.CreateNewCanvas(canvasWidth, canvasHeight)

	// first, make a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	for _, p := range b.particles {
		// make a circle at p's position with the appropriate width
		scalingFactor := float64(canvasHeight) / b.height

		c.SetFillColor(canvas.MakeColor(p.red, p.green, p.blue))

		c.Circle(p.position.x*scalingFactor, p.position.y*scalingFactor, p.radius*scalingFactor)

		c.Fill()
	}
	return canvas.GetImage(c)
}
