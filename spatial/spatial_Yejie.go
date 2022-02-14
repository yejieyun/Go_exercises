package main

import ( 
	"bufio"
	"fmt"
	"os"
	"strconv"
)


// The data stored in a single cell of a field
type Cell struct {
	strategy  string //represents "C" or "D" corresponding to the type of prisoner in the cell
	score float64 //represents the score of the cell based on the prisoner's relationship with neighboring cells
}

// The game board is a 2D slice of Cell objects 
type GameBoard [][]Cell

// main function 
func main() {
	fmt.Println("spatial games")
	//initial CD assignment 
	initialAssignmentFile := os.Args[1] 
	//points rewarded (float)
	bPointsRewarded, err := strconv.ParseFloat(os.Args[2],64)
	if err != nil {
		panic("Error: Problem converting b (points) parameter to a float.")
	}
	// steps = generations (int)
	steps, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("Error: Problem converting steps parameter to an integer.")
	}
	fmt.Println("Parameters read in successfully!")

	// read initial assignment from file and save on gameboard
	initialBoard := ReadBoardFromFile(initialAssignmentFile)
	fmt.Println("initial board created from file")

	//play spatial games 
	// boards is a collection of boards ** DIFFERENT FROM BOARD
	boards := PlaySpatialGames(initialBoard, bPointsRewarded, steps)
	//save separate for final board image 
	finalBoard := FinalSpatialGameResult(initialBoard, bPointsRewarded, steps)
	fmt.Println("spatial games completed")

	// create slice of board images - Use cellular automata given functions 
	imglist := DrawGameBoards(boards, 5) 
	fmt.Println("game board list of images created")
	// create final board image CHECK IF THIS WORKS 
	finalImg := DrawFinalGameBoard(finalBoard, 5) 

	//create a gif using the images in imglist 
	gifhelper.ImagesToGIF(imglist, outputFile)
	fmt.Println("Yay! gif created")

}



