package main

import (
	"fmt"
	"gifhelper"
	"log"
	"time"
)

type Particle struct {
	position         OrderedPair
	name             string
	radius           float64 // not necessary?
	diffusionRate    float64 // length of single step
	red, green, blue uint8   // color function of type?
}

type Board struct {
	width, height float64
	particles     []*Particle
}

type OrderedPair struct {
	x, y float64
}

func (b *Board) CopyBoard() *Board {
	var newBoard Board

	newBoard.width = b.width
	newBoard.height = b.height
	newBoard.particles = make([]*Particle, len(b.particles))

	for i, p := range b.particles {
		newBoard.particles[i] = p.CopyParticle()
	}

	return &newBoard
}

func (p *Particle) CopyParticle() *Particle {
	var p2 Particle

	p2 = *p // shallow copy ok because all fields are elementary

	return &p2
}

func main() {
	fmt.Println("Particle simulator.")

	fmt.Println("Generating random particles and initializing board.")

	numParticles := 1_000_000
	boardWidth := 1000.0
	boardHeight := 1000.0
	particleRadius := 5.0
	diffusionRate := 1.0
	numSteps := 100
	canvasWidth := 300
	frequency := 10

	random := false // make true if we want to scatter across board

	initialBoard := InitializeBoard(boardWidth, boardHeight, numParticles, particleRadius, diffusionRate, random)

	fmt.Println("Running simulation in serial.")

	/*
		// timing our two approaches
		fmt.Println("Running algorithm serially.")
		start := time.Now()
		boardsSerial := UpdateBoards(initialBoard, numSteps, false)
		elapsed := time.Since(start)
		log.Printf("Serial algorithm took %s", elapsed)
		images := AnimateSystem(boardsSerial, canvasWidth, frequency)
		gifhelper.ImagesToGIF(images, "boardSerial")

	*/
	fmt.Println("Running algorithm in parallel.")
	start2 := time.Now()
	boardsParallel := UpdateBoards(initialBoard, numSteps, true)
	elapsed2 := time.Since(start2)
	log.Printf("Parallel algorithm took %s", elapsed2)
	images := AnimateSystem(boardsParallel, canvasWidth, frequency)
	gifhelper.ImagesToGIF(images, "boardParallel")
}
