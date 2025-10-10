package simptab

import (
	"fmt"
	"math"
	"strconv"
)

type SimplexTable struct {
	Basis   []string    // базисные переменные
	Headers []string    // свободные переменные
	Table   [][]float64 // сама таблица
	View    bool
}

func (simptab SimplexTable) Print(rowPrint int, colPrint int) { // -1 без подстветки

	red := "\033[31m"
	reset := "\033[0m"

	fmt.Printf("%-4s", "B\\F") // вывод свободных переменных
	for i, val := range simptab.Headers {
		if i == colPrint {
			fmt.Printf("%s%8s%s", red, val, reset)
		} else {
			fmt.Printf("%8s", val)
		}

	}
	fmt.Println()

	for i, row := range simptab.Table {
		if i == rowPrint {
			fmt.Printf("%s%-4s%s", red, simptab.Basis[i], reset) //вывод  базиса
			for _, val := range row {
				fmt.Printf("%s%8.2f%s", red, val, reset)
			}
			fmt.Println()
		} else {
			fmt.Printf("%-4s", simptab.Basis[i]) // вывод базиса
			for j, val := range row {
				if j == colPrint {
					fmt.Printf("%s%8.2f%s", red, val, reset)
				} else {
					fmt.Printf("%8.2f", val)
				}
			}
			fmt.Println()
		}
	}
}

func (simptab SimplexTable) MakeKanonView() {
	fmt.Println("Задача:")
	fmt.Printf("%s = ", simptab.Basis[len(simptab.Basis)-1])

	lastRow := simptab.Table[len(simptab.Table)-1]
	for j, val := range lastRow {
		if !simptab.View {
			val = -val
		}
		if j == 0 {
			continue //пропускаем свободный член
		}
		if val == 0 {
			continue
		}

		if j > 1 && val > 0 {
			fmt.Print(" + ")
		} else if val < 0 {
			fmt.Print(" - ")
			val = -val
		}

		fmt.Printf("%.2g * %s", val, simptab.Headers[j])
	}
	if simptab.View {
		fmt.Println(" MAX")
	} else {
		fmt.Println(" MIN")
	}

	fmt.Println()

	for i, val := range simptab.Table {
		if i != len(simptab.Table)-1 {

			first := true
			for j, element := range val {
				if j == 0 {
					continue
				}
				if element == 0 {
					continue 
				}

				if !first {
					if element > 0 {
						fmt.Print(" + ")
					} else {
						fmt.Print(" - ")
						element = -element
					}
				} else {
					// первый элемент — без ведущего "+"
					if element < 0 {
						fmt.Print(" - ")
						element = -element
					}
					first = false
				}

				fmt.Printf("%.2g * %s", element, simptab.Headers[j])
			}
			fmt.Printf(" + %s = %.2g", simptab.Basis[i], simptab.Table[i][0])
			fmt.Println()
		}

	}
}

// цель оптимизации
// true - max
// false - min
func NewTable(c_vector []float64, b_vector []float64, a_matrix [][]float64, view bool) *SimplexTable {
	xCount := len(c_vector)
	bCount := len(b_vector)

	headers := []string{"S"}
	for i := 1; i <= xCount; i++ {
		headers = append(headers, "X"+strconv.Itoa(i))
	}

	basis := make([]string, 0, bCount+1)
	for i := 1; i <= bCount; i++ {
		basis = append(basis, "X"+strconv.Itoa(xCount+i))
	}
	basis = append(basis, "F")

	table := make([][]float64, bCount+1)
	for i := 0; i < bCount; i++ {
		table[i] = make([]float64, xCount+1)
		table[i][0] = b_vector[i]

		for j := 0; j < xCount; j++ {
			table[i][j+1] = a_matrix[i][j]
		}
	}
	table[bCount] = make([]float64, xCount+1)
	table[bCount][0] = 0
	for i := 0; i < xCount; i++ {
		if view {
			table[bCount][i+1] = c_vector[i]
		} else {
			table[bCount][i+1] = -c_vector[i]
		}

	}

	return &SimplexTable{
		Basis:   basis,
		Headers: headers,
		Table:   table,
		View:    view,
	}
}

