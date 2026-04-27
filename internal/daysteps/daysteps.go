package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parseData := make([]string, 2)

	split := strings.Split(data, ",")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("Неверный формат")
	}
	parseData[0] = split[0]
	parseData[1] = split[1]

	if parseData[0] == "" {
		return 0, 0, fmt.Errorf("Не указано количество шагов")
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка парсинга шагов: %w", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("Количество шагов должно быть больше 0")
	}

	if parseData[1] == "" {
		return 0, 0, fmt.Errorf("Время не указано")
	}

	duration, err := time.ParseDuration(parseData[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка парсинга времени: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("Длительность должна быть больше 0")
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distance := stepLength * float64(steps)
	distanceKm := distance / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %0.2f км.\nВы сожгли %0.2f ккал.\n", steps, distanceKm, calories)
}
