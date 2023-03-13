package main

//  0 1 2 3 4 5
// [* * * * * * * * * * * * * * * * * * *  ]
//  |----0---||----1---||---2----||---3----|
//  ^         ^
//  |         |
//  start     end

/*
105 items
10 processors
chunksize = 10
chunk0 = 0..10
chunk1 = 10..20
chunk2 = 20..30
...
chunk9 = 90..100

*/
will work for "any" f:

	[a[f(i):f(i+1)] for i in range(n)]

Take f(i) = i*k. This is the equivelent to what we wrote originally in go:

	[a[i*k:(i+1)*k] for i in range(n)]

Consider:
	f(i) = i*k+c for some constant c:
	[a[i*k+5:(i+1)*k +5]]

What will happen?

Need: 
	(A) f(0) = 0 (so we start at the start of the array)
	(B) f(n+1) = length of the array (so we end at the end)

Consider: f(i) = i*k+i

	[a[i*k+i:(i+1)*k+i+1] for i in range(n)]

Has property (A), and increases the size by 1 (due to i+1) but not (B). Take:

	f(i) = i*k+min(i,m)

	[a[i*k+min(i, m):(i+1)*k+min(i+1, m)] for i in range(n)]

For the first m bins, add 1; for the remaning bins, add nothing.

def split(a, n):
    k, m = divmod(len(a), n)
    return (a[i*k+min(i, m):(i+1)*k+min(i+1, m)] for i in range(n))

// fill in parallel code here
func (b *Board) DiffuseParallel(numProcs int) {
	chunkSize := len(b.particles) / numProcs
	finished := make(chan bool)
	for i := 0; i <= numProcs; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end >= len(b.particles) {
			end = len(b.particles)
		}

		go DiffuseOneProc(b.particles[start:end], finished)
	}
	for i := 0; i <= numProcs; i++ {
		<-finished
	}
}

func DiffuseOneProc(particles []*Particle, done chan bool) {
	//fmt.Println("Processing chunk", i, "in generation", gen)
	for _, p := range particles {
		p.RandStep()
	}
	done <- true
}

func (p *Particle) RandStepFast() {
}
