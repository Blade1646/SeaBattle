	// Головний пакет програми
	package main

	// Підключення необхідних бібліотек
	import (
		"strings"
		"strconv"

		"github.com/gdamore/tcell/v2"
	)

	// Очищення головного екрану меню гри
	func clearMainMenuScreen(screen tcell.Screen) {
		for i := 3; i <= screenSizeY-5; i++ {
			for j := 1; j < screenSizeX-2; j++ {
				screen.SetContent(j, i, ' ', nil, tcell.StyleDefault)
			}
		}
	}

	// Встановлення опції режиму гри (екран, режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ))
	func setGameModeState(screen tcell.Screen, state int) {
		clearMainMenuScreen(screen)

		// Друк тексту підказки
		printStringLeft(screen, 3, "Please, choose a game mode to play (tick [x] one):")

		// Вибір режиму відповідно до вхідного параметру
		switch state {
		case 0:
			printStringLeft(screen, 5, "   [ ]   vs player")
			printStringLeft(screen, 6, "   [ ]   vs AI")
		case 1:
			printStringLeft(screen, 5, "   [x]   vs player")
			printStringLeft(screen, 6, "   [ ]   vs AI")
			paintCellsRowWithColor(screen, 5, 5, 1, tcell.ColorRed)
			setGameDataState(screen, 0, 1)
		case 2:
			printStringLeft(screen, 5, "   [ ]   vs player")
			printStringLeft(screen, 6, "   [x]   vs AI")
			paintCellsRowWithColor(screen, 5, 6, 1, tcell.ColorRed)
			setGameDataState(screen, 0, 2)
		}
	}

	// Встановлення опції збереження гри (екран, збереження гри (0 - не обрано, 1 - створення нової гри, 2 - завантаження існуючої гри), режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ))
	func setGameDataState(screen tcell.Screen, state int, gameModeState int) {
		// Назва файлу збереження відповідно до режиму гри
		var filename string
		switch gameModeState {
		case 1:
			filename = "save1.txt"
		case 2:
			filename = "save2.txt"
		}

		// Друк тексту підказки
		printStringLeft(screen, 8, "Start new or continue existing game? (tick [x] one):")

		// Вибір збереження відповідно до вхідного параметру
		switch state {
		case 0:
			printStringLeft(screen, 10, "   [ ]   New game")
			printStringLeft(screen, 11, "   [ ]   Load game")
		case 1:
			printStringLeft(screen, 10, "   [x]   New game")
			printStringLeft(screen, 11, "   [ ]   Load game")
			paintCellsRowWithColor(screen, 5, 10, 1, tcell.ColorRed)
			setGameBoatsAmount(screen, 0, 0, 0)
		case 2:
			if !saveExists(filename) {
				return
			}
			printStringLeft(screen, 10, "   [ ]   New game")
			printStringLeft(screen, 11, "   [x]   Load game")
			paintCellsRowWithColor(screen, 5, 11, 1, tcell.ColorRed)

			// Очищення попереднього екрану та початок гри відповідно до налаштувань
			gameBegin = true
			eraseScreenDynamically(screen)
			beginPlay(screen, gameModeState, state, 0)
		}

		// Сірий кольор недоступної опції якщо відсутнє наявне збереження гри
		if !saveExists(filename) {
			paintCellsRowWithColor(screen, 4, 11, 15, tcell.ColorGray)
		}
	}

	// Встановлення опції кількості кораблів (екран, кількість кораблів, режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ), збереження гри (0 - не обрано, 1 - створення нової гри, 2 - завантаження існуючої гри))
	func setGameBoatsAmount(screen tcell.Screen, amount int, gameModeState int, gameDataState int) {

		// Текст підказки
		boatsAmountQuestionString := "Please, choose an amount of boats for each side (select [xx] one):"
		boatsAmountOptionsString := " 4    5    6    7    8    9    10 "

		// Приведення рядків до однієї довжини
		diff := len(boatsAmountQuestionString) - len(boatsAmountOptionsString)
		boatsAmountOptionsStringLong := strings.Repeat(" ", diff/2) + boatsAmountOptionsString + strings.Repeat(" ", diff-diff/2)

		// Друк тексту підказки
		printStringLeft(screen, 13, boatsAmountQuestionString)
		printStringLeft(screen, 15, boatsAmountOptionsStringLong)

		// Якщо користувач обрав одну з перелічених кількостей
		if amount > 0 {

			// Отримання індексу обраного значення
			amountIndex := getStringIndex(screen, 15, strconv.Itoa(amount))

			// Виділення обраного значення
			screen.SetContent(amountIndex-1, 15, '[', nil, tcell.StyleDefault)
			screen.SetContent(amountIndex+len(strconv.Itoa(amount)), 15, ']', nil, tcell.StyleDefault)

			paintCellsRowWithColor(screen, amountIndex, 15, len(strconv.Itoa(amount)), tcell.ColorRed)

			// Очищення попереднього екрану та початок гри відповідно до налаштувань
			gameBegin = true
			eraseScreenDynamically(screen)
			beginPlay(screen, gameModeState, gameDataState, amount)
		}
	}

