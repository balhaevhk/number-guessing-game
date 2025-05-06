package game

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type GameResult struct {
	Level        string    `json:"level"`
	MaxNumber    int       `json:"maxNumber"`
	AttemptsUsed int       `json:"attemptsUsed"`
	Win          bool      `json:"win"`
	Timestamp    string `json:"timestamp"`
}

// начало игры
func StartGame() (random int, attempts int, levelGame string) {
	fmt.Println("Вы начали игру \"Угадай число\"")
	fmt.Println("Для начала выберите номер уровеня:")
	fmt.Println("1 - easy; 2 - medium; 3 - hard")

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
			fmt.Println("1 - easy; 2 - medium; 3 - hard")
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
		fmt.Print("Введите ваше число: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		s := strings.ReplaceAll(input, " ", "")
		guess, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Вы ввели не число. Повторите")
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
			fmt.Println("Вы угадали!🙌")
			return true, attUsed
		}
		attempts--
		if attempts > 0 {
			fmt.Println("У вас осталось", attempts, "попыток")
		}
		if attempts == 0 {
			fmt.Println("Попытки закончились")
			fmt.Println("Вы проиграли!😢")
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
	const filePath = "../../storage/results.json"

	var results []GameResult

	// Проверка, существует ли файл и чтение старых данных
	data, err := os.ReadFile(filePath)
	if err != nil {
		// Если файл не существует, не выводим ошибку, просто создаем новый с результатами
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Файл не найден, будет создан новый.")
		} else {
			return fmt.Errorf("не удалось прочитать файл: %w", err)
		}
	}

	// Если файл существует, пытаемся распарсить его
	if len(data) > 0 {
		if err := json.Unmarshal(data, &results); err != nil {
			return fmt.Errorf("ошибка при чтении JSON: %w", err)
		}
	}

	// Добавляем новый результат
	results = append(results, result)

	// Сохраняем обратно
	newData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге: %w", err)
	}

	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	return nil
}

