package util

import "fmt"

// FmtPrint outputs fmt string array.
func FmtPrint(arr [][]string) {
	if len(arr) == 0 {
		return
	}

	rows := len(arr)
	cols := len(arr[0])
	lens := make([][]int, rows)

	for i := 0; i < rows; i++ {
		lens[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			lens[i][j] = len(arr[i][j])
		}
	}

	colMax := make([]int, cols)
	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			if colMax[j] < lens[i][j] {
				colMax[j] = lens[i][j]
			}
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Print(arr[i][j])
			padding := colMax[j] - lens[i][j] + 2
			for p := 0; p < padding; p++ {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}
