package main

import (
	"fmt"
	"math"
)

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse

	for i := 0; i < numGens; i++ {
		var galaxy Galaxy
		galaxy = append(galaxy, timePoints[i].stars...)

		var universeTree QuadTree

		initialQuad := GetQuadrant(galaxy)
		universeRoot := MakeTree(galaxy, initialQuad)
		universeTree.root = universeRoot

		timePoints[i+1] = UpdateUniverse(timePoints[i], universeTree.root, theta, time)
		fmt.Println("timepoints:", i)
		fmt.Println("universe width is:", timePoints[i].width)
		timePoints[i].PrettyPrint()
		fmt.Println("barn worked")
	}

	return timePoints
}

//split input quadrant into equal size 4 pieces
func SplitQuadrant(rootQuad Quadrant) []Quadrant {
	quadrants := make([]Quadrant, 4)

	for i := 0; i < 4; i++ {
		quadrants[i].width = rootQuad.width / 2.0
	}
	quadrants[0].x = rootQuad.x
	quadrants[0].y = rootQuad.y + rootQuad.width/2.0

	quadrants[1].x = rootQuad.x + rootQuad.width/2.0
	quadrants[1].y = rootQuad.y + rootQuad.width/2.0

	quadrants[2].x = rootQuad.x
	quadrants[2].y = rootQuad.y

	quadrants[3].x = rootQuad.x + rootQuad.width/2.0
	quadrants[3].y = rootQuad.y

	return quadrants
}

//find the smallest quadrant with all the stars - used once at beginning
func GetQuadrant(stars Galaxy) Quadrant {
	var quad Quadrant
	var biggest OrderedPair
	quad = Quadrant{
		x: 9e100,
		y: 9e100,
	}

	for s := range stars {

		if stars[s].position.x < quad.x {
			quad.x = stars[s].position.x
		}
		if stars[s].position.y < quad.y {
			quad.y = stars[s].position.y
		}

		if stars[s].position.x > biggest.x {
			biggest.x = stars[s].position.x
		}
		if stars[s].position.y > biggest.y {
			biggest.y = stars[s].position.y
		}

	}

	xLen := biggest.x - quad.x
	yLen := biggest.y - quad.y

	if xLen < yLen {
		xLen = yLen
	}
	quad.width = xLen

	return quad
}

//read galaxy to construct tree -- output: root node of the galaxy input
func MakeTree(stars Galaxy, quad Quadrant) *Node {
	if len(stars) == 1 { // when only 1 star in the quad
		var node Node
		node.star = stars[0]
		node.sector = quad
		return &node

	}

	var NE, NW, SE, SW Galaxy
	var mass float64
	var weightedPosition OrderedPair
	var centerOfMass OrderedPair

	// calculate quad position root mass/weighted position
	for s := range stars {
		//find total mass & weighted sum of x, y of all stars to find root
		mass = mass + stars[s].mass
		weightedPosition.x = weightedPosition.x + stars[s].position.x*stars[s].mass
		weightedPosition.y = weightedPosition.y + stars[s].position.y*stars[s].mass
	}
	//calculate root position
	centerOfMass.x = weightedPosition.x / mass
	centerOfMass.y = weightedPosition.y / mass
	//assign root as dummy star
	var rootStar Star
	rootStar.position.x = centerOfMass.x
	rootStar.position.y = centerOfMass.y
	rootStar.mass = mass

	// assign star root as a dummy star with children
	var rootNode Node
	rootNode.star = &rootStar
	rootNode.sector = quad

	//find stars that belong to 4 sections of the quadrant
	for ss := range stars {
		if stars[ss].position.x < quad.x+quad.width/2.0 { //West
			if stars[ss].position.y < quad.y+quad.width/2.0 {
				SW = append(SW, stars[ss])
			} else {
				NW = append(NW, stars[ss])
			}
		}

		if stars[ss].position.x >= quad.x+quad.width/2.0 { //East
			if stars[ss].position.y < quad.y+quad.width/2.0 {
				SE = append(SE, stars[ss])
			} else {
				NE = append(NE, stars[ss])
			}
		}
	}

	// split quad into 4 quadrants
	directions := SplitQuadrant(quad)

	if len(NW) != 0 {
		rootNode.children[0] = MakeTree(NW, directions[0])
	}
	if len(NE) != 0 {
		rootNode.children[1] = MakeTree(NE, directions[1])
	}
	if len(SW) != 0 {
		rootNode.children[2] = MakeTree(SW, directions[2])
	}
	if len(SE) != 0 {
		rootNode.children[3] = MakeTree(SE, directions[3])
	}
	return &rootNode

}

