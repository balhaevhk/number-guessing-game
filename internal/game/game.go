package game

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

type GameResult struct {
	Level        string    `json:"level"`
	MaxNumber    int       `json:"maxNumber"`
	AttemptsUsed int       `json:"attemptsUsed"`
	Win          bool      `json:"win"`
	Timestamp    string 	 `json:"timestamp"`
}

// начало игры
func StartGame() (random int, attempts int, levelGame string) {
	fmt.Println("Вы начали игру \"Угадай число\"")
	fmt.Println("Для начала выберите номер уровеня:")
	fmt.Println(green("1 - easy; ", yellow("2 - medium; "), red("3 - hard")))

	reader := bufio.NewReader(os.Stdin)
	var number int
	var err error

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		number, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Введите корректное значение")
			continue	
		}

		random, attempts, levelGame, err = Complexity(number)
		if err != nil {
			fmt.Println("Ошибка:", err)
			fmt.Println("Введите один из предложенных уровней")
			fmt.Println(green("1 - easy; ", yellow("2 - medium; "), red("3 - hard")))
			continue
		}
		break
	}
	return random, attempts, levelGame

}

// начало игры
func Game(randomNumber int, attempts int) (bool, int) {
	listAttempts := make([]int, 0)
	fmt.Println("У вас всего", attempts, "попыток")
	fmt.Println("Нужно угадать число от 1 до", randomNumber)
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(randomNumber) + 1
	reader := bufio.NewReader(os.Stdin)
	var attUsed int
	for attUsed = 1; attempts > 0; attUsed++ {
		fmt.Print(yellow("Введите ваше число: "))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		s := strings.ReplaceAll(input, " ", "")
		guess, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Вы ввели не число. Повторите")
			continue
		} else if (guess > randomNumber) {
			fmt.Println("Вы ввели число больше максимального")
			continue
		}
		listAttempts = append(listAttempts, guess)
		if random < guess {
			fmt.Println("Секретное число меньше👇")
			fmt.Println(LastAppempts(listAttempts))
			Tips(random, guess)
		} else if random > guess{
			fmt.Println("Секретное число больше👆")
			fmt.Println(LastAppempts(listAttempts))
			Tips(random, guess)
		} else {
			fmt.Println(green("Вы угадали!🙌"))
			return true, attUsed
		}
		attempts--
		if attempts > 0 {
			fmt.Println("У вас осталось", attempts, "попыток")
		}
		if attempts == 0 {
			fmt.Println("Попытки закончились")
			fmt.Println(red("Вы проиграли!😢"))
			fmt.Println("А правильный ответ был", random)
			return false, attUsed
		}
	}
	return false, attUsed
}

// подсказки
func Tips(random int, guess int) {
	diff := random - guess
	if diff < 0 {
		diff = -diff
	}
	if diff <= 5 {
		fmt.Println("Горячо 🔥")
	} else if diff <= 15 {
		fmt.Println("Тепло  🙂")
	} else {
		fmt.Println("Холодно ❄️")
	}
}

// прошлые попытки
func LastAppempts(attempts []int) string {
	var s []string

	for _, v := range attempts {
		s = append(s, strconv.Itoa(v))
	}
	result := strings.Join(s, ", ")
	return fmt.Sprintf("Ранее вы ввели %s", result)
}

// сложность игры
func Complexity(level int) (random int, attempts int, levelGame string,  err error) {
	switch level {
	case 1:
		random = 50
		attempts = 15
		levelGame = "easy"
	case 2:
		random = 100
		attempts = 10
		levelGame = "medium"
	case 3:
		random = 200
		attempts = 5
		levelGame = "hurd"
	default:
		err = fmt.Errorf("такого уровня нет")
	}
	return random, attempts, levelGame, err
}

// повтор игры
func GameAgain() (start bool) {
	fmt.Println("Хотите сыграть? 1 - да; 0 - нет")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		s := strings.ReplaceAll(input, " ", "")
		num, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Введите 1 или 0")
			continue
		}
		if num == 1 {
			return  true
		} else if num == 0 {
			return false
		}
	}
}

// сохранение результата
func SaveResult(result GameResult) error {
	// Получаем текущую рабочую директорию (ту, откуда запущена программа)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить текущую директорию: %w", err)
	}

	// Собираем путь к results.json
	filePath := filepath.Join(wd, "../../storage", "results.json")

	var results []GameResult

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Файл не найден, будет создан новый.")
		} else {
			return fmt.Errorf("не удалось прочитать файл: %w", err)
		}
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &results); err != nil {
			return fmt.Errorf("ошибка при чтении JSON: %w", err)
		}
	}

	results = append(results, result)

	newData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге: %w", err)
	}

	// Создаем директорию storage, если её нет (на всякий случай)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию storage: %w", err)
	}

	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	return nil
}