func (simptab *SimplexTable) DeepCopy() *SimplexTable {
	newBasis := make([]string, len(simptab.Basis))
	copy(newBasis, simptab.Basis)

	newHeaders := make([]string, len(simptab.Headers))
	copy(newHeaders, simptab.Headers)

	newTable := make([][]float64, len(simptab.Table))
	for i := range simptab.Table {
		newTable[i] = make([]float64, len(simptab.Table[i]))
		copy(newTable[i], simptab.Table[i])
	}

	return &SimplexTable{
		Basis:   newBasis,
		Headers: newHeaders,
		Table:   newTable,
		View:    simptab.View,
	}
}

func (simptab SimplexTable) CheckSupportSolution() bool {
	for i := 0; i < len(simptab.Table); i++ {
		if simptab.Table[i][0] < 0 {
			return false
		}
	}
	return true
}

func (simptab *SimplexTable) FindSupportSolution() bool {// додделать

	fmt.Println("Оппорное решение: ")
	for i, h := range simptab.Headers {
		if i == 0 {
			fmt.Printf("F = %v; ", simptab.Table[len(simptab.Table)-1][0])
			continue
		}
		fmt.Printf("%v = 0; ", h)
	}
	for i, h := range simptab.Basis {
		if i != len(simptab.Basis)-1 {
			fmt.Printf("%v = %.2g; ", h, simptab.Table[i][0])
		}
	}
	fmt.Println()

	return true
}

func (simptab *SimplexTable) CheckOptimized() bool {
	for _, val := range simptab.Table[len(simptab.Table)-1] {
		if val > 0 {
			return false
		}
	}
	return true
}

//targetLine - целевая строка - последняя строка

func (simptab *SimplexTable) DoSimplexMethod() {

	checkTable := simptab.DeepCopy()
	simptab.FindSupportSolution()

	var count int = 1
	var sizeTable int = len(simptab.Table)

	for !simptab.CheckOptimized() {

		if count > 100 {
			fmt.Println("Превышено число итераций")
			return
		}

		var targetLine []float64 = simptab.Table[sizeTable-1]
		var max float64 = 0
		razrech_stolb := -1
		var razrech_string int

		for j := 1; j < len(targetLine); j++ {
			if targetLine[j] > max {
				max = targetLine[j]
				razrech_stolb = j
			} //
		}

		if razrech_stolb == -1 {
			fmt.Println("Не нашел разрешающий столбец") // нет входящего столбца
			break
		}

		max = math.MaxFloat64
		// первый столбец это столбец начальных ограничений
		for i, val := range simptab.Table {
			if i != sizeTable-1 && val[razrech_stolb] > 0 && val[0] > 0 && val[0]/val[razrech_stolb] < max { 
				max = val[0] / val[razrech_stolb]
				razrech_string = i
			}
		}

		if max == math.MaxFloat64 { // если разрешающей строки нет
			fmt.Println("Решения нет")
			return
		}

		pivot := simptab.Table[razrech_string][razrech_stolb] //pivot - разрешающее значение

		// новая таблица
		newTable := make([][]float64, sizeTable)
		for i := range newTable {
			newTable[i] = make([]float64, len(targetLine))
		}

		for i := 0; i < sizeTable; i++ {
			for j := 0; j < len(targetLine); j++ {
				switch {
				case i == razrech_string && j == razrech_stolb:
					newTable[i][j] = 1.0 / pivot
				case i == razrech_string && j != razrech_stolb:
					newTable[i][j] = simptab.Table[i][j] / pivot
				case i != razrech_string && j == razrech_stolb:
					newTable[i][j] = -simptab.Table[i][j] / pivot
				default:
					newTable[i][j] = simptab.Table[i][j] - (simptab.Table[i][razrech_stolb] * simptab.Table[razrech_string][j] / pivot)
				}
			}
		}

		//заменяем таблицу шапки
		simptab.Table = newTable
		simptab.Basis[razrech_string], simptab.Headers[razrech_stolb] = simptab.Headers[razrech_stolb], simptab.Basis[razrech_string]

		fmt.Printf("----------- Итерация %v -----------", count)
		fmt.Println()
		fmt.Printf("Разрешающий столбец = %v, Разрешающая строка = %v", razrech_stolb+1, razrech_string+1)
		fmt.Println()
		count++
		simptab.Print(razrech_string, razrech_stolb)
		fmt.Println()

	}
	simptab.getAnswerAndCheck(checkTable)

}

