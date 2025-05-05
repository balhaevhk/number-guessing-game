package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// начало игры
func Game() {
	attempts := 10
	fmt.Println("Вы начали игру \"Угадай число\"")
	fmt.Println("У вас всего", attempts, "попыток")
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(100) + 1
	reader := bufio.NewReader(os.Stdin)
	for  {
		fmt.Print("Введите ваше число: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		s := strings.ReplaceAll(input, " ", "")
		guess, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Вы ввели не число. Повторите")
			continue
		}
		if random < guess {
			fmt.Println("Секретное число меньше👇")
			Tips(random, guess)
		} else if random > guess{
			fmt.Println("Секретное число больше👆")
			Tips(random, guess)
		} else {
			fmt.Println("Вы угадали!🙌")
			return
		}
		attempts--
		if attempts > 0 {
			fmt.Println("У вас осталось", attempts, "попыток")
		}
		if attempts == 0 {
			fmt.Println("Попытки закончились")
			fmt.Println("Вы проиграли!😢")
			fmt.Println("А правильный ответ был", random)
			return
		}
	}
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



