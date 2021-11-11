package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const path = "./3/landscape.txt"

type landscape struct {
	width  int
	height int
	trees  [][]bool
}

type toboggan struct {
	x        int
	y        int
	treesHit int
}

func main() {
	landscape, error := parseLandscape(path)
	if error != nil {
		fmt.Println("File could not be opened! -> ", error)
		return
	}

	paths := [][]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}}

	allTreesHit := 1
	for _, movingRange := range paths {
		allTreesHit *= landscape.travel(movingRange[0], movingRange[1])
	}

	fmt.Println(allTreesHit)
}

func (l landscape) travel(moveX, moveY int) int {
	toboggan := toboggan{x: 0, y: 0}
	for toboggan.y < l.height {
		toboggan.move(l, moveX, moveY)
	}

	return toboggan.treesHit
}

func (t *toboggan) move(l landscape, moveX, moveY int) {
	// Complete if at the bottom of the map
	if t.y+moveY > l.height {
		fmt.Println("End reached")
		return
	}

	t.y += moveY
	t.x += moveX

	// Reset X if off the right edge of the map
	if t.x > l.width-1 {
		//fmt.Println("Resetting X:", t.x, t.x-l.width)
		t.x = t.x - l.width
	}

	//fmt.Println("Moved to:", t.y, t.x)

	if l.trees[t.y][t.x] {
		t.treesHit++
		//fmt.Println("Tree HIT")
	}
}

func parseLandscape(path string) (*landscape, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileLinesCount, err := lineCounter(file)
	if err != nil {
		return nil, err
	}

	landscape := landscape{height: fileLinesCount}
	trees := make([][]bool, fileLinesCount+1)

	// Reset file pointer to read from beginning
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)

	line := 0
	width := 0
	for scanner.Scan() {
		width = len(scanner.Text())
		trees[line] = make([]bool, width)
		for i, char := range scanner.Text() {
			trees[line][i] = string(char) == "#"
		}

		line += 1
	}

	landscape.width = width
	landscape.trees = trees

	return &landscape, scanner.Err()
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
