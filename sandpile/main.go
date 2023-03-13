package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

type SandBoard [][]int

func main() {
	rand.Seed(time.Now().UnixNano())
	numProcs := runtime.NumCPU()

	// board is size x size
	sizeString := os.Args[1]
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		fmt.Println("error converting size")
	}
	initialBoard := InitializeBoard(size, size)
	fmt.Println("initial board of size ", len(initialBoard), " created")

	//pile = number of coins
	pileString := os.Args[2]
	pile, err := strconv.Atoi(pileString)
	if err != nil {
		fmt.Println("error converting pile")
	}
	fmt.Println("number of coins is ", pile)

	placement := os.Args[3]
	if placement == "central" {
		initialBoard[size/2][size/2] = pile
	} else if placement == "random" {
		RandomPlacement(initialBoard, pile, size)
	}
	//make copy of initialBoard for parallel
	initialBoardParallel := CopyBoard(initialBoard)
	fmt.Println("placement of coins finished")

	serialStart := time.Now()
	finalBoardSerial := SerialTopple(initialBoard)
	serialElapsed := time.Since(serialStart)
	fmt.Println("Serial took ", serialElapsed.Seconds(), " seconds")

	parallelStart := time.Now()
	finalBoardParallel := ParallelTopple(initialBoardParallel, numProcs)
	parallelElapsed := time.Since(parallelStart)
	fmt.Println("Parallel took ", parallelElapsed.Seconds(), " seconds")

	fmt.Println("topple finished. now drawing gameboards")
	cellWidth := 10
	img := DrawGameBoard(finalBoardSerial, cellWidth)

	//create png file
	f, err := os.Create("serial.png")
	if err != nil {
		fmt.Println("error: cannot create png file")
	}
	//write image to png file
	err = png.Encode(f, img)
	if err != nil {
		fmt.Println("error:cannot write board image(serial)")
	}

	//create image object of parallel final sandpile
	cellWidth2 := 10
	img2 := DrawGameBoard(finalBoardParallel, cellWidth2)
	//create png file
	ff, err := os.Create("parallel.png")
	if err != nil {
		fmt.Println("error: cannot create png file")
	}
	//write image to png file
	errP := png.Encode(ff, img2)
	if errP != nil {
		fmt.Println("error:cannot write board image (parallel)")
	}

}

func InitializeBoard(rowsize, colsize int) SandBoard {
	board := make(SandBoard, rowsize)
	for r := range board {
		board[r] = make([]int, colsize)
	}
	return board
}

func CopyBoard(board SandBoard) SandBoard {
	size := len(board)
	newBoard := make(SandBoard, size)
	for r := range newBoard {
		newBoard[r] = board[r]
	}
	return newBoard
}

func RandomPlacement(board SandBoard, pile int, size int) {

	//choose 100 random [r][c]
	var coin, row, col []int
	for i := 0; i < 100; i++ {
		row = append(row, rand.Intn(size))
		col = append(col, rand.Intn(size))
	}

	// genereate a list size 100 with random numbers that sum to pile
	totalCoin := pile
	for j := 0; j < 99; j++ {
		seats := 100 - j
		maxCoin := totalCoin + 1 - seats
		a := rand.Intn(maxCoin) + 1
		coin = append(coin, a)
		totalCoin = totalCoin - a
	}
	coin = append(coin, totalCoin)

	//place the numbers into sanboard
	for ii := 0; ii < 100; ii++ {
		board[row[ii]][col[ii]] = board[row[ii]][col[ii]] + coin[ii]
	}

}

func (board SandBoard) ToppleOnce() {
	rowSize := len(board)
	colSize := len(board[0])
	newBoard := InitializeBoard(rowSize, colSize)
	for r := 0; r < rowSize; r++ {
		for c := 0; c < colSize; c++ {
			if board[r][c] != 0 {
				newBoard[r][c] = newBoard[r][c] + board[r][c]%4
			}

			q := board[r][c] / 4

			if q != 0 {
				if r-1 >= 0 {
					newBoard[r-1][c] = newBoard[r-1][c] + q
				}
				if r+1 < rowSize {
					newBoard[r+1][c] = newBoard[r+1][c] + q
				}
				if c-1 >= 0 {
					newBoard[r][c-1] = newBoard[r][c-1] + q
				}
				if c+1 < colSize {
					newBoard[r][c+1] = newBoard[r][c+1] + q
				}
			}

		}
		board[r] = newBoard[r]
	}

}

func SerialTopple(board SandBoard) SandBoard {

	size := len(board)
	for {
		finished := true
		board.ToppleOnce()

		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				if board[r][c] >= 4 {
					finished = false
				}
			}
			if !finished {
				break
			}
		}

		if finished {
			fmt.Println("escaped!")
			break
		}
	}
	fmt.Println("the correct one is")
	return board
}

func (board SandBoard) MultiToppleOnce(start, end int, channel chan []int, done chan bool) {
	numRows := len(board)
	numCols := len(board[0])
	newBoard := InitializeBoard(numRows, numCols)

	for r := start; r < end; r++ {
		for c := 0; c < numCols; c++ {
			if board[r][c] != 0 {
				newBoard[r][c] = newBoard[r][c] + board[r][c]%4
			}

			q := board[r][c] / 4
			if q != 0 {
				if r-1 >= start && r-1 >= 0 {
					newBoard[r-1][c] = newBoard[r-1][c] + q
				} else if r-1 >= 0 && r-1 < start {
					temp := make([]int, 3)
					temp[0] = r - 1
					temp[1] = c
					temp[2] = q
					channel <- temp
				}

				if r+1 < end && r+1 < numRows {
					newBoard[r+1][c] = newBoard[r+1][c] + q
				} else if r+1 >= end && r+1 < numRows {
					temp := make([]int, 3)
					temp[0] = r + 1
					temp[1] = c
					temp[2] = q
					channel <- temp
				}

				if c-1 >= 0 {
					newBoard[r][c-1] = newBoard[r][c-1] + q
				}
				if c+1 < numCols {
					newBoard[r][c+1] = newBoard[r][c+1] + q
				}
			}

		}
		board[r] = newBoard[r]
	}
	done <- true
}

func ParallelTopple(board SandBoard, numProcs int) SandBoard {
	numRows := len(board)
	numCols := len(board[0])
	chunkSize := numRows / numProcs

	channel := make(chan []int, numCols*(numProcs+1))
	done := make(chan bool, numProcs)

	for {
		finished := true
		for i := 0; i <= numProcs; i++ {
			start := i * chunkSize
			end := (i + 1) * chunkSize
			if end >= numRows {
				end = numRows
			}
			go board.MultiToppleOnce(start, end, channel, done)

		}

		for i := 0; i <= numProcs; i++ {
			<-done
		}

		len1 := len(channel)
		for i := 0; i < len1; i++ {
			spill := <-channel
			board[spill[0]][spill[1]] = board[spill[0]][spill[1]] + spill[2]
		}

		for r := 0; r < numRows; r++ {
			for c := 0; c < numCols; c++ {
				if board[r][c] >= 4 {
					finished = false
				}
			}
			if !finished {
				break
			}
		}

		if finished {
			fmt.Println("escaped!")
			break
		}
	}

	return board
}

func (board SandBoard) Spill(spill [][]int, start, end int, over chan bool) {
	for i := start; i < end; i++ {
		if spill[i][2] != 0 {
			board[spill[i][0]][spill[i][1]] = board[spill[i][0]][spill[i][1]] + spill[i][2]
		}

	}
	over <- true
}
