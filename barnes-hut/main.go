package main

import (
	"fmt"
	"gifhelper"
	"math"
	"os"
)

// moonscale - scale up stars to make it more visible
//const moonscale = 6
func (u *Universe) PrettyPrint() {
	for _, b := range u.stars {
		b.PrettyPrint()
	}
}
func (b *Star) PrettyPrint() {
	fmt.Printf(" %v %v %v %v\n",

		b.position, b.velocity, b.acceleration, b.blue,
	)
}

func main() {

	var timePoints []*Universe

	if os.Args[1] == "jupiter" {
		jupiterUniverse := CreateJupiterSystem()
		jupiterUniverse.PrettyPrint()
		fmt.Print("\n")

		numGens := 500
		time := 2e10 //2e14 initially
		theta := 0.5

		timePoints = BarnesHut(jupiterUniverse, numGens, time, theta)
		fmt.Println("timepoint size is ", len(timePoints))
		fmt.Println("Simulation run. Now drawing images.")
		canvasWidth := 1000 //1000
		frequency := 1
		scalingFactor := 1e15 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
		imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

		fmt.Println("Images drawn. Now generating GIF.")
		gifhelper.ImagesToGIF(imageList, "jupiter")
		fmt.Println("GIF drawn.")

	} else if os.Args[1] == "galaxy" {
		g0 := InitializeGalaxy(500, 4e21, 7e22, 2e22)
		width := 1.0e23
		galaxies := []Galaxy{g0}
		initialUniverse := InitializeUniverse(galaxies, width)
		numGens := 200
		time := 2e17 //2e14
		theta := 0.5

		timePoints = BarnesHut(initialUniverse, numGens, time, theta)
		fmt.Println("Simulation run. Now drawing images.")
		canvasWidth := 500 //1000
		frequency := 1
		scalingFactor := 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
		imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

		fmt.Println("Images drawn. Now generating GIF.")
		gifhelper.ImagesToGIF(imageList, "galaxy")
		fmt.Println("GIF drawn.")

	} else if os.Args[1] == "collision" {

		g0 := InitializeGalaxy(500, 4e21, 7e22, 2e22)
		g1 := InitializeGalaxy(500, 4e21, 6e22, 4e22)

		width := 1.0e23
		galaxies := []Galaxy{g0, g1}

		initialUniverse := InitializeUniverse(galaxies, width)

		numGens := 200
		time := 3e18
		theta := 0.3

		timePoints = BarnesHut(initialUniverse, numGens, time, theta)
		fmt.Println("Simulation run. Now drawing images.")
		canvasWidth := 1000
		frequency := 1
		scalingFactor := 2e11
		imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

		fmt.Println("Images drawn. Now generating GIF.")
		gifhelper.ImagesToGIF(imageList, "collision")
		fmt.Println("GIF drawn.")

	}

}

func CreateJupiterSystem() *Universe {
	var jupiter, io, europa, ganymede, callisto Star

	jupiter.red, jupiter.green, jupiter.blue = 223, 227, 202
	io.red, io.green, io.blue = 249, 249, 165
	europa.red, europa.green, europa.blue = 132, 83, 52
	ganymede.red, ganymede.green, ganymede.blue = 76, 0, 153
	callisto.red, callisto.green, callisto.blue = 0, 153, 76

	jupiter.mass = 1.898 * math.Pow(10, 27)
	io.mass = 8.9319 * math.Pow(10, 22)
	europa.mass = 4.7998 * math.Pow(10, 22)
	ganymede.mass = 1.4819 * math.Pow(10, 23)
	callisto.mass = 1.0759 * math.Pow(10, 23)

	// multiplied radius by 10 ? Fix?
	jupiter.radius = 71000000
	io.radius = 18210000
	europa.radius = 15690000
	ganymede.radius = 26310000
	callisto.radius = 24100000

	jupiter.position.x, jupiter.position.y = 2000000000, 2000000000
	io.position.x, io.position.y = 2000000000-421600000, 20000000000
	europa.position.x, europa.position.y = 2000000000, 2000000000+670900000
	ganymede.position.x, ganymede.position.y = 2000000000+1070400000, 2000000000
	callisto.position.x, callisto.position.y = 2000000000, 2000000000-1882700000

	jupiter.velocity.x, jupiter.velocity.y = 0, 0
	io.velocity.x, io.velocity.y = 0, -17320
	europa.velocity.x, europa.velocity.y = -13740, 0
	ganymede.velocity.x, ganymede.velocity.y = 0, 10870
	callisto.velocity.x, callisto.velocity.y = 8200, 0

	var jupiterUniverse Universe
	jupiterUniverse.width = 1.0e23
	jupiterUniverse.AddStar(jupiter)
	jupiterUniverse.AddStar(io)
	jupiterUniverse.AddStar(europa)
	jupiterUniverse.AddStar(ganymede)
	jupiterUniverse.AddStar(callisto)

	return &jupiterUniverse
}

func (u *Universe) AddStar(s Star) {
	u.stars = append(u.stars, &s)
}
