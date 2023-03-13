package main

import (
	"math"
	"math/rand"
)

func (b *Board) DiffuseSerial() {
	for _, p := range b.particles {
		p.RandStep()
	}
}

func (p *Particle) RandStep() {
	stepLength := p.diffusionRate
	angle := rand.Float64() * 2 * math.Pi
	p.position.x += stepLength * math.Cos(angle)
	p.position.y += stepLength * math.Sin(angle)
}
