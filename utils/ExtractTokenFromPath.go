package utils

import (
	"fmt"
	"strings"
)

func ExtractTokenFromPath(path string) (string, error) {
	prefix := "/api/account/verify/"

	// Проверяем, что путь начинается с нужного префикса
	if !strings.HasPrefix(path, prefix) {
		return "", fmt.Errorf("invalid path format")
	}

	// Извлекаем часть после префикса
	token := strings.TrimPrefix(path, prefix)

	// Удаляем возможные слеши в начале/конце
	token = strings.Trim(token, "/")

	// Проверяем, что токен не пустой
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}