//read initial assignment files 
//ReadBoardFromFile takes a filename as a string and reads in the data provided
//in this file, returning a game board.
func ReadBoardFromFile(filename string) GameBoard {
	
	//open file 
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//read through file to make board 
	scanner := bufio.NewScanner(file)

	// ready only the first line for dimensions
	var line1 int
	for scanner.Scan() {
		if line1 == 0 {
		string dims := scanner.Text()
		}

	}

	// split and convert dims (dimensions) into 1 integer as the initial boards are a square
	var x []string = strings.Split(dims)
	fmt.Println("The dimensions of initial board is", x)

	x, err = strconv.Atoi(x)
	if err != nil {
		panic("Error: Problem converting steps parameter to an integer.")
	}
	
	// the row/column number of the initial board is x/y
	x := x[0]
	y := x[1]

// make gameboard (2D slice of Cell objects)
board := make(GameBoard, x)
for i := 0; i <= x; i++ {
	board[i] = make([]Cell,y)
}
	i := 1 // skip row 0 
	for scanner.Scan() {
		currentLine := scanner.Text()
		// currentArray := make([]string, 0) // Cell.strategy-- array of strings
		var strategy string
		for j := range currentLine {
			letter = currentLine[i] // 
			board[i][j].strategy = letter
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return board
}

//update board to series of boards (cosider each board through each step// how many points goes where can be considered in UpdateBoards)
func PlaySpatialGames(initialBoard, steps int) {
	
	boards := make([]Gameboard, steps+1)
	boards[0] = initialBoard

	for i := 1; i<=steps; i++ {
		boards[i] = UpdateBoard(boards[i-1], bPointsRewarded)
	}

	return boards
}

func UpdateBoard(currBoard Gameboard, bPointsRewarded float64) GameBoard {
	//create new empty gameboard
	numRows := CountRows(currBoard)
	numCols := CountCols(currBoard)
	newBoard :=InitializeBoard(numRows-1, numCols) // first row shows dimension
	 
	//add new score values to each cell
	for r := 1; r< numRows-2; r++ {
		//repeat for each row 
		for c := 1; c < numCols-1; c++ {
			//repeat for each column 
			// makes currBoard[r][c]
			newBoard[r][c] = UpdateScore(currBoard, r, c, bPointsRewarded) //updates score and saves it in newBoard
			newBoard[r][c] = UpdateStrategy(newBoard, r, c) // updates strategy of newBoard
		}
	}

	return newBoard
	// new updated board 
}

func UpdateScore(currBoard Gameboard, r int, c int, bPointsRewarded float64) Gameboard {
	allRows := 0
	allCols := 0
	for allRows = 1; allRows <= r; allRows++{
		for all Cols = 0; allCols <= c; allCols++ {
			for i = -1; i <= 1; i ++{
				for j = -1; j<= 1; j++{
				//loops for r-1,r,r+1
					if currBoard[r][c].strategy == "C" {
						if currBoard[r][c].strategy == currBoard[r+i][c+j].strategy {
							currBoard[r][c].score = currBoard[r][c].score + 1
						}
		
					}
					else {
						if currBoard[r][c].strategy != currBoard[r+i][c+j].strategy {
							currBoard[r][c].score = currBoard[r][c].score + bPointsRewarded
						}
					}
						
					}
				}
			}
		}
	}


	
	return currBoard
}

func UpdateStrategy(currBoard Gameboard, r int, c int) Gameboard {
	allRows := 0
	allCols := 0
	for allRows = 1; allRows <= r; allRows++n {
		for all Cols = 0; allCols <= c; allCols++{
			for i = -1; i <= 1; i ++{
				for j = -1; j<= 1; j++{

				}
			}

		}
	}
	
}



func InitializeBoard(numRows, numCols int) Gameboard {
	// make a 2-D slice (default values = false)
	var board GameBoard
	board = make(GameBoard, numRows)
	// now we need to make the rows too
	for r := range board {
		board[r] = make([]int, numCols)
	}

	return board
}

func DrawGameBoards(boards []GameBoard, cellWidth int) []image.Image {
	numGenerations := len(boards)
	imageList := make([]image.Image, numGenerations)
	for i := range boards {
		imageList[i] = DrawGameBoard(boards[i], cellWidth)
	}
	return imageList
}

func DrawFinalGameBoard(finalBoard GameBoard, cellWidth int) image.Image {
	image := make image.Image
	image = DrawGameBoard(finalBoard, cellWidth)

	return image
}

func DrawGameBoard(board GameBoard, cellWidth int) image.Image {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewPalettedCanvas(width, height, nil)

	// declare colors
	darkGray := MakeColor(50, 50, 50)
	// black := MakeColor(0, 0, 0)
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)
	green := MakeColor(0, 255, 0)
	yellow := MakeColor(255, 255, 0)
	magenta := MakeColor(255, 0, 255)
	white := MakeColor(255, 255, 255)
	cyan := MakeColor(0, 255, 255)

	/*
		//set the entire board as black
		c.SetFillColor(gray)
		c.ClearRect(0, 0, height, width)
		c.Clear()
	*/

	/*
		// draw the grid lines in white
		c.SetStrokeColor(white)
		DrawGridLines(c, cellWidth)
	*/

	// fill in colored squares
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 0 {
				c.SetFillColor(darkGray)
			} else if board[i][j] == 1 {
				c.SetFillColor(blue)
			} else if board[i][j] == 2 {
				c.SetFillColor(red)
			} else if board[i][j] == 3 {
				c.SetFillColor(green)
			} else if board[i][j] == 4 {
				c.SetFillColor(yellow)
			} else if board[i][j] == 5 {
				c.SetFillColor(magenta)
			} else if board[i][j] == 6 {
				c.SetFillColor(white)
			} else if board[i][j] == 7 {
				c.SetFillColor(cyan)
			} else {
				panic("Error: Out of range value " + string(board[i][j]) + " in board when drawing board.")
			}
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}

	return GetImage(c)
}
