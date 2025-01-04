// Головний пакет програми
package main

// Підключення необхідних бібліотек
import (
	"io"
	"os"
	"strings"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

// Завдання власних типів даних

// Дані збереження
type SaveData struct {
	// Ігрове поле першого гравця
	player1GameField [][]GameCell
	// Ігрове поле другого гравця
	player2GameField [][]GameCell
	// Гравець який ходить (true - гравець 1, false - гравець 2)
	playerToMove bool
	// Пересування корабля гравця (true - переміщено, false - не переміщено)
	boatMoved bool
}

// Конвертація стилю у число для збереження (стиль) (число)
func tcellStyleToInt(style tcell.Style) int {
	// Вибір стилю відповідно до вхідного параметру
	switch style {
	case tcell.StyleDefault.Foreground(tcell.ColorGreen):
		return 1
	case tcell.StyleDefault.Foreground(tcell.ColorRed):
		return 2
	case tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor()):
		return 3
	case tcell.StyleDefault.Foreground(tcell.ColorBlue):
		return 4
	default:
		return 0
	}
}

// Конвертація числа для збереження у стиль (число) (стиль)
func intToTcellStyle(int int) tcell.Style {
	// Вибір числа відповідно до вхідного параметру
	switch int {
	case 1:
		return tcell.StyleDefault.Foreground(tcell.ColorGreen)
	case 2:
		return tcell.StyleDefault.Foreground(tcell.ColorRed)
	case 3:
		return tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor())
	case 4:
		return tcell.StyleDefault.Foreground(tcell.ColorBlue)
	default:
		return tcell.StyleDefault
	}
}

// Конвертація ігрового поля у текст для збереження (ігрове поле) (рядок)
func gameFieldToText(gameField [][]GameCell) string {
	// Вихідний текстовий рядок
	gameFieldText := ""

	// Проходження через всі клітини ігрового поля
	for _, row := range gameField {
		for _, cell := range row {
			// Запис через кому значень символу та кольору видимої та прихованої клітини ігрового поля з крапкокомою в кінці
			gameFieldCellText := string(cell.VisibleCell.Rune) + "," + strconv.Itoa(tcellStyleToInt(cell.VisibleCell.Style)) + "," +
			                     string(cell.HiddenCell.Rune) + "," + strconv.Itoa(tcellStyleToInt(cell.HiddenCell.Style)) + ";"
			gameFieldText += gameFieldCellText
		}
		// Запис у вихідний рядок без останнього роздільного символу рядка з переходом на новий в кінці
		gameFieldText = gameFieldText[:len(gameFieldText)-1] + "\n"
	}

	// Повернення отриманого текстового рядку без останнього роздільного символу
	return gameFieldText[:len(gameFieldText)-1]
}

// Конвертація тексту для збереження у ігрове поле (рядок) (ігрове поле)
func textToGameField(text string) [][]GameCell {
	// Ініціалізація вихідної ігрової матриці
	gameField := make([][]GameCell, battlefieldSize*2)
	for i := 0; i < battlefieldSize*2; i++ {
		gameField[i] = make([]GameCell, battlefieldSize)
	}

	// Розділення вхідного тексту на рядки
	rowsText := strings.Split(text, "\n")
	// Проходження через кожний рядок
	for i, rowText := range rowsText {
		// Розділення рядка на підрядки за крапкокомою
		cellsText := strings.Split(rowText, ";")
		// Проходження через кожну клітину
		for j, cellText := range cellsText {
			// Розділення клітини на підрядки за комою
			gameCellTextArray := strings.Split(cellText, ",")

			// Встановлення значень символу та кольору видимої та прихованої клітини ігрового поля
			gameField[i][j].VisibleCell = Cell{
				Rune:  []rune(gameCellTextArray[0])[0],
				Style: intToTcellStyle(int(gameCellTextArray[1][0] - '0')),
			}
			gameField[i][j].HiddenCell = Cell{
				Rune:  []rune(gameCellTextArray[2])[0],
				Style: intToTcellStyle(int(gameCellTextArray[3][0] - '0')),
			}
		}
	}

	// Повернення отриманої ігрової матриці
	return gameField
}

// Визначення, чи існує дійсний файл збереження (назва файлу) (існування (true - існує, false - відсутній))
func saveExists(filename string) bool {
	// Отримання файлу збереження
	file, err := os.Open(filename)
	// Якщо сталася помилка при відкритті
	if err != nil {
		return false
	}
	// Відкладне закриття відкритого файлу
	defer file.Close()

	// Отримання мета-даних файлу
	fileInfo, err := file.Stat()
	// Якщо сталася помилка при отриманні
	if err != nil {
		return false
	}

	// Якщо файл пустий
	if fileInfo.Size() == 0 {
		return false
	  // Якщо є актуальне збереження гри
	} else {
		return true
	}
}

// Створення нового або перезапис існуючого на пустий файл збереження (назва файлу)
func createSaveFile(filename string) {
	// Створення файлу збереження
	file, _ := os.Create(filename)
	// Відкладне закриття створеного файлу
	defer file.Close()
}

// Читання файла збереження (рядок) (дані збереження)
func readSaveFile(filename string) SaveData {
	// Отримання файлу збереження
	file, _ := os.Open(filename)
	// Відкладне закриття відкритого файлу
	defer file.Close()

	// Читання файлу у масив байтів
	fileContent, _ := io.ReadAll(file)
	// Масив частин перетвореного з масиву байтів текстового рядку, виділених за роздільником "\n---\n"
	gameFieldsTextArray := strings.Split(string(fileContent), "\n---\n")

	// Повернення даних збереження перетворених у програмні типи даних
	return SaveData{
		textToGameField(gameFieldsTextArray[0]),
		textToGameField(gameFieldsTextArray[1]),
		// Анонімні функції отримання з рядку булевого значення
		func() bool {
			result, _ := strconv.ParseBool(gameFieldsTextArray[2])
			return result
		}(),
		func() bool {
			result, _ := strconv.ParseBool(gameFieldsTextArray[3])
			return result
		}(),
	}
}

// Запис у файл збереження (рядок, дані збереження)
func writeSaveFile(filename string, saveData SaveData) {
	// Створення нового або перезапис існуючого файлу збереження
	file, _ := os.Create(filename)
	// Відкладне закриття відкритого файлу
	defer file.Close()

	// Перетворення вхідних типів даних збереження у текстовий та їх з'єднання між собою через роздільник "\n---\n"
	gameFieldsTextArray := gameFieldToText(saveData.player1GameField) + "\n---\n" +
	                       gameFieldToText(saveData.player2GameField) + "\n---\n" +
						   strconv.FormatBool(saveData.playerToMove) + "\n---\n" +
						   strconv.FormatBool(saveData.boatMoved)

	// Запис рядка у файл
	_, _ = file.WriteString(gameFieldsTextArray)
}