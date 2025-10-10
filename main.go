package main

import (
	"fmt"
	"gopraktika/simptab"
)

//var matrica [][]float64

func main() {
	c_vector := []float64{7, 8, 3}
	b_vector := []float64{4, 7, 8}

	a_matrix := [][]float64{{3,1,1},{1,4,0},{0,0.5,2}}


	table := simptab.NewTable(c_vector, b_vector, a_matrix, true)
	table.Print(-1,-1)
	table.MakeKanonView()
	fmt.Println()

	table.DoSimplexMethod()
}