func(simptab *SimplexTable) getAnswerAndCheck( checkTable *SimplexTable){

	finalOtvet := make(map[string]float64)

	fmt.Println("Оптимизированное решение: ")
	for i, h := range simptab.Headers {
		if i == 0 {
			if simptab.View {
				fmt.Printf("F = %.4g; ", -simptab.Table[len(simptab.Table)-1][0]) // в случае max необходимо инвертировать ответ
				continue
			} else {
				fmt.Printf("F = %.4g; ", simptab.Table[len(simptab.Table)-1][0])
				continue
			}

		}
		fmt.Printf("%s = 0; ", h)
		finalOtvet[h] = 0
	}
	for i, h := range simptab.Basis {
		if i != len(simptab.Basis)-1 {
			fmt.Printf("%v = %.4g; ", h, simptab.Table[i][0])
			finalOtvet[h] = simptab.Table[i][0]
		}
	}
	fmt.Println()
	fmt.Println()

	fmt.Println("Проверка решения (система ограничений):")
	fmt.Println()
	for i := 0; i < len(checkTable.Table)-1; i++ {
		var lhs string
		var sum float64 = 0
		for j := 1; j < len(checkTable.Headers); j++ {
			coef := checkTable.Table[i][j]

			varName := checkTable.Headers[j]
			val := finalOtvet[varName]
			sum += coef * val

			// печатаем слагаемое
			if lhs != "" && coef >= 0 {
				lhs += " + "
			} else if coef < 0 {
				lhs += " - "
				coef = -coef
			}
			lhs += fmt.Sprintf("%.4g*%s", coef, varName)
		}
		sum += finalOtvet[checkTable.Basis[i]]
		fmt.Printf("%s = %.4g  (правая часть: %.4g)\n", lhs, sum, checkTable.Table[i][0])
	}
	fmt.Println("Проверка решения (целевая функция):")
	fmt.Println()
	var sum float64 = 0
	var lhs string
	targetRow := checkTable.Table[len(checkTable.Table)-1]
	for i := 1; i < len(targetRow); i++ {
		coef := targetRow[i]

		varName := checkTable.Headers[i]
		val := finalOtvet[varName]
		sum += coef * val

		// печатаем слагаемое
		if checkTable.View { // для max
			if lhs != "" && coef >= 0 {
				lhs += " + "
			} else if coef < 0 {
				lhs += " - "
				coef = -coef
			}
			lhs += fmt.Sprintf("%.4g*%s", coef, varName)
		} else { // для min
			if lhs != "" && coef >= 0 {
				lhs += " - "
			} else if coef < 0 && lhs != "" {
				lhs += " + "
				coef = -coef
			}
			lhs += fmt.Sprintf("%.4g*%s", coef, varName)
		}

	}
	if checkTable.View {
		fmt.Printf("%s = %.4g  (правая часть: %.4g)\n", lhs, sum, -simptab.Table[len(simptab.Table)-1][0])
	} else {
		fmt.Printf("%s = %.4g  (правая часть: %.4g)\n", lhs, sum, simptab.Table[len(simptab.Table)-1][0])
	}
}