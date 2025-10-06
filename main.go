package main

import (
	"fmt"
	"gopraktika/simptab"
)

//var matrica [][]float64

func main() {
	//bCount - число элементов вектора ограничений В
	//xCount - число переменных, оно же число  элементов вектора целевой функции

	//var bCount, xCount int
	/*
	fmt.Println("Введите число ограничений B и число переменных X (через пробел): ")
	fmt.Scan(&bCount, &xCount)

	c_vector := make([]float64, xCount)
	fmt.Println("Введите коэффициенты целевой функции (через пробел): ")
	for i := 0; i < xCount; i++ {
		fmt.Scan(&c_vector[i])
	}

	a_matrix := make([][]float64, bCount)
	fmt.Println("Введите матрицу ограничений (B строк по X чисел): ")
	for i := 0; i < bCount; i++ {
		a_matrix[i] = make([]float64, xCount)
		for j := 0; j < xCount; j++ {
			fmt.Scan(&a_matrix[i][j])
		}
	}

	b_vector := make([]float64, bCount)
	fmt.Println("Введите вектор ограничений (B чисел):")
	for i := 0; i < bCount; i++ {
		fmt.Scan(&b_vector[i])
	}
	//Тут мы как бы должны добавить базисные переменные и выразить их через свободные но в силу ограничения <= в условии они все с + и все базовые
	*/
	//bCount, xCount = 3, 3
	c_vector := []float64{5, 6, 4}
	b_vector := []float64{7, 8, 6}

	a_matrix := [][]float64{{1,1,1},{1,3,0},{0,0.5,4}}


	table := simptab.NewTable(c_vector, b_vector, a_matrix, true)
	table.Print(-1,-1)
	table.MakeKanonView()
	fmt.Println()

	table.DoSimplexMethod()
}
