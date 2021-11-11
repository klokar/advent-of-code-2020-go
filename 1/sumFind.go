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
	found    bool
	foundNum int
	wg       sync.WaitGroup
}

func main() {
	numbers, error := readFile(fileName)
	if error != nil {
		fmt.Println("File could not be opened! -> ", error)
		return
	}

	stats := stats{numbers: numbers, found: false}
	for number := numbers.Front(); number != nil && stats.found == false; number = number.Next() {
		stats.current = number.Value.(int)
		stats.compared = number.Next()

		for stats.compared != nil {
			stats.runWorkers()
		}
	}

	fmt.Println(fmt.Sprintf("Answer: Looking for: %d & %d. Multiplies to: %d", stats.current, stats.foundNum, stats.current*stats.foundNum))
}

func (s *stats) runWorkers() {
	for i := 0; i < noOfWorkers; i++ {

		// When notifying WG about new goroutine be sure to only increment by 1
		s.wg.Add(1)

		i := i
		go func() {
			defer s.wg.Done()
			checking := s.compared
			if checking != nil {
				s.compared = s.compared.Next()
				found, number := check(i, s.current, checking.Value.(int))

				if found {
					s.found = true
					s.foundNum = number
				}
			}
		}()
	}

	s.wg.Wait()
}

func check(workerId, num1, num2 int) (bool, int) {
	fmt.Println(fmt.Sprintf("Worker:%d, Comparing: %d & %d", workerId, num1, num2))
	return num1+num2 == required, num2
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
