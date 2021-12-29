package main

import (
	"fmt"
	"math/rand"
	"time"
)

const ()

var r *rand.Rand

func main() {
	fmt.Println("starting")
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var sum int
	for i := 0; i < 30; i++ {
		sum = -1
		for sum < 0 {
			values := getValues()
			indexes := getValues()
			var line string
			for x, v := range values {
				sum = sum + calc(v, indexes[x])
				line = line + fmt.Sprintf("%d-%d, ", v, indexes[x])
			}
			if sum >= 0 {
				fmt.Printf("%d: %s\r\n", sum, line)
			}
		}
	}
}

func getValues() []int {
	values := make([]int, 3)
	for i := 0; i < 3; i++ {
		values[i] = -1
	}
	for i := 0; i < 3; i++ {
		var value int = -1
		for value < 0 {
			value = r.Intn(8)
			for x := 0; x < 3; x++ {
				if value == values[x] {
					value = -1
					break
				}
			}
		}
		values[i] = value
	}
	return values
}

func calc(i int, index int) int {
	switch index {
	case 0:
		i = i
	case 1:
		i = i + 7
	case 2:
		i = i * 3
	case 3:
		i = i - 5
	case 4:
		i = i * -1
	case 5:
		i = i * 4
	case 6:
		i = i / 2
	case 7:
		i = i * 2
	}
	return i
}
