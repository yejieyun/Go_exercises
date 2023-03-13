package main

import (
	"math/rand"
	"runtime"
)

func UpdateBoards(initialBoard *Board, numSteps int, isParallel bool) []*Board {
	boards := make([]*Board, numSteps+1)
	boards[0] = initialBoard

	for i := 1; i <= numSteps; i++ {
		boards[i] = boards[i-1].UpdateBoard(isParallel)
	}

	return boards
}

func (b *Board) UpdateBoard(isParallel bool) *Board {
	newBoard := b.CopyBoard()

	if isParallel {
		numProcs := runtime.NumCPU()
		newBoard.DiffuseParallel(numProcs)
	} else {
		newBoard.DiffuseSerial()
	}

	return newBoard
}

func InitializeBoard(boardWidth, boardHeight float64, numParticles int, particleRadius float64, diffusionRate float64, random bool) *Board {
	var b Board

	b.width = boardWidth
	b.height = boardHeight

	b.particles = make([]*Particle, numParticles)

	for i := range b.particles {
		var p Particle
		if random {
			p.position.x = rand.Float64() * boardWidth
			p.position.y = rand.Float64() * boardHeight
		} else {
			// default: non-random: assign all to center of board
			p.position.x = boardWidth / 2
			p.position.y = boardHeight / 2
		}
		p.radius = particleRadius
		p.diffusionRate = diffusionRate
		p.red, p.green, p.blue = 255, 255, 255
		b.particles[i] = &p
	}

	return &b
}
