// Головний пакет програми
package main

// Підключення необхідних бібліотек
import (
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

// Завдання власних типів даних

// Перевірка можливості та виконання руху корабля
type PossibleMovement struct {
	// Можливість руху (true - дозволено виконати дію, false - заборонено)
	canMove func([][]GameCell, Boat) bool
	// Виконання руху
	move    func([][]GameCell, Boat) Boat
}

// Масив можливих рухів корабля
var possibleMovements = []PossibleMovement{
	{canMoveUp, moveUp},
	{canMoveDown, moveDown},
	{canMoveLeft, moveLeft},
	{canMoveRight, moveRight},
	{canRotateLeft, rotateLeft},
	{canRotateRight, rotateRight},
}

// Завдання констант
const (
	// Початкова координата х меж поля першого гравця
	startCol1 int = (screenSizeX*0.5 - (battlefieldSize * 2) - 2) / 2
	// Початкова координата у меж поля першого гравця
	startRow1 int = (screenSizeY - battlefieldSize - 2) / 2

	// Початкова координата х меж поля другого гравця
	startCol2 int = (screenSizeX*1.5 - (battlefieldSize * 2) - 2) / 2
	// Початкова координата у меж поля другого гравця
	startRow2 int = (screenSizeY - battlefieldSize - 2) / 2
)

// Завдання глобальних змінних

// Ігрова матриця першого гравця
var player1GameField [][]GameCell
// Ігрова матриця другого гравця
var player2GameField [][]GameCell
// Етап гри (0 - очікування підтвердження перед початком налаштування, 0.5 - налаштування кораблів першим гравцем, 1 - очікування підтвердження перед початком налаштування другим гравцем, 1.5 - налаштування кораблів другим гравцем, 2 - завершення налаштування та очікування підтвердження перед наступним кроком, 2.1 - хід першого гравця, 2.2 - хід другого гравця, 2.3 - хід ШІ, 3 - кінець гри)
var gameState float32
// Гравець який ходить (true - гравець 1, false - гравець 2)
var playerToMove bool
// Пересування корабля гравця (true - переміщено, false - не переміщено)
var boatMoved bool
// Обраний корабель гравця що ходить
var selectedBoat Boat

// Друк меж поля (екран, гравець (true - гравець 1, false - гравець 2))
func printGameFieldBorder(screen tcell.Screen, player bool) {
	// Встановлення початкових координат поля відповідно до вхідного параметру гравця
	var startCol int
	var startRow int
	if player {
		startCol = startCol1
		startRow = startRow1
	} else {
		startCol = startCol2
		startRow = startRow2
	}

	// Малювання верхнього та нижнього кордонів поля
	for i := 0; i <= battlefieldSize*2+1; i++ {
		screen.SetContent(startCol+i, startRow, '▄', nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
		screen.SetContent(startCol+i, startRow+battlefieldSize+1, '▀', nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	}

	// Малювання лівого та правого кордонів поля
	for i := 1; i <= battlefieldSize; i++ {
		screen.SetContent(startCol, startRow+i, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
		screen.SetContent(startCol+(battlefieldSize*2)+1, startRow+i, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	}

	// Встановлення позначення гравця та його кольору відповідно до вхідного параметру гравця
	var playerNumber string
	var playerColor tcell.Color
	if player {
		playerNumber = "Player 1"
		playerColor = tcell.ColorBlue
	} else {
		playerNumber = "Player 2"
		playerColor = tcell.ColorRed
	}

	// Підписання поля
	playerFieldString := playerNumber + " " + "field"
	diff := battlefieldSize*2 + 2 - len(playerFieldString)
	playerFieldStringLong := strings.Repeat(" ", diff/2) + playerFieldString + strings.Repeat(" ", diff-diff/2)
	for i, r := range playerFieldStringLong {
		screen.SetContent(startCol+i, startRow-2, r, nil, tcell.StyleDefault)
	}
	paintCellsRowWithColor(screen, getStringIndex(screen, startRow-2, playerNumber), startRow-2, len(playerNumber), playerColor)
}

// Оновлення візуального інтерфейсу гри (екран)
func updateInterface(screen tcell.Screen) {
	printPlayerTitle(screen)
	printPlayerControls(screen)
}

// Друк заголовчого підпису (екран)
func printPlayerTitle(screen tcell.Screen) {
	clearPlayerTitleScreen(screen)

	// Встановлення позначення гравця що ходить та його кольору
	var playerNumber string
	var playerColor tcell.Color
	if playerToMove {
		playerNumber = "Player 1"
		playerColor = tcell.ColorBlue
	} else {
		playerNumber = "Player 2"
		playerColor = tcell.ColorRed
	}
	// Встановлення дії гравця відповідно до етапу гри
	var playerAction string
	if gameState < 2 {
		playerAction = "setup"
	} else if gameState < 3 {
		playerAction = "move"
	} else {
		playerAction = "wins!🥇"
	}

	// Друк підпису
	printStringMiddle(screen, 1, playerNumber+" "+playerAction)
	paintCellsRowWithColor(screen, getStringIndex(screen, 1, playerNumber), 1, len(playerNumber), playerColor)

	// Вибір підказки відповідно до етапу гри
	switch gameState {
	// Очікування підтвердження дії
	case 0, 1, 2:
		printStringMiddle(screen, 2, "Press [Space] when ready")
	// Налаштування кораблів перед початком ігрового процесу
	case 0.5, 1.5:
		printStringMiddle(screen, 2, "Press [Space] when finished")
		printStringMiddle(screen, 4, "Setup your boats!")
	// Здійснення ходу - вибір, рух та постріл у ворожий корабль
	case 2.1, 2.2:
		printStringMiddle(screen, 2, "Press [Space] if finished")
		printStringMiddle(screen, 4, "Select, move your OR/AND fire enemy's boat!")
	// Перемога гравця, кінець гри
	case 3:
		printStringMiddle(screen, 2, "Press [Space] to exit")
	}

	paintCellsRowWithColor(screen, getStringIndex(screen, 2, "Space"), 2, len("Space"), tcell.ColorWhite.TrueColor())
}

// Друк підказки з елементами управління (екран)
func printPlayerControls(screen tcell.Screen) {
	clearPlayerControlsScreen(screen)

	// Якщо етап гри нецілочисельний (дійовий) та хід не за ШІ
	if float64(gameState) != math.Trunc(float64(gameState)) && gameState != 2.3 {
		// Масиви можливих дій та команд їх виконання
		actions := []string{"Select", "Move up", "Move down", "Move left", "Move right", "Rotate left", "Rotate right"}
		var controls []string
		if playerToMove {
			controls = []string{"LMB", "W", "S", "A", "D", "Q", "E"}
		} else {
			controls = []string{"LMB", "8", "5", "4", "6", "7", "9"}
		}

		// Якщо ігровий процес було розпочато
		if gameState > 2 {
			// Додавання опції пострілу
			actions[0] += "/Fire"
		}

		// Пошук найбільшої довжини дій
		maxActionLength := 0
		for _, action := range actions {
			if len(action) > maxActionLength {
				maxActionLength = len(action)
			}
		}
		// Пошук найбільшої довжини команд
		maxControlLength := 0
		for _, control := range controls {
			if len(control) > maxControlLength {
				maxControlLength = len(control)
			}
		}

		// Приведення всіх елементів до однієї довжини
		for i, action := range actions {
			diff := maxActionLength - len(action)
			actions[i] = strings.Repeat(" ", diff/2) + action + strings.Repeat(" ", diff-diff/2)
		}
		for i, control := range controls {
			diff := maxControlLength - len(control)
			controls[i] = strings.Repeat(" ", diff/2) + control + strings.Repeat(" ", diff-diff/2)
		}

		// Друк підказки посередині екрану
		printStringMiddle(screen, startRow1+2, controls[0]+" "+actions[0])
		for i := 1; i < len(actions); i++ {
			printStringMiddle(screen, startRow1+2+1+i, controls[i]+" "+actions[i])
		}
	}
}

// Очищення частини екрану з заголовчим підписом (екран)
func clearPlayerTitleScreen(screen tcell.Screen) {
	for i := 1; i <= 4; i++ {
		for j := 1; j < screenSizeX-2; j++ {
			screen.SetContent(j, i, ' ', nil, tcell.StyleDefault)
		}
	}
}

// Очищення частини екрану з підказкою до елементів управління (екран)
func clearPlayerControlsScreen(screen tcell.Screen) {
	for i := startRow1 + 1; i <= startRow1+battlefieldSize; i++ {
		for j := startCol1 + battlefieldSize*2 + 1 + 1; j < startCol2-1; j++ {
			screen.SetContent(j, i, ' ', nil, tcell.StyleDefault)
		}
	}
}

// Визначення, чи знаходиться клітина в межах поля (координата х клітни, координата у клітини) (true - в межах, false - за межами)
func inGameField(cellX int, cellY int) bool {
	// Якщо значення є меншим за 0 або перевищує розмір ігрового поля
	if cellX < 0 || cellX >= battlefieldSize*2 || cellY < 0 || cellY >= battlefieldSize {
		return false
	} else {
		return true
	}
}

// Визначення, чи є клітина частиною корабля (ігрове поле, координата х клітни, координата у клітини) (true - корабль, false - вільна ділянка)
func isBoat(gameField [][]GameCell, cellX int, cellY int) bool {
	if !inGameField(cellX, cellY) {
		return false
	}
	if gameField[cellX][cellY].VisibleCell.Rune == '█' {
		return true
	}
	return false
}

// Визначення, чи є корабль пошкодженим (корабль) (true - пошкоджений, false - цілий)
func isDamaged(boat Boat) bool {
	for i := 0; i < len(boat.partsDamaged); i++ {
		if boat.partsDamaged[i] {
			return true
		}
	}
	return false
}

// Визначення, чи є корабль знищеним (корабль) (true - знищений, false - у грі)
func isDestroyed(boat Boat) bool {
	for i := 0; i < len(boat.partsDamaged); i++ {
		if !boat.partsDamaged[i] {
			return false
		}
	}
	return true
}

// Визначення, чи є корабль пересувним (ігрове поле, корабль) (true - пересувний, false - нерухомий)
func isMoveable(gameField [][]GameCell, boat Boat) bool {
	for _, possibleMovement := range possibleMovements {
		if possibleMovement.canMove(gameField, boat) {
			return true
		}
	}
	return false
}

// Визначення, чи є клітина ділянки вільною для розміщення корабля (ігрове поле, координата х клітини ділянки, координата у клітини ділянки) (true - вільна, false - зайнята)
func isAreaFree(gameField [][]GameCell, areaX int, areaY int) bool {
	if !inGameField(areaX, areaY) {
		return false
	}
	// Перебір суміжних клітин
	for i := areaX - 2; i <= areaX+2; i += 2 {
		for j := areaY - 1; j <= areaY+1; j++ {
			if isBoat(gameField, i, j) {
				return false
			}
		}
	}
	return true
}

// Визначення, чи існують пересувні кораблі на ігровому полі (ігрове поле) (true - існують, false - відсутні)
func areMoveableBoatsExist(gameField [][]GameCell) bool {
	for x, row := range gameField {
		for y, cell := range row {
			// Якщо корабль та не знищений
			if cell.VisibleCell.Rune == '█' && cell.VisibleCell.Style != tcell.StyleDefault.Foreground(tcell.ColorRed) {
				if isMoveable(gameField, getBoatByCell(gameField, x-x%2, y)) {
					return true
				}
			}
		}
	}
	return false
}

// Визначення, чи залишаються непоцілені клітини на ігровому полі (ігрове поле) (true - залишаються, false - немає)
func areUnfiredCellsExist(gameField [][]GameCell) bool {
	for _, row := range gameField {
		for _, cell := range row {
			// Якщо в клітину не стріляли
			if cell.HiddenCell.Style == tcell.StyleDefault {
				return true
			}
		}
	}
	return false
}

// Визначення, чи всі кораблі на ігровому полі знищені (ігрове поле) (true - кораблі знищені, false - залишаються кораблі у грі)
func isFieldCleared(gameField [][]GameCell) bool {
	for _, row := range gameField {
		for _, cell := range row {
			// Якщо корабль та не знищений
			if cell.VisibleCell.Rune == '█' && cell.VisibleCell.Style != tcell.StyleDefault.Foreground(tcell.ColorRed) {
				return false
			}
		}
	}
	return true
}

// Визначення, чи може бути розміщений корабль на ігровому полі (ігрове поле, корабль) (true - може, false - не може)
func canPlaceBoat(gameField [][]GameCell, boat Boat) bool {
	// Змінні для поступового проходження через весь корабль
	counterX := boat.headX
	counterY := boat.headY
	sizeCounter := 0
	
	// Поки залишаються неопрацьовані частини корабля
	for sizeCounter < len(boat.partsDamaged) {
		if isAreaFree(gameField, counterX, counterY) {
			// Інкрементація лічильників
			if boat.orientation {
				counterY++
			} else {
				counterX = counterX + 2
			}
			sizeCounter++
		} else {
			return false
		}
	}
	return true
}

// Отримання всього корабля за координатою (ігрове поле, координата х клітини, координата у клітини) (корабль)
func getBoatByCell(gameField [][]GameCell, cellX int, cellY int) Boat {
	// Поля майбутньої структури корабля
	var boatHeadX int
	var boatHeadY int
	var boatOrientation bool
	var boatPartsDamaged []bool

	// Змінні для поступового проходження через весь корабль
	var counterX int
	var counterY int

	// Встановлення верхньої лівої клітини як головної
	for counterX = cellX; isBoat(gameField, counterX, cellY); counterX -= 2 {
		boatHeadX = counterX
	}
	for counterY = cellY; isBoat(gameField, cellX, counterY); counterY-- {
		boatHeadY = counterY
	}

	// Визначення орієнтації корабля
	if isBoat(gameField, boatHeadX+2, boatHeadY) {
		boatOrientation = false
	} else {
		boatOrientation = true
	}

	// Перевстановлення лічильників
	counterX = boatHeadX
	counterY = boatHeadY
	// Поки поточна клітина є кораблем
	for isBoat(gameField, counterX, counterY) {
		// Змінна станом пошкодженості поточної клітини (true - пошкоджена, false - ціла)
		var boatPartDamaged bool
		// Якщо частина корабля пошкоджена
		if gameField[counterX][counterY].VisibleCell.Style == tcell.StyleDefault.Foreground(tcell.ColorRed) {
			boatPartDamaged = true
		} else {
			boatPartDamaged = false
		}
		// Додавання значення до загального масиву всіх клітин корабля
		boatPartsDamaged = append(boatPartsDamaged, boatPartDamaged)
		// Інкрементація лічильників
		if boatOrientation {
			counterY++
		} else {
			counterX += 2
		}
	}

	// Повернення отриманого корабля
	return Boat{
		headX:        boatHeadX,
		headY:        boatHeadY,
		partsDamaged: boatPartsDamaged,
		orientation:  boatOrientation,
	}
}

// Додавання корабля до ігрової матриці (ігрове поле, корабль)
func placeBoat(gameField [][]GameCell, boat Boat) {
	// Визначення, чи є корабль пошкодженим
	boatDamaged := isDamaged(boat)

	// Лічильники для проходження всієї довжини корабля
	counterX := boat.headX
	counterY := boat.headY
	sizeCounter := 0

	// Поки залишаються неопрацьовані частини корабля
	for sizeCounter < len(boat.partsDamaged) {
		// Встановлення значення клітини корабля
		gameField[counterX][counterY].VisibleCell.Rune = '█'
		gameField[counterX+1][counterY].VisibleCell.Rune = '█'

		// Встановлення кольору клітини корабля відповідно до його пошкодженості
		if boat.partsDamaged[sizeCounter] {
			gameField[counterX][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorRed)
			gameField[counterX+1][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorRed)
		} else if boatDamaged {
			gameField[counterX][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor())
			gameField[counterX+1][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor())
		} else {
			gameField[counterX][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorGreen)
			gameField[counterX+1][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorGreen)
		}

		// Інкрементація лічильників циклу
		if boat.orientation {
			counterY++
		} else {
			counterX = counterX + 2
		}
		sizeCounter++
	}
}

// Видалення корабля з ігрової матриці (ігрове поле, корабль)
func removeBoat(gameField [][]GameCell, boat Boat) {
	// Лічильники для проходження всієї довжини корабля
	counterX := boat.headX
	counterY := boat.headY
	sizeCounter := 0

	// Поки залишаються неопрацьовані частини корабля
	for sizeCounter < len(boat.partsDamaged) {

		gameField[counterX][counterY].VisibleCell = Cell{}
		gameField[counterX+1][counterY].VisibleCell = Cell{}

		// Інкрементація лічильників циклу
		if boat.orientation {
			counterY++
		} else {
			counterX = counterX + 2
		}
		sizeCounter++
	}
}

// Визначення, чи може корабль переміститися наверх (ігрове поле, корабль) (true - може, false - ні)
func canMoveUp(gameField [][]GameCell, boat Boat) bool {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Відкладене повернення початкового корабля на ігрове поле
	defer placeBoat(gameField, boat)
	// Якщо корабль вертикальний або одноклітинний
	if boat.orientation || len(boat.partsDamaged) == 1 {
		// Якщо можна розмістити корабль клітиною вище
		if canPlaceBoat(gameField, Boat{
			headX:        boat.headX,
			headY:        boat.headY - 1,
			partsDamaged: boat.partsDamaged,
			orientation:  boat.orientation,
		}) {
			return true
		}
	}
	return false
}

// Визначення, чи може корабль переміститися вниз (ігрове поле, корабль) (true - може, false - ні)
func canMoveDown(gameField [][]GameCell, boat Boat) bool {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Відкладене повернення початкового корабля на ігрове поле
	defer placeBoat(gameField, boat)
	// Якщо корабль вертикальний або одноклітинний
	if boat.orientation || len(boat.partsDamaged) == 1 {
		// Якщо можна розмістити корабль клітиною нижче
		if canPlaceBoat(gameField, Boat{
			headX:        boat.headX,
			headY:        boat.headY + 1,
			partsDamaged: boat.partsDamaged,
			orientation:  boat.orientation,
		}) {
			return true
		}
	}
	return false
}

// Визначення, чи може корабль переміститися вліво (ігрове поле, корабль) (true - може, false - ні)
func canMoveLeft(gameField [][]GameCell, boat Boat) bool {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Відкладене повернення початкового корабля на ігрове поле
	defer placeBoat(gameField, boat)
	// Якщо корабль горизонтальний або одноклітинний
	if !boat.orientation || len(boat.partsDamaged) == 1 {
		// Якщо можна розмістити корабль клітиною лівіше
		if canPlaceBoat(gameField, Boat{
			headX:        boat.headX - 2,
			headY:        boat.headY,
			partsDamaged: boat.partsDamaged,
			orientation:  boat.orientation,
		}) {
			return true
		}
	}
	return false
}

// Визначення, чи може корабль переміститися вправо (ігрове поле, корабль) (true - може, false - ні)
func canMoveRight(gameField [][]GameCell, boat Boat) bool {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Відкладене повернення початкового корабля на ігрове поле
	defer placeBoat(gameField, boat)
	// Якщо корабль горизонтальний або одноклітинний
	if !boat.orientation || len(boat.partsDamaged) == 1 {
		// Якщо можна розмістити корабль клітиною правіше
		if canPlaceBoat(gameField, Boat{
			headX:        boat.headX + 2,
			headY:        boat.headY,
			partsDamaged: boat.partsDamaged,
			orientation:  boat.orientation,
		}) {
			return true
		}
	}
	return false
}

// Визначення, чи може корабль розвернутися вліво (ігрове поле, корабль) (true - може, false - ні)
func canRotateLeft(gameField [][]GameCell, boat Boat) bool {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Відкладене повернення початкового корабля на ігрове поле
	defer placeBoat(gameField, boat)
	// Якщо корабль не одноклітинний
	if len(boat.partsDamaged) != 1 {
		// Якщо корабль вертикальний
		if boat.orientation {
			// Якщо можна розмістити корабль з іншою орієнтацією
			if canPlaceBoat(gameField, Boat{
				headX:        boat.headX,
				headY:        boat.headY,
				partsDamaged: boat.partsDamaged,
				orientation:  !boat.orientation,
			}) {
				return true
			}
		  // Якщо корабль горизонтальний
		} else {
			// Якщо можна розмістити корабль з іншою орієнтацією 
			if canPlaceBoat(gameField, Boat{
				headX:        boat.headX,
				headY:        boat.headY - (len(boat.partsDamaged) - 1),
				partsDamaged: boat.partsDamaged,
				orientation:  !boat.orientation,
			}) {
				return true
			}
		}
	}
	return false
}

// Визначення, чи може корабль розвернутися вправо (ігрове поле, корабль) (true - може, false - ні)
func canRotateRight(gameField [][]GameCell, boat Boat) bool {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Відкладене повернення початкового корабля на ігрове поле
	defer placeBoat(gameField, boat)
	// Якщо корабль не одноклітинний
	if len(boat.partsDamaged) != 1 {
		// Якщо корабль вертикальний
		if boat.orientation {
			// Якщо можна розмістити корабль з іншою орієнтацією
			if canPlaceBoat(gameField, Boat{
				headX:        boat.headX - (len(boat.partsDamaged)-1)*2,
				headY:        boat.headY,
				partsDamaged: boat.partsDamaged,
				orientation:  !boat.orientation,
			}) {
				return true
			}
		  // Якщо корабль горизонтальний
		} else {
			// Якщо можна розмістити корабль з іншою орієнтацією
			if canPlaceBoat(gameField, Boat{
				headX:        boat.headX,
				headY:        boat.headY,
				partsDamaged: boat.partsDamaged,
				orientation:  !boat.orientation,
			}) {
				return true
			}
		}
	}
	return false
}

// Переміщення корабля наверх (ігрове поле, корабль) (корабль)
func moveUp(gameField [][]GameCell, boat Boat) Boat {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Новий переміщений корабль
	movedBoat := Boat{
		headX:        boat.headX,
		headY:        boat.headY - 1,
		partsDamaged: boat.partsDamaged,
		orientation:  boat.orientation,
	}
	// Розміщення нового корабля
	placeBoat(gameField, movedBoat)
	// Повернення нового корабля
	return movedBoat
}

// Переміщення корабля вниз (ігрове поле, корабль) (корабль)
func moveDown(gameField [][]GameCell, boat Boat) Boat {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Новий переміщений корабль
	movedBoat := Boat{
		headX:        boat.headX,
		headY:        boat.headY + 1,
		partsDamaged: boat.partsDamaged,
		orientation:  boat.orientation,
	}
	// Розміщення нового корабля
	placeBoat(gameField, movedBoat)
	// Повернення нового корабля
	return movedBoat
}

// Переміщення корабля вліво (ігрове поле, корабль) (корабль)
func moveLeft(gameField [][]GameCell, boat Boat) Boat {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Новий переміщений корабль
	movedBoat := Boat{
		headX:        boat.headX - 2,
		headY:        boat.headY,
		partsDamaged: boat.partsDamaged,
		orientation:  boat.orientation,
	}
	// Розміщення нового корабля
	placeBoat(gameField, movedBoat)
	// Повернення нового корабля
	return movedBoat
}

// Переміщення корабля вправо (ігрове поле, корабль) (корабль)
func moveRight(gameField [][]GameCell, boat Boat) Boat {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Новий переміщений корабль
	movedBoat := Boat{
		headX:        boat.headX + 2,
		headY:        boat.headY,
		partsDamaged: boat.partsDamaged,
		orientation:  boat.orientation,
	}
	// Розміщення нового корабля
	placeBoat(gameField, movedBoat)
	// Повернення нового корабля
	return movedBoat
}

// Поворот корабля вліво (ігрове поле, корабль) (корабль)
func rotateLeft(gameField [][]GameCell, boat Boat) Boat {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Новий переміщений корабль
	var movedBoat Boat
	// Якщо корабль був вертикальний
	if boat.orientation {
		movedBoat = Boat{
			headX:        boat.headX,
			headY:        boat.headY,
			partsDamaged: boat.partsDamaged,
			orientation:  !boat.orientation,
		}
	  // Якщо корабль був горизонтальний
	} else {
		movedBoat = Boat{
			headX:        boat.headX,
			headY:        boat.headY - (len(boat.partsDamaged) - 1),
			partsDamaged: boat.partsDamaged,
			orientation:  !boat.orientation,
		}
	}
	// Розміщення нового корабля
	placeBoat(gameField, movedBoat)
	// Повернення нового корабля
	return movedBoat
}

// Поворот корабля вправо (ігрове поле, корабль) (корабль)
func rotateRight(gameField [][]GameCell, boat Boat) Boat {
	// Видалення початкового корабля з ігрового поля
	removeBoat(gameField, boat)
	// Новий переміщений корабль
	var movedBoat Boat
	// Якщо корабль був вертикальний
	if boat.orientation {
		movedBoat = Boat{
			headX:        boat.headX - (len(boat.partsDamaged)-1)*2,
			headY:        boat.headY,
			partsDamaged: boat.partsDamaged,
			orientation:  !boat.orientation,
		}
	  // Якщо корабль був горизонтальний
	} else {
		movedBoat = Boat{
			headX:        boat.headX,
			headY:        boat.headY,
			partsDamaged: boat.partsDamaged,
			orientation:  !boat.orientation,
		}
	}
	// Розміщення нового корабля
	placeBoat(gameField, movedBoat)
	// Повернення нового корабля
	return movedBoat
}

// Виділення корабля на полі (ігрове поле, корабль)
func selectBoat(gameField [][]GameCell, boat Boat) {
	// Якщо можливе виконання руху, позначити напрямок на полі
	if canMoveUp(gameField, boat) {
		gameField[boat.headX][boat.headY-1].VisibleCell.Rune = '↑'
	}
	if canMoveDown(gameField, boat) {
		gameField[boat.headX+1][boat.headY+len(boat.partsDamaged)].VisibleCell.Rune = '↓'
	}
	if canMoveLeft(gameField, boat) {
		gameField[boat.headX-1][boat.headY].VisibleCell.Rune = '←'
	}
	if canMoveRight(gameField, boat) {
		gameField[boat.headX+len(boat.partsDamaged)*2][boat.headY].VisibleCell.Rune = '→'
	}
	if canRotateLeft(gameField, boat) {
		if boat.orientation {
			gameField[boat.headX+2][boat.headY].VisibleCell.Rune = '↺'
		} else {
			gameField[boat.headX][boat.headY-1].VisibleCell.Rune = '↺'
		}
	}
	if canRotateRight(gameField, boat) {
		if boat.orientation {
			gameField[boat.headX-1][boat.headY].VisibleCell.Rune = '↻'
		} else {
			gameField[boat.headX][boat.headY+1].VisibleCell.Rune = '↻'
		}
	}
}

// Виконання пострілу кораблем (ігрове поле, координата х пострілу, координата у пострілу) (true - влучання, false - промах)
func fire(gameField [][]GameCell, cellX int, cellY int) bool {
	// Якщо влучення в ціль
	if gameField[cellX][cellY].VisibleCell.Rune == '█' {
		// Фарбування клітини в яку потрапив постріл
		gameField[cellX][cellY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorRed)
		gameField[cellX+1][cellY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorRed)
		// Отримання структури корабля після поранення
		firedBoat := getBoatByCell(gameField, cellX, cellY)

		// Лічильники для проходження всієї довжини корабля
		counterX := firedBoat.headX
		counterY := firedBoat.headY
		sizeCounter := 0

		// Поки залишаються неопрацьовані частини корабля
		for sizeCounter < len(firedBoat.partsDamaged) {
			// Якщо корабль був непошкоджений
			if gameField[counterX][counterY].VisibleCell.Style == tcell.StyleDefault.Foreground(tcell.ColorGreen) {
				// Встановлення кольору видимої клітини корабля
				gameField[counterX][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor())
				gameField[counterX+1][counterY].VisibleCell.Style = tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor())
			}

			// Інкрементація лічильників циклу
			if firedBoat.orientation {
				counterY++
			} else {
				counterX = counterX + 2
			}
			sizeCounter++
		}

		// Якщо пострілом корабль було знищено
		if isDestroyed(firedBoat) {
			// Лічильники для проходження всієї довжини корабля
			counterX = firedBoat.headX
			counterY = firedBoat.headY
			sizeCounter := 0

			// Поки залишаються неопрацьовані частини корабля
			for sizeCounter < len(firedBoat.partsDamaged) {
				// Встановлення значення та кольору навколишніх клітин корабля
				for i := counterX - 2; i <= counterX+2; i += 2 {
					for j := counterY - 1; j <= counterY+1; j++ {
						if inGameField(i, j) {
							if !isBoat(gameField, i, j) {
								hiddenCell := Cell{
									Rune:  '░',
									Style: tcell.StyleDefault.Foreground(tcell.ColorBlue),
								}
								gameField[i][j].HiddenCell = hiddenCell
								gameField[i+1][j].HiddenCell = hiddenCell
							}
						}
					}
				}
				// Встановлення значення та кольору прихованої клітини корабля як знищеної
				hiddenCell := Cell{
					Rune:  '▓',
					Style: tcell.StyleDefault.Foreground(tcell.ColorRed),
				}
				gameField[counterX][counterY].HiddenCell = hiddenCell
				gameField[counterX+1][counterY].HiddenCell = hiddenCell

				// Інкрементація лічильників циклу
				if firedBoat.orientation {
					counterY++
				} else {
					counterX = counterX + 2
				}
				sizeCounter++
			}
		  // Якщо корабль було поранено
		} else {
			// Встановлення значення та кольору прихованої клітини корабля як пошкодженої
			hiddenCell := Cell{
				Rune:  '░',
				Style: tcell.StyleDefault.Foreground(tcell.ColorLightPink.TrueColor()),
			}
			gameField[cellX][cellY].HiddenCell = hiddenCell
			gameField[cellX+1][cellY].HiddenCell = hiddenCell
		}
		// Повернення значення влучання
		return true
	  // Якщо промах	
	} else {
		// Встановлення значення та кольору прихованої клітини як незайнятої
		hiddenCell := Cell{
			Rune:  '░',
			Style: tcell.StyleDefault.Foreground(tcell.ColorBlue),
		}
		gameField[cellX][cellY].HiddenCell = hiddenCell
		gameField[cellX+1][cellY].HiddenCell = hiddenCell
		// Повернення значення промаху
		return false
	}
}

// Прибирання виділення з корабля (ігрове поле, корабль)
func unselectBoat(gameField [][]GameCell, boat Boat) {
	// Координати хвоста (останньої клітини) корабля
	boatTailX := boat.headX
	boatTailY := boat.headY
	// Встановлення координат хвоста (останньої клітини) корабля
	if boat.orientation {
		boatTailY += len(boat.partsDamaged) - 1
	} else {
		boatTailX += (len(boat.partsDamaged) - 1) * 2
	}
	// Проходження через всі навколишні клітини корабля
	for x := boat.headX - 1; x <= boatTailX+2; x++ {
		for y := boat.headY - 1; y <= boatTailY+1; y++ {
			if inGameField(x, y) {
				if !isBoat(gameField, x, y) {
					gameField[x][y].VisibleCell = Cell{}
				}
			}
		}
	}
}

// Створення випадкового ігрового поля (кількість кораблів) (ігрове поле)
func createRandomGameField(boatsAmount int) [][]GameCell {
	// Ініціалізація ігрової матриці
	gameField := make([][]GameCell, battlefieldSize*2)
	for i := 0; i < battlefieldSize*2; i++ {
		gameField[i] = make([]GameCell, battlefieldSize)
	}

	// Змінні для поступового заповнення ігрового поля
	placedBoats := 0
	boatSize := 1
	boatMaxSize := 4

	// Якщо залишаються нерозміщені кораблі
	for placedBoats < boatsAmount {
		// Випадковий вибір початкової клітини та орієнтації корабля
		randX := rand.Intn(battlefieldSize * 2)
		randY := rand.Intn(battlefieldSize)
		randBool := rand.Intn(2) == 1

		// Створення випадкового корабля
		boat := Boat{
			headX:        randX - randX%2,
			headY:        randY,
			partsDamaged: make([]bool, boatSize),
			orientation:  randBool,
		}

		if canPlaceBoat(gameField, boat) {
			// Розташування корабля у перевірено вільному просторі
			placeBoat(gameField, boat)
			placedBoats++
			// Встановлення нової розмірності для наступного корабля (1=>2=>3=>4=>1=>2=>3=>1=>2=>1)
			if boatSize == boatMaxSize {
				boatMaxSize--
				boatSize = 1
			} else {
				boatSize++
			}
		}
	}

	// Повернення створеного ігрового поля
	return gameField
}

// Друк вмісту ігрового поля (екран, стан відкритості (0 - повністю приховане, 1 - обмежено відкрите, 2 - повністю відкрите), стан динаміки (0 - відсутність динаміки, 1 - динамічний друк з верхнього лівого кута, 2 - динамічний друк з нижнього правого кута), початковий стовпець, початковий рядок, ігрове поле)
func printGameField(screen tcell.Screen, openState int, dynamicState int, startCol int, startRow int, gameField [][]GameCell) {
	// Послідовний друк збережених значень комірок консолі починаючи з кута
	for k := 0; k < battlefieldSize*2-1; k++ {
		for i := 0; i <= k; i++ {
			j := k - i
			// Якщо дана ітерація в межах ігрового поля
			if i < battlefieldSize && j < battlefieldSize {
				// Поточний стовпець
				var c int
				// Поточний рядок
				var r int
				if dynamicState == 2 {
					// Відлік з кінця
					c = (battlefieldSize - 1 - i) * 2
					r = battlefieldSize - 1 - j
				} else {
					// Відлік з початку
					c = i * 2
					r = j
				}
				// Ліва частина клітини
				var cellL Cell
				// Права частина клітини
				var cellR Cell
				// Вибір в залежності від відкритості поля
				switch openState {
				case 0:
					cellL = Cell{Rune: '▒', Style: tcell.StyleDefault.Foreground(tcell.ColorGray)}
					cellR = Cell{Rune: '▒', Style: tcell.StyleDefault.Foreground(tcell.ColorGray)}
				case 1:
					cellL = gameField[c][r].HiddenCell
					cellR = gameField[c+1][r].HiddenCell
				case 2:
					cellL = gameField[c][r].VisibleCell
					cellR = gameField[c+1][r].VisibleCell
					// Якщо поточна клітина не є клітиною корабля
					if cellL.Rune != '█' {
						// Якщо прихована клітина має якесь значення
						if gameField[c][r].HiddenCell.Rune != 0 && gameField[c+1][r].HiddenCell.Rune != 0 {
							cellL.Style = tcell.StyleDefault.Foreground(tcell.ColorBlue)
							cellR.Style = tcell.StyleDefault.Foreground(tcell.ColorBlue)
							if cellL.Rune == 0 && cellR.Rune == 0 {
								cellL.Rune = '░'
								cellR.Rune = '░'
							}
						}
					}
				}
				// Одночасний друк лівої та правої частини клітини
				screen.SetContent(startCol+1+c, startRow+1+r, cellL.Rune, nil, cellL.Style)
				screen.SetContent(startCol+1+c+1, startRow+1+r, cellR.Rune, nil, cellR.Style)
			}
		}
		if dynamicState != 0 {
			// Невелика затримка для забезпечення динамічної послідовності друку
			screen.Show()
			time.Sleep(time.Millisecond * 40)
		}
	}
}

// Рух корабля за можливістю (екран, режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ), початковий стовпець, початковий рядок, ігрове поле, можливий рух)
func moveIfPossible(screen tcell.Screen, gameModeState int, startCol int, startRow int, gameField [][]GameCell,
	possibleMovement PossibleMovement) {
	// Якщо можна здійснити дію та корабль ще не був переміщений
	if possibleMovement.canMove(gameField, selectedBoat) && !boatMoved {
		unselectBoat(gameField, selectedBoat)
		selectedBoat = possibleMovement.move(gameField, selectedBoat)
		// Якщо етап гри має дробову частину 0.5 (етап налаштування кораблів)
		if float64(gameState)-0.5 == math.Trunc(float64(gameState)) {
			selectBoat(gameField, selectedBoat)
			printGameField(screen, 2, 0, startCol, startRow, gameField)
		} else {
			if isDamaged(selectedBoat) {
				printGameField(screen, 2, 0, startCol, startRow, gameField)
				makeMove(screen, gameModeState)
			} else {
				printGameField(screen, 2, 0, startCol, startRow, gameField)
				boatMoved = true
			}
		}
	}
}

// Подія натискання символьної клавіші (екран, режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ), символ)
func keyRuneEvent(screen tcell.Screen, gameModeState int, r rune) {
	// Якщо не обрано корабль або гра вже закінчилася
	if selectedBoat.partsDamaged == nil || gameState == 3 {
		return
	}

	// Змінні для позначення поточного ігрового поля та його координат
	var gameField [][]GameCell
	var startCol int
	var startRow int

	if playerToMove {
		gameField = player1GameField
		startCol = startCol1
		startRow = startRow1
		// Якщо натиснута клавіша не є літерою (не належить до команд управління першого гравця)
		if !unicode.IsLetter(r) {
			return
		}
	} else {
		gameField = player2GameField
		startCol = startCol2
		startRow = startRow2
		// Якщо натиснута клавіша не є цифрою (не належить до команд управління другого гравця)
		if !unicode.IsDigit(r) {
			return
		}
	}

	// Вибір нижнього регістру символу з клавіші
	switch unicode.ToLower(r) {
	case 'w', 'ц', '8':
		// Рух наверх
		moveIfPossible(screen, gameModeState, startCol, startRow, gameField, possibleMovements[0])
	case 's', 'ы', 'і', '5':
		// Рух вниз
		moveIfPossible(screen, gameModeState, startCol, startRow, gameField, possibleMovements[1])
	case 'a', 'ф', '4':
		// Рух вліво
		moveIfPossible(screen, gameModeState, startCol, startRow, gameField, possibleMovements[2])
	case 'd', 'в', '6':
		// Рух вправо
		moveIfPossible(screen, gameModeState, startCol, startRow, gameField, possibleMovements[3])
	case 'q', 'й', '7':
		// Поворот вліво
		moveIfPossible(screen, gameModeState, startCol, startRow, gameField, possibleMovements[4])
	case 'e', 'у', '9':
		// Поворот вправо
		moveIfPossible(screen, gameModeState, startCol, startRow, gameField, possibleMovements[5])
	}
}

// Подія натискання лівої кнопки миші (екран, режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ), позиція по х, позиція по у)
func mouseClickEvent(screen tcell.Screen, gameModeState int, x int, y int) {
	if gameState == 3 {
		return
	}

	// Змінні для позначення ігрових полів та їх координат для гравців - атакуючої та обороняючої сторони
	var playerGameField [][]GameCell
	var opponentGameField [][]GameCell
	var playerStartCol int
	var playerStartRow int
	var opponentStartCol int
	var opponentStartRow int

	if playerToMove {
		playerGameField = player1GameField
		opponentGameField = player2GameField
		playerStartCol = startCol1
		playerStartRow = startRow1
		opponentStartCol = startCol2
		opponentStartRow = startRow2
	} else {
		playerGameField = player2GameField
		opponentGameField = player1GameField
		playerStartCol = startCol2
		playerStartRow = startRow2
		opponentStartCol = startCol1
		opponentStartRow = startRow1
	}

	// Якщо етап гри нецілочисельний (дійовий)
	if float64(gameState) != math.Trunc(float64(gameState)) {
		if selectedBoat.partsDamaged != nil {
			unselectBoat(playerGameField, selectedBoat)
			printGameField(screen, 2, 0, playerStartCol, playerStartRow, playerGameField)
		}

		// Якщо натиснута позиція в межах поля гравця що ходить
		if inGameField(x-1-playerStartCol, y-1-playerStartRow) {
			if !boatMoved {
				// Встановлення нормалізованих координат клітини ігрового поля гравця за позицією миші в консолі
				boatX := (x - 1 - playerStartCol) - (x-1-playerStartCol)%2
				boatY := y - playerStartRow - 1
				// Якщо клітина корабля що у грі
				if isBoat(playerGameField, boatX, boatY) && !isDestroyed(getBoatByCell(playerGameField, boatX, boatY)) {
					selectedBoat = getBoatByCell(playerGameField, boatX, boatY)
					selectBoat(playerGameField, selectedBoat)
					printGameField(screen, 2, 0, playerStartCol, playerStartRow, playerGameField)
				}
			}
		  // Якщо натиснута позиція в межах поля суперника
		} else if inGameField(x-1-opponentStartCol, y-1-opponentStartRow) && selectedBoat.partsDamaged != nil {
			// Якщо ігровий процес було розпочато
			if gameState > 2 {
				// Встановлення нормалізованих координат клітини ігрового поля суперника за позицією миші в консолі
				cellX := (x - 1 - opponentStartCol) - (x - 1 - opponentStartCol)%2
				cellY := y - 1 - opponentStartRow
				// Постріл, якщо влучний
				if fire(opponentGameField, cellX, cellY) {
					printGameField(screen, 1, 0, opponentStartCol, opponentStartRow, opponentGameField)
					unselectBoat(playerGameField, selectedBoat)
					printGameField(screen, 2, 0, playerStartCol, playerStartRow, playerGameField)
					boatMoved = true
					// Якщо добитий останній корабль
					if isFieldCleared(opponentGameField) {
						gameState = 3
						updateInterface(screen)
					}
				  // Якщо промах
				} else {
					printGameField(screen, 1, 0, opponentStartCol, opponentStartRow, opponentGameField)
					makeMove(screen, gameModeState)
				}
			}
		}
	}
}

// Здійснення ходу ШІ (екран)
func AIMove(screen tcell.Screen) {
	// Безперервний цикл до виконання умови виходу
	for true {
		// Вибір 2 випадкових значень в межах поля
		randX := rand.Intn(battlefieldSize * 2)
		randY := rand.Intn(battlefieldSize)
		// Нормалізація у координати ігрового поля
		boatX := randX - randX%2
		boatY := randY
		if isBoat(player2GameField, boatX, boatY) {
			// Випадковий корабль за випадковими координатами
			randBoat := getBoatByCell(player2GameField, boatX, boatY)
			if !isDestroyed(randBoat) {
				// Якщо існують пересувні кораблі на ігровому полі та корабль ще не було переміщено
				if areMoveableBoatsExist(player2GameField) && !boatMoved {
					if isMoveable(player2GameField, randBoat) {
						// Якщо корабль непошкоджений, або виконання руху замість пострілу з ймовірністю 50% для пошкодженого корабля
						if !isDamaged(randBoat) || rand.Intn(2) == 1 {
							// Безперервний цикл до виконання умови виходу
							for true {
								// Випадковий індекс для масиву можливих рухів
								randI := rand.Intn(6)
								// Якщо рух можна здійснити
								if possibleMovements[randI].canMove(player2GameField, randBoat) {
									// Затримка на виконання ходу ШІ
									screen.Show()
									time.Sleep(time.Second)
									// Рух корабля
									possibleMovements[randI].move(player2GameField, randBoat)
									boatMoved = true
									// Якщо переміщений корабль був пошкодженим - завершення ходу, вихід з функції
									if isDamaged(randBoat) {
										gameState = 2
										return
									  // Якщо переміщений корабль цілий, продовження виконання, перехід до пострілу
									} else {
										break
									}
								}
							}
						}
					  // Якщо корабль нерухомий - продовження пошуку пересувного
					} else {
						continue
					}
				}

				// Безперервний цикл до виконання умови виходу
				for true {
					// Вибір 2 випадкових значень в межах поля
					randX2 := rand.Intn(battlefieldSize * 2)
					randY2 := rand.Intn(battlefieldSize)
					// Нормалізація у координати ігрового поля
					boatX2 := randX2 - randX2%2
					boatY2 := randY2
					// Якщо постріл не у вже знищений корабль, та не в клітину що вже була поцілена, якщо залищаються неперевірені
					if player1GameField[boatX2][boatY2].HiddenCell.Style != tcell.StyleDefault.Foreground(tcell.ColorRed) &&
						(player1GameField[boatX2][boatY2].HiddenCell.Style == tcell.StyleDefault || !areUnfiredCellsExist(player1GameField)) {
						printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
						// Затримка на виконання ходу ШІ
						screen.Show()
						time.Sleep(time.Second)
						// Постріл, якщо влучний
						if fire(player1GameField, boatX2, boatY2) {
							// Якщо добитий останній корабль
							if isFieldCleared(player1GameField) {
								gameState = 3
								updateInterface(screen)
								return
							  // Якщо ще залишаються кораблі в грі - продовження вогню
							} else {
								continue
							}
						  // Якщо промах
						} else {
							// Завершення ходу, вихід з функції
							gameState = 2
							return
						}
					  // Якщо в клітині знищений корабль, або вона вже була поцілена, коли залишаються неперевірені - продовження пошуку
					} else {
						continue
					}
				}
			}
		}
	}
}

// Виконання чергового ходу (екран, режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ))
func makeMove(screen tcell.Screen, gameModeState int) {
	boatMoved = false
	// Вибір поточного етапу гри, перехід до наступного
	switch gameState {
	case 0:
		gameState = 0.5
		if gameModeState == 2 {
			printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
		} else {
			printGameField(screen, 2, 1, startCol1, startRow1, player1GameField)
		}
	case 0.5:
		unselectBoat(player1GameField, selectedBoat)
		selectedBoat = Boat{}
		printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
		if gameModeState == 2 {
			gameState = 2.1
			printGameField(screen, 1, 1, startCol2, startRow2, player2GameField)
			playerToMove = true
		} else {
			gameState = 1
			printGameField(screen, 0, 2, startCol1, startRow1, player1GameField)
			playerToMove = false
		}
	case 1:
		gameState = 1.5
		printGameField(screen, 2, 2, startCol2, startRow2, player2GameField)
	case 1.5:
		unselectBoat(player2GameField, selectedBoat)
		selectedBoat = Boat{}
		printGameField(screen, 2, 0, startCol2, startRow2, player2GameField)
		gameState = 2
		printGameField(screen, 0, 1, startCol2, startRow2, player2GameField)
		// Випадковий вибір гравця що ходить першим у грі
		playerToMove = rand.Intn(2) == 1
	case 2:
		if gameModeState == 2 {
			printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
			printGameField(screen, 1, 0, startCol2, startRow2, player2GameField)
			if playerToMove {
				gameState = 2.1
			} else {
				gameState = 2.3
				// Підтвердження переходу ходу за ШІ
				makeMove(screen, gameModeState)
			}
		} else {
			if playerToMove {
				gameState = 2.1
				printGameField(screen, 2, 1, startCol1, startRow1, player1GameField)
				printGameField(screen, 1, 1, startCol2, startRow2, player2GameField)
			} else {
				gameState = 2.2
				printGameField(screen, 2, 2, startCol2, startRow2, player2GameField)
				printGameField(screen, 1, 2, startCol1, startRow1, player1GameField)
			}
		}
	case 2.1:
		unselectBoat(player1GameField, selectedBoat)
		selectedBoat = Boat{}
		printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
		if gameModeState == 2 {
			gameState = 2.3
			playerToMove = !playerToMove
			updateInterface(screen)
			// Підтвердження переходу ходу за ШІ
			makeMove(screen, gameModeState)
		} else {
			gameState = 2
			printGameField(screen, 0, 2, startCol2, startRow2, player2GameField)
			printGameField(screen, 0, 2, startCol1, startRow1, player1GameField)
			playerToMove = !playerToMove
		}
	case 2.2:
		unselectBoat(player2GameField, selectedBoat)
		selectedBoat = Boat{}
		printGameField(screen, 2, 0, startCol2, startRow2, player2GameField)
		gameState = 2
		printGameField(screen, 0, 1, startCol1, startRow1, player1GameField)
		printGameField(screen, 0, 1, startCol2, startRow2, player2GameField)
		playerToMove = !playerToMove
	case 2.3:
		AIMove(screen)
		printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
		if gameState != 3 {
			playerToMove = !playerToMove
			updateInterface(screen)
			// Підтвердження переходу ходу за ШІ
			makeMove(screen, gameModeState)
		}
	case 3:
		// Очищення попереднього екрану та повернення в меню
		gameBegin = false
		eraseScreenDynamically(screen)
		beginMenu(screen)
	}

	updateInterface(screen)
}

// Початок гри
func beginPlay(screen tcell.Screen, gameModeState int, gameDataState int, boatsAmount int) {
	// Підтвердження виходу з гри (true - підтверджено, false - не підтверджено)
	confirmExitState := false

	// Назва файлу збереження відповідно до режиму гри
	var filename string
	switch gameModeState {
	case 1:
		filename = "save1.txt"
	case 2:
		filename = "save2.txt"
	}

	printAuthor(screen)

	setDefaultExitState(screen)

	printGameFieldBorder(screen, true)
	printGameFieldBorder(screen, false)

	// Створення нової гри
	if gameDataState == 1 {
		createSaveFile(filename)
		player1GameField = createRandomGameField(boatsAmount)
		player2GameField = createRandomGameField(boatsAmount)
		playerToMove = true
		boatMoved = false
		gameState = 0
	  // Завантаження існуючої гри
	} else {
		saveData := readSaveFile(filename)
		player1GameField = saveData.player1GameField
		player2GameField = saveData.player2GameField
		playerToMove = saveData.playerToMove
		boatMoved = saveData.boatMoved
		gameState = 2
	}

	printGameField(screen, 0, 0, startCol1, startRow1, player1GameField)
	printGameField(screen, 0, 0, startCol2, startRow2, player2GameField)

	updateInterface(screen)

	// Режим гри проти ШІ
	if gameModeState == 2 {
		makeMove(screen, gameModeState)
	}

	drawScreenDynamically(screen, false)

	// Обробник подій
	for {
		// Якщо гра закінчилася
		if !gameBegin {
			break
		}

		// Очікування подій
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		// Натискання клавіші
		case *tcell.EventKey:
			// Вибір клавіші залежно від її типу
			switch ev.Key() {
			// Звичайна символьна клавіша
			case tcell.KeyRune:
				// Клавіша Space
				if ev.Rune() == ' ' {
					makeMove(screen, gameModeState)
				  // Клавіша управління
				} else {
					keyRuneEvent(screen, gameModeState, ev.Rune())
				}

			// Клавіша Esc
			case tcell.KeyEscape:
				if confirmExitState {
					confirmExitState = false
					setDefaultExitState(screen)
				} else {
					confirmExitState = true
					dataSaved := gameState != 3
					setConfirmExitState(screen, dataSaved)
				}
			// Клавіша Enter
			case tcell.KeyEnter:
				if confirmExitState {
					// Якщо етап гри нецілочисельний (дійовий)
					if float64(gameState) != math.Trunc(float64(gameState)) {
						if playerToMove {
							unselectBoat(player1GameField, selectedBoat)
							selectedBoat = Boat{}
							printGameField(screen, 2, 0, startCol1, startRow1, player1GameField)
						} else {
							unselectBoat(player2GameField, selectedBoat)
							selectedBoat = Boat{}
							printGameField(screen, 2, 0, startCol2, startRow2, player2GameField)
						}
					}
					// Якщо власне ігровий процес самої гри розпочато, але ще не завершено
					if gameState >= 2 && gameState < 3 {
						writeSaveFile(filename, SaveData{
							player1GameField,
							player2GameField,
							playerToMove,
							boatMoved,
						})
					}

					// Очищення попереднього екрану та повернення в меню
					gameBegin = false
					eraseScreenDynamically(screen)
					beginMenu(screen)
				}
			}
		// Рух або натисканні миші
		case *tcell.EventMouse:
			// Отримання натиснутої кнопки та позиції миші
			button := ev.Buttons()
			x, y := ev.Position()

			// Якщо натиснута ліва кнопка
			if button&tcell.Button1 != 0 {
				mouseClickEvent(screen, gameModeState, x, y)
			}
		}

		// Оновлення екрану
		screen.Show()
	}
}