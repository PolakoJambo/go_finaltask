package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	//Не уверен что эта костанта нужна здесь => lenStep = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parseData := make([]string, 3)
	split := strings.Split(data, ",")

	if len(split) != 3 {
		return 0, "", 0, fmt.Errorf("Неверный формат")
	}

	parseData[0] = split[0]
	parseData[1] = split[1]
	parseData[2] = split[2]

	if parseData[0] == "" {
		return 0, "", 0, fmt.Errorf("Не указано количество шагов")
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка парсинга шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Количество шагов должно быть больше 0")
	}

	if parseData[1] == "" {
		return 0, "", 0, fmt.Errorf("Активность не указана")
	}

	action := parseData[1]

	if parseData[2] == "" {
		return 0, "", 0, fmt.Errorf("Время не указано")
	}

	duration, err := time.ParseDuration(parseData[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка парсинга времени: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Время должно быть больше 0")
	}
	return steps, action, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * float64(height)
	return (stepLength * float64(steps)) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	durationInHours := duration.Hours()
	return distance(steps, height) / durationInHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, action, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Некорректно получена информация о тренировке: %w", err)
	}

	speed := meanSpeed(steps, height, duration)
	dist := distance(steps, height)
	var calories float64

	switch action {

	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("Ошибка: %w", err)
		}

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("Ошибка: %w", err)
		}

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", action)
	}
	hours := duration.Hours()

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", action, hours, dist, speed, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Длительность бега должна быть больше 0")
	}
	speed := meanSpeed(steps, height, duration)

	durationInMinute := duration.Minutes()

	calories := (weight * speed * durationInMinute) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Длительность ходьбы должна быть больше 0")
	}
	speed := meanSpeed(steps, height, duration)

	durationInMinute := duration.Minutes()

	calories := (weight * speed * durationInMinute) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