// Початок роботи з меню гри (екран)
func beginMenu(screen tcell.Screen) {
	// Підтвердження виходу з гри (true - підтверджено, false - не підтверджено)
	confirmExitState := false
	// Режим гри (0 - не обрано, 1 - проти гравця, 2 - проти ШІ)
	gameModeState := 0
	// Дані гри (0 - не обрано, 1 - створення нової гри, 2 - завантаження існуючої гри)
	gameDataState := 0

	// Друк заголовку гри
	printStringMiddle(screen, 1, "Welcome to the game \"Admiral\"!")
	paintCellsRowWithColor(screen, getStringIndex(screen, 1, "Admiral"), 1, len("Admiral"), tcell.ColorBlue) 

	printAuthor(screen)

	setDefaultExitState(screen)
	setGameModeState(screen, 0)

	drawScreenDynamically(screen, true)

	// Обробник подій
	for {
		// Якщо гра вже почалася
		if gameBegin {
			break
		}

		// Очікування подій
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		// Натискання клавіші
		case *tcell.EventKey:
			// Клавіша Esc
			if ev.Key() == tcell.KeyEscape {
				if confirmExitState {
					confirmExitState = false
					setDefaultExitState(screen)
				} else {
					confirmExitState = true
					setConfirmExitState(screen, false)
				}
			}
			// Клавіша Enter
			if ev.Key() == tcell.KeyEnter {
				if confirmExitState {
					return
				}
			}

		// Рух або натисканні миші
		case *tcell.EventMouse:
			// Отримання натиснутої кнопки та позиції миші
			button := ev.Buttons()
			x, y := ev.Position()

			// Якщо натиснута ліва кнопка
			if button&tcell.Button1 != 0 {
				// Якщо стовпець опцій меню
				if x == 5 {

					// Отримання лівого та правого символу від натиснутої позиції миші
					leftR, _, _, _ := screen.GetContent(x-1, y)
					rightR, _, _, _ := screen.GetContent(x+1, y)

					// Якщо пункт меню
					if leftR == '[' && rightR == ']' {

						// Отримання всього рядка обраної опції
						rowString := ""
						for col := 0; col < screenSizeX; col++ {
							r, _, _, _ := screen.GetContent(col, y)
							rowString += string(r)
						}

						// Визначення обраного пункту меню

						// Обрання режиму гри проти гравця
						if strings.Contains(strings.ToLower(rowString), "player") && gameModeState != 1 {
							gameModeState = 1
							gameDataState = 0
							setGameModeState(screen, 1)
						  // Обрання режиму гри проти ШІ
						} else if strings.Contains(strings.ToLower(rowString), "ai") && gameModeState != 2 {
							gameModeState = 2
							gameDataState = 0
							setGameModeState(screen, 2)
						  // Обрання даних гри - створення нової
						} else if strings.Contains(strings.ToLower(rowString), "new") && gameDataState != 1 {
							gameDataState = 1
							setGameDataState(screen, 1, gameModeState)
						  // Обрання даних гри - завантаження існуючої
						} else if strings.Contains(strings.ToLower(rowString), "load") && gameDataState != 2 {
							gameDataState = 2
							setGameDataState(screen, 2, gameModeState)
						}
					}
				  // Якщо рядок кількості кораблів
				} else if y == 15 {

					// Отримання символу комірки натиснутої позиції миші
					rune, _, _, _ := screen.GetContent(x, y)

					// Якщо число
					if rune >= '0' && rune <= '1' {
						setGameBoatsAmount(screen, 10, gameModeState, gameDataState)
					} else if rune >= '2' && rune <= '9' {
						setGameBoatsAmount(screen, int(rune-'0'), gameModeState, gameDataState)
					}
				}
			}

		}

		// Оновлення екрану
		screen.Show()
	}
}