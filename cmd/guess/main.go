package main

import (
	"fmt"
	"guess/internal/game"
	"time"

)

func main() {
	start := true
	for start {
		random, attempts, levelGame := game.StartGame()
		res, attemptsUsed := game.Game(random, attempts)

		
		result := game.GameResult{
			Level:        levelGame,
			MaxNumber:    random,
			AttemptsUsed: attemptsUsed,
			Win:          res,
			Timestamp:    time.Now().Format("02:01:2006 15:04:05"),
		}

		if err := game.SaveResult(result); err != nil {
			fmt.Println("Ошибка при сохранении результата:", err)
		} else {
			fmt.Println("Результат игры успешно сохранён!")
		}

		start = game.GameAgain()
	}
}



