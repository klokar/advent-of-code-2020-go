package main

import (
	"bufio"
	linkedList "container/list"
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	fileName    = "./1/inputs.txt"
	required    = 2020
	noOfWorkers = 12
)

type stats struct {
	numbers  *linkedList.List
	current  int
	compared *linkedList.Element
	found    []int
}

func main() {
	numbers, error := readFile(fileName)
	if error != nil {
		fmt.Println("File could not be opened! -> ", error)
		return
	}

	stats := stats{numbers: numbers}
	for number := numbers.Front(); number != nil && len(stats.found) == 0; number = number.Next() {
		stats.current = number.Value.(int)
		stats.compared = numbers.Front()

		for stats.compared != nil {
			stats.runWorkers()
		}
	}

	if len(stats.found) != 0 {
		fmt.Println(fmt.Sprintf(
			"Answer: Looking for: %d & %d & %d. Multiplies to: %d",
			stats.found[0], stats.found[1], stats.found[2], stats.found[0]*stats.found[1]*stats.found[2],
		))
	}
}

func (s *stats) runWorkers() {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {

		// When notifying WG about new goroutine be sure to only increment by 1
		wg.Add(1)

		i := i
		go func() {
			defer wg.Done()
			checking := s.compared
			if checking != nil {
				s.compared = s.compared.Next()

				for number3 := s.numbers.Front(); number3 != nil && len(s.found) == 0; number3 = number3.Next() {
					found, number1, number2, number3 := check(i, s.current, checking.Value.(int), number3.Value.(int))

					if found {
						s.found = []int{number1, number2, number3}
					}
				}
			}
		}()
	}

	wg.Wait()
}

func check(workerId, num1, num2, num3 int) (bool, int, int, int) {
	//fmt.Println(fmt.Sprintf("Worker:%d, Comparing: %d & %d & %d", workerId, num1, num2, num3))
	return num1+num2+num3 == required, num1, num2, num3
}

func readFile(path string) (*linkedList.List, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := linkedList.New()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		numeric, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Error parsing number - " + scanner.Text())
			return nil, err
		}

		lines.PushBack(numeric)
	}

	return lines, scanner.Err()
}
