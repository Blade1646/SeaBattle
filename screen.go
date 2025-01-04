// Головний пакет програми
package main

// Підключення необхідних бібліотек
import (
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Створення екрану консолі (екран)
func createScreen() tcell.Screen {
	// Створення нового екрану
	screen, _ := tcell.NewScreen()
	// Ініціалізація екрану
	screen.Init()
	// Активація миші
	screen.EnableMouse()
	// Повернення створеного екрану
	return screen
}

// Друк рядку зліва (екран, рядок консолі, рядок тексту)
func printStringLeft(screen tcell.Screen, row int, str string) {
	// Друк починаючи з початку рядка
	for i, r := range str {
		screen.SetContent(i+1, row, r, nil, tcell.StyleDefault)
	}
}

// Друк рядку справа (екран, рядок консолі, рядок тексту)
func printStringRight(screen tcell.Screen, row int, str string) {
	// Обчислення початкового стовпця для орієнтації рядку праворуч відносно екрану
	startCol := screenSizeX - len(str) - 1
	for i, r := range str {
		screen.SetContent(startCol+i, row, r, nil, tcell.StyleDefault)
	}
}

// Друк рядку посередині (екран, рядок консолі, рядок тексту)
func printStringMiddle(screen tcell.Screen, row int, str string) {
	// Обчислення початкового стовпця для орієнтації рядку посередині екрану
	startCol := (screenSizeX - len(str)) / 2
	for i, r := range str {
		screen.SetContent(startCol+i, row, r, nil, tcell.StyleDefault)
	}
}

// Отримання індексу підрядка в рядку (екран, рядок консолі, рядок тексту) (індекс)
func getStringIndex(screen tcell.Screen, row int, str string) int {
	// Отримання тексту всього рядка
	rowAsString := ""
	for j := 0; j < screenSizeX; j++ {
		r, _, _, _ := screen.GetContent(j, row)
		rowAsString += string(r)
	}

	// Повернення індексу підрядка в рядку
	return strings.Index(rowAsString, str)
}

// Малювання комірок консолі за кольором (екран, стовпець, рядок, довжина, колір)
func paintCellsRowWithColor(screen tcell.Screen, startCol int, row int, length int, color tcell.Color) {
	for i := startCol; i < startCol+length; i++ {
		r, _, _, _ := screen.GetContent(i, row)
		screen.SetContent(i, row, r, nil, tcell.StyleDefault.Foreground(color))
	}
}

// Динамічне очищення всього екрану (екран)
func eraseScreenDynamically(screen tcell.Screen) {
	// Ініціалізація матриці очищених комірок консолі (true - стерта, false - роздрукована)
	field := make([][]bool, screenSizeX)
	for i := range field {
		field[i] = make([]bool, screenSizeY)
		for j := range field[i] {
			field[i][j] = false
		}
	}
	// Лічильник стертих комірок
	counter := 0
	// Поки всі комірки не очіщені
	for counter < screenSizeX*screenSizeY {
		// Випадковий вибір координат х та у комірки
		randCol := rand.Intn(screenSizeX)
		randRow := rand.Intn(screenSizeY)
		// Якщо комірка ще не була очищена
		if !field[randCol][randRow] {
			field[randCol][randRow] = true
			counter++
			screen.SetContent(randCol, randRow, ' ', nil, tcell.StyleDefault)
			// Невелика затримка для забезпечення динамічної послідовності очистки
			screen.Show()
			time.Sleep(time.Millisecond / 10)
		}
	}
}

// Динамічне друкування всього екрану (екран, верх (true - друкування починаючи зверху, false - друкування починаючи зліва))
func drawScreenDynamically(screen tcell.Screen, top bool) {
	// Збереження початкових значень комірок консолі та її очищення
	initialData := make([][]Cell, screenSizeX)
	for i := 0; i < screenSizeX; i++ {
		initialData[i] = make([]Cell, screenSizeY)
		for j := 0; j < screenSizeY; j++ {
			rune, _, style, _ := screen.GetContent(i, j)
			initialData[i][j] = Cell{rune, style}
			screen.SetContent(i, j, ' ', nil, tcell.StyleDefault)
		}
	}
	screen.Show()
	// Послідовний друк збережених значень комірок консолі по рядках
	if top {
		for i := 0; i < screenSizeY; i++ {
			for j := 0; j < screenSizeX; j++ {
				cell := initialData[j][i]
				screen.SetContent(j, i, cell.Rune, nil, cell.Style)
			}
			// Невелика затримка для забезпечення динамічної послідовності друку
			screen.Show()
			time.Sleep(time.Millisecond * 40)
		}
		// Послідовний друк збережених значень комірок консолі по стовпцях
	} else {
		for i := 0; i < screenSizeX; i++ {
			for j := 0; j < screenSizeY; j++ {
				cell := initialData[i][j]
				screen.SetContent(i, j, cell.Rune, nil, cell.Style)
			}
			// Невелика затримка для забезпечення динамічної послідовності друку
			screen.Show()
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// Друк автору (екран)
func printAuthor(screen tcell.Screen) {
	printStringRight(screen, screenSizeY-2, "Vladyslav Rakitin")
	printStringRight(screen, screenSizeY-3, "PZPI-21-11 student")
	printStringRight(screen, screenSizeY-4, "Developed by:")
}

// Очищення панелі екрану виходу з гри (екран)
func clearExitScreen(screen tcell.Screen) {
	// Проходження в циклі всіх комірок області
	for i := screenSizeY - 4; i <= screenSizeY-2; i++ {
		for j := 1; j < screenSizeX/2; j++ {
			screen.SetContent(j, i, ' ', nil, tcell.StyleDefault)
		}
	}
}

// Малювання панелі екрану виходу з гри (екран)
func paintExitScreen(screen tcell.Screen) {
	for i := screenSizeY - 4; i <= screenSizeY-2; i++ {
		// Отримання індексів назв клавіш для малювання
		escIndex := getStringIndex(screen, i, "Esc")
		enterIndex := getStringIndex(screen, i, "Enter")

		paintCellsRowWithColor(screen, escIndex, i, len("Esc"), tcell.ColorRed)
		if enterIndex != -1 {
			paintCellsRowWithColor(screen, enterIndex, i, len("Enter"), tcell.ColorGreen)
		}
	}
}

// Встановлення базового режиму виходу з гри (екран)
func setDefaultExitState(screen tcell.Screen) {
	clearExitScreen(screen)
	printStringLeft(screen, screenSizeY-4, "   Esc   Press to exit game")
	paintExitScreen(screen)
}

// Встановлення режиму підказки виходу з гри (екран, збереження даних (true - відбудеться, false - не станеться))
func setConfirmExitState(screen tcell.Screen, dataSaved bool) {
	clearExitScreen(screen)

	// Текст підказки
	confirmExitQuestionString := "Are you sure you want to exit?"
	confirmExitOptionsString := "   Enter   Accept exit    Esc   Decline exit"

	// Приведення рядків до однієї довжини
	diff := len(confirmExitOptionsString) - len(confirmExitQuestionString)
	confirmExitQuestionStringLong := strings.Repeat(" ", diff/2) + confirmExitQuestionString + strings.Repeat(" ", diff-diff/2)

	// Друк тексту підказки
	printStringLeft(screen, screenSizeY-4, confirmExitQuestionStringLong)
	printStringLeft(screen, screenSizeY-2, confirmExitOptionsString)

	// Оповіщення якщо при виході з гри поточні дані будуть збережені
	if dataSaved {
		dataSavedString := "(your data will be saved!)"
		diff := len(confirmExitOptionsString) - len(dataSavedString)
		dataSavedStringLong := strings.Repeat(" ", diff/2) + dataSavedString + strings.Repeat(" ", diff-diff/2)
		printStringLeft(screen, screenSizeY-3, dataSavedStringLong)
	}

	paintExitScreen(screen)
}