func Distance(s1, s2 Star) float64 {
	dx := s1.position.x - s2.position.x
	dy := s1.position.y - s2.position.y
	return math.Sqrt(dx*dx + dy*dy)
}

func ComputeGravityForce(s1, s2 Star) OrderedPair {
	var vector OrderedPair
	d := Distance(s1, s2)
	deltaX := s2.position.x - s1.position.x
	deltaY := s2.position.y - s1.position.y
	F := G * s1.mass * s2.mass / (d * d)

	vector.x = F * deltaX / d
	vector.y = F * deltaY / d

	return vector
}

//compute the net force acting on star s
//root is the root of the entire univ tree
func ComputeNetForce(root *Node, s Star, theta float64) OrderedPair {
	var netForce OrderedPair //BUG - where do I define this var?

	if root == nil {
		return netForce
	}
	//when root is not star
	if Distance(*root.star, s) != 0 {
		// if root is leaf, no need to consider s/d conditions
		if len(root.children) == 0 || root.sector.width/Distance(*root.star, s) <= theta { // the root is a terminal node(leaf)
			netForce.x = netForce.x + ComputeGravityForce(*root.star, s).x
			netForce.y = netForce.y + ComputeGravityForce(*root.star, s).y

		} else if root.sector.width/Distance(*root.star, s) > theta { // go down tree
			NetForceChild0 := ComputeNetForce(root.children[0], s, theta)
			netForce.x = netForce.x + NetForceChild0.x
			netForce.y = netForce.y + NetForceChild0.y

			NetForceChild1 := ComputeNetForce(root.children[1], s, theta)
			netForce.x = netForce.x + NetForceChild1.x
			netForce.y = netForce.y + NetForceChild1.y

			NetForceChild2 := ComputeNetForce(root.children[2], s, theta)
			netForce.x = netForce.x + NetForceChild2.x
			netForce.y = netForce.y + NetForceChild2.y

			NetForceChild3 := ComputeNetForce(root.children[3], s, theta)
			netForce.x = netForce.x + NetForceChild3.x
			netForce.y = netForce.y + NetForceChild3.y

		}

	}

	return netForce
}

//Update Functions

func (s *Star) NewAcceleration(root *Node, theta float64) OrderedPair {
	F := ComputeNetForce(root, *s, theta)
	return OrderedPair{
		x: F.x / s.mass,
		y: F.y / s.mass,
	}
}

func (s *Star) NewVelocity(t float64) OrderedPair {
	return OrderedPair{
		x: s.velocity.x + s.acceleration.x*t,
		y: s.velocity.y + s.acceleration.y*t,
	}
}

func (s *Star) NewPosition(t float64) OrderedPair {
	return OrderedPair{
		x: s.position.x + s.velocity.x*t + 0.5*s.acceleration.x*t*t,
		y: s.position.y + s.velocity.y*t + 0.5*s.acceleration.x*t*t,
	}
}

func (s *Star) Update(root *Node, theta float64, t float64) {
	acc := s.NewAcceleration(root, theta)
	vel := s.NewVelocity(t)
	pos := s.NewPosition(t)

	s.acceleration, s.velocity, s.position = acc, vel, pos
}

func UpdateUniverse(univ *Universe, root *Node, theta, t float64) *Universe {
	newUniverse := CopyUniverse(univ)
	for b := range univ.stars {
		newUniverse.stars[b].Update(root, theta, t)
	}
	return newUniverse

}
