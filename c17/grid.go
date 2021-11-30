package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	path           = "./c17/start.txt"
	activeSymbol   = "#"
	inactiveSymbol = "."
	bootCycles     = 6
)

func main() {
	currentStates, width, err := loadGrid(path)
	if err != nil {
		fmt.Println("Could not read file!")
	}

	displayAndCount(*currentStates, 0, width)

	for iteration := 1; iteration <= bootCycles; iteration++ {
		previewStates := make(map[string]bool)
		for x := 0 - iteration; x < width+iteration; x++ {
			for y := 0 - iteration; y < width+iteration; y++ {
				for z := 0 - iteration; z <= iteration; z++ {
					for w := 0 - iteration; w <= iteration; w++ {
						currentState := (*currentStates)[key(x, y, z, w)]
						newState := calculateState(currentState, activeNeighbors(currentStates, x, y, z, w))
						previewStates[key(x, y, z, w)] = newState
					}
				}
			}
		}

		currentStates = &previewStates
	}

	displayAndCount(*currentStates, bootCycles, width)
}

func displayAndCount(states map[string]bool, iterations, width int) {
	fmt.Println("----------")
	count := 0
	for w := 0 - iterations; w <= 0+iterations; w++ {
		for z := 0 - iterations; z <= 0+iterations; z++ {
			fmt.Println(fmt.Sprintf("z=%d, w=%d", z, w))
			for y := 0 - iterations; y < width+iterations; y++ {
				for x := 0 - iterations; x < width+iterations; x++ {
					active := states[key(x, y, z, w)]
					if active {
						count++
					}
					fmt.Print(symbol(active))
				}
				fmt.Print("\n")
			}
		}
	}
	fmt.Println("Active count: ", count)
}

func calculateState(active bool, neighbours int) bool {
	if active && (neighbours < 2 || neighbours > 3) {
		active = false
	} else if !active && neighbours == 3 {
		active = true
	}

	return active
}

func activeNeighbors(currentState *map[string]bool, xC, yC, zC, wC int) int {
	count := 0
	for x := xC - 1; x <= xC+1; x++ {
		for y := yC - 1; y <= yC+1; y++ {
			for z := zC - 1; z <= zC+1; z++ {
				for w := wC - 1; w <= wC+1; w++ {
					if (*currentState)[key(x, y, z, w)] && key(x, y, z, w) != key(xC, yC, zC, wC) {
						count++
					}
				}
			}
		}
	}

	return count
}

func loadGrid(path string) (*map[string]bool, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	cubes := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	width := 0
	y := 0
	z := 0
	w := 0
	for scanner.Scan() {
		width = len(scanner.Text())
		for x, rune := range scanner.Text() {
			cubes[key(x, y, z, w)] = isActive(string(rune))
		}
		y++
	}

	return &cubes, width, scanner.Err()
}

func key(x, y, z, w int) string {
	return fmt.Sprintf("%d,%d,%d,%d", x, y, z, w)
}

func isActive(char string) bool {
	return char == activeSymbol
}

func symbol(active bool) string {
	if active {
		return activeSymbol
	}

	return inactiveSymbol
}
