// Головний пакет програми
package main

// Підключення необхідних бібліотек
import (
	"github.com/gdamore/tcell/v2"
)

// Завдання власних типів даних

// Комірка консолі з символом та його кольором
type Cell struct {
	// Символ
	Rune  rune
	// Кольор
	Style tcell.Style
}
// Клітина ігрового поля, що складається з двох комірок - видимої для гравця і прихованої його оппонента
type GameCell struct {
	// Видима комірка
	VisibleCell Cell
	// Прихована комірка
	HiddenCell  Cell
}
// Корабль, для якого є головна клітина, тіло та напрямок
type Boat struct {
	// Координата х голови
	headX        int
	// Координата у голови
	headY        int
	// Масив всіх клітин зі станом їх пошкодженості (true - пошкоджена, false - ціла)
	partsDamaged []bool
	// Напрямок (true - вертикальний, false - горизонтальний)
	orientation  bool
}

// Завдання констант
const (
	// Розмір вікна по x
	screenSizeX     = 120
	// Розмір вікна по y 
	screenSizeY     = 30
	// Розмір ігрового поля
	battlefieldSize = 10
)

// Завдання глобальних змінних

// Початок гри (true - гра почалася, false - знаходження в меню)
var gameBegin bool = false

// Точка входу програми
func main() {
	screen := createScreen()
	beginMenu(screen)
}
