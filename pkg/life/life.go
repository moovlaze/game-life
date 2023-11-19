package life

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type World struct {
	Height, Width int
	Cells         [][]bool
}

func (w *World) SaveState(filename string) error {

	if filename == "" {
		return fmt.Errorf("error")
	}

	file, err := os.Create(filename)
	defer file.Close()

	for i, arr := range w.Cells {
		str := ""
		for _, el := range arr {
			if el {
				str += "1"
			} else {
				str += "0"
			}
		}

		if i != len(w.Cells)-1 {
			str += "\n"
		}

		file.WriteString(str)
	}

	return err
}

func (w *World) LoadState(filename string) error {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cells := [][]bool{}

	countHeight := 0
	for scanner.Scan() {

		cells = append(cells, []bool{})

		for _, num := range scanner.Text() {
			cells[countHeight] = append(cells[countHeight], string(num) == "1")
		}

		countHeight++
	}

	for i := len(cells) - 1; i >= 1; i-- {
		if len(cells[i]) != len(cells[i-1]) {
			return fmt.Errorf("error: cells is empty")
		}
	}

	w.Cells = cells

	return nil
}

func (w World) String() string {
	resultStr := ""

	for _, arr := range w.Cells {

		for _, el := range arr {
			if el {
				resultStr += "\xF0\x9F\x9F\xA9"
			} else {
				resultStr += "\xF0\x9F\x94\xA5"
			}
		}

		resultStr += "\n"
	}

	return resultStr
}

func (w *World) Neighbours(x, y int) int {
	n := 0

	for i := y - 1; i <= y+1; {
		if i < 0 {
			i = len(w.Cells) - 1
		}

		for j := x - 1; j <= x+1; {
			if j == x && i == y {
				j++
				continue
			}
			if j < 0 {
				j = len(w.Cells[i]) - 1
			}
			if w.Cells[i][j] {
				n++
			}
			if j == len(w.Cells[i])-1 {
				j = -1
			}
			if j == 0 && x == len(w.Cells[i])-1 {
				break
			}

			j++
		}

		if i == len(w.Cells)-1 {
			i = -1
		}

		if i == 0 && y == len(w.Cells)-1 {
			break
		}

		i++
	}

	return n
}

func (w *World) Next(x, y int) bool {
	n := w.Neighbours(x, y)      // получим количество живых соседей
	alive := w.Cells[y][x]       // текущее состояние клетки
	if n < 4 && n > 1 && alive { // если соседей двое или трое, а клетка жива
		return true // то следующее состояние — жива
	}
	if n == 3 && !alive { // если клетка мертва, но у неё трое соседей
		return true // клетка оживает
	}

	return false // в любых других случаях — клетка мертва
}

func NextState(oldWorld, newWorld *World) {
	// переберём все клетки, чтобы понять, в каком они состоянии
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			// для каждой клетки получим новое состояние
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}
}

func (w *World) RandInit(percentage int) {
	// Количество живых клеток
	numAlive := percentage * w.Height * w.Width / 100
	// Заполним живыми первые клетки
	w.fillAlive(numAlive)
	// Получаем рандомные числа
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// Рандомно меняем местами
	for i := 0; i < w.Height*w.Width; i++ {
		randRowLeft := r.Intn(w.Width)
		randColLeft := r.Intn(w.Height)
		randRowRight := r.Intn(w.Width)
		randColRight := r.Intn(w.Height)

		w.Cells[randRowLeft][randColLeft] = w.Cells[randRowRight][randColRight]
	}
}

func (w *World) fillAlive(num int) {
	aliveCount := 0
	for j, row := range w.Cells {
		for k := range row {
			w.Cells[j][k] = true
			aliveCount++
			if aliveCount == num {

				return
			}
		}
	}
}

func NewWorld(height, width int) (*World, error) {
	if height < 0 || width < 0 {
		return &World{}, fmt.Errorf("[ERROR] числа отрицательные")
	}

	cells := make([][]bool, height)

	for i := range cells {
		cells[i] = make([]bool, width)
	}

	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}, nil
}
