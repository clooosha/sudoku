package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Size = 9

type Cell struct {
	Value          int
	PossibleValues []int
}

func NewCell() Cell {
	cell := Cell{}
	cell.Value = 0
	cell.PossibleValues = make([]int, Size)
	for i := 0; i < Size; i++ {
		cell.PossibleValues[i] = i + 1
	}
	return cell
}

func CopyCell(c *Cell) *Cell {
	cell := Cell{}
	cell.Value = c.Value
	cell.PossibleValues = make([]int, len(c.PossibleValues))
	for i := 0; i < len(c.PossibleValues); i++ {
		cell.PossibleValues[i] = c.PossibleValues[i]
	}
	return &cell
}

func (c *Cell) SetValue(value int) {
	if value != 0 {
		c.Value = value
		c.PossibleValues = []int{value}
	}
}

func (cell *Cell) checkValue() error {
	switch len(cell.PossibleValues) {
	case 0:
		return CellError(fmt.Sprintf("empty mass %v", cell.Value))
	case 1:
		cell.Value = cell.PossibleValues[0]
		return nil
	default:
		return nil
	}
}

func (cell Cell) String() string {
	if cell.Value == 0 {
		return fmt.Sprintf("%v", cell.PossibleValues)
	}
	return fmt.Sprintf("%v", cell.Value)
}

type CellError string

func (c CellError) Error() string {
	return string(c)
}

func deleteValue(pSlice *[]int, value int) {
	for i, elem := range *pSlice {
		if elem == value {
			if i != 0 {
				*pSlice = append((*pSlice)[:i], (*pSlice)[i+1:len(*pSlice)]...)
			} else {
				*pSlice = (*pSlice)[i+1 : len(*pSlice)]
			}
			return
		}
	}
}

func isValue(slice []int, value int) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return true
		}
	}
	return false
}

type Sudoku struct {
	array [Size][Size]Cell
}

func CopySudoku(s *Sudoku) *Sudoku {
	sudoku := Sudoku{}
	for x := 0; x < Size; x++ {
		for y := 0; y < Size; y++ {
			sudoku.array[x][y] = *CopyCell(&s.array[x][y])
		}
	}
	return &sudoku
}

func (s *Sudoku) SetData(array [Size][Size]int) {
	for x := 0; x < Size; x++ {
		for y := 0; y < Size; y++ {
			cell := NewCell()
			cell.SetValue(array[x][y])
			s.array[x][y] = cell
		}
	}
}

func (s *Sudoku) checkRow(value int, row int, column int) (bool, error) {
	isFindValue := false
	for y := 0; y < Size; y++ {
		if y != column {
			oldValue := s.array[row][y].Value
			deleteValue(&s.array[row][y].PossibleValues, value)
			if err := s.array[row][y].checkValue(); err != nil {
				return isFindValue, err
			} else {
				if (oldValue == 0) && (s.array[row][y].Value) != 0 {
					isFindValue = true
				}
			}
		}
	}
	return isFindValue, nil
}

func (s *Sudoku) checkColumn(value int, row int, column int) (bool, error) {
	isFindValue := false
	for x := 0; x < Size; x++ {
		if x != row {
			oldValue := s.array[x][column].Value
			deleteValue(&s.array[x][column].PossibleValues, value)
			if err := s.array[x][column].checkValue(); err != nil {
				return isFindValue, err
			} else {
				if (oldValue == 0) && (s.array[x][column].Value != 0) {
					isFindValue = true
				}
			}
		}
	}
	return isFindValue, nil
}

func (s *Sudoku) checkSquare(value int, row int, column int) (bool, error) {
	isFindValue := false
	startRow := row - row%(Size/3)
	startColumn := column - column%(Size/3)
	for x := startRow; x < startRow+(Size/3); x++ {
		for y := startColumn; y < startColumn+(Size/3); y++ {
			if !((x == row) && (y == column)) {
				oldValue := s.array[x][y].Value
				deleteValue(&s.array[x][y].PossibleValues, value)
				if err := s.array[x][y].checkValue(); err != nil {
					return isFindValue, err
				} else {
					if (oldValue == 0) && (s.array[x][y].Value != 0) {
						isFindValue = true
					}
				}
			}
		}
	}
	return isFindValue, nil
}

func (s *Sudoku) deleteValueInArray() error {
	for x := 0; x < Size; x++ {
		for y := 0; y < Size; y++ {
			value := s.array[x][y].Value
			if value != 0 {
				isFindValue := false
				if isFind, err := s.checkRow(value, x, y); err != nil {
					return err
				} else {
					isFindValue = isFindValue || isFind
				}
				if isFind, err := s.checkColumn(value, x, y); err != nil {
					return err
				} else {
					isFindValue = isFindValue || isFind
				}
				if isFind, err := s.checkSquare(value, x, y); err != nil {
					return err
				} else {
					isFindValue = isFindValue || isFind
				}
				if isFindValue {
					if err := s.deleteValueInArray(); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (s *Sudoku) IsDone() (int, int) {
	x := -1
	y := -1
	for x := 0; x < Size; x++ {
		for y := 0; y < Size; y++ {
			if s.array[x][y].Value == 0 {
				return x, y
			}
		}
	}
	return x, y
}

func (s *Sudoku) Calculate() error {
	if err := s.deleteValueInArray(); err != nil {
		return err
	}
	for i := 0; i < Size; i++ {
		if isFindValue := s.checkOneValueInRow(i); isFindValue {
			if err := s.deleteValueInArray(); err != nil {
				return err
			}
		}
		if isFindValue := s.checkOneValueInColumn(i); isFindValue {
			if err := s.deleteValueInArray(); err != nil {
				return err
			}
		}
		if isFindValue := s.checkOneValueInSquare(i); isFindValue {
			if err := s.deleteValueInArray(); err != nil {
				return err
			}
		}
	}
	if x, y := s.IsDone(); (x != -1) && (y != -1) {
		newS := CopySudoku(s)
		PosibleValue := s.array[x][y].PossibleValues[0]
		newS.array[x][y].SetValue(PosibleValue)
		if err := newS.Calculate(); err == nil {
			*s = *newS
			return nil
		} else {
			deleteValue(&s.array[x][y].PossibleValues, PosibleValue)
			return s.Calculate()
		}
	}
	return nil
}

func (sudoku *Sudoku) checkOneValueInRow(row int) bool {
	isFindValue := false
	for value := 1; value < Size+1; value++ {
		count := 0
		index := -1
		for i := 0; i < Size; i++ {
			if sudoku.array[row][i].Value == 0 && isValue(sudoku.array[row][i].PossibleValues, value) {
				count++
				index = i
			}
		}
		if count == 1 {
			sudoku.array[row][index].SetValue(value)
			isFindValue = true
		}
	}
	return isFindValue
}

func (sudoku *Sudoku) checkOneValueInColumn(column int) bool {
	isFindValue := false
	for value := 1; value < Size+1; value++ {
		count := 0
		index := -1
		for i := 0; i < Size; i++ {
			if sudoku.array[i][column].Value == 0 && isValue(sudoku.array[i][column].PossibleValues, value) {
				count++
				index = i
			}
		}
		if count == 1 {
			sudoku.array[index][column].SetValue(value)
			isFindValue = true
		}
	}
	return isFindValue
}

func (sudoku *Sudoku) checkOneValueInSquare(num int) bool {
	isFindValue := false
	startRow := 3 * (num / 3)
	startColumn := 3 * (num % 3)
	for value := 1; value < Size+1; value++ {
		count := 0
		indexX := -1
		indexY := -1
		for x := startRow; x < startRow+3; x++ {
			for y := startColumn; y < startColumn+3; y++ {
				if sudoku.array[x][y].Value == 0 && isValue(sudoku.array[x][y].PossibleValues, value) {
					count++
					indexX = x
					indexY = y
				}
			}
		}
		if count == 1 {
			sudoku.array[indexX][indexY].SetValue(value)
			isFindValue = true
		}
	}
	return isFindValue
}

func (s Sudoku) String() string {
	var str string
	for x := 0; x < Size; x++ {
		for y := 0; y < Size; y++ {
			str = fmt.Sprintf("%s%s ", str, s.array[x][y])
		}
		str = fmt.Sprintln(str)
	}
	return str
}

func ReadSudoku(array *[Size][Size]int) int {
	reader := bufio.NewReader(os.Stdin)
	for x := 0; x < Size; x++ {
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		row := strings.Split(text, " ")
		for y := 0; y < len(row); y++ {
			if row[y] == "_" {
				array[x][y] = 0
			} else {
				if value, err := strconv.ParseInt(row[y], 0, 64); err != nil {
					fmt.Println(err)
					return -1
				} else {
					array[x][y] = int(value)
				}
			}
		}

	}
	return 0
}

func main() {
	var array [Size][Size]int
	if err := ReadSudoku(&array); err != 0 {
		fmt.Println("Error")
		return
	}
	var sudoku Sudoku
	sudoku.SetData(array)
	if err := sudoku.Calculate(); err != nil {
		fmt.Println(err)
	}
	fmt.Print(sudoku)
}
