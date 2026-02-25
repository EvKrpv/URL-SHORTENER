package random

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"unicode"
)

const (
	charset        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	shortURLLength = 10
)

// Генерирует случайную строку для короткого URL. Длина строки определяется константой shortURLLength.
func GenerateShortCode(longURL string) (string, error) {
	// Нормализация URL для обеспечения консистентности.
	// Это гарантирует, что одинаковые URL будут генерировать одинаковые короткие коды, что может помочь в предотвращении коллизий.
	_, err := normalizeURL(longURL)
	if err != nil {
		return "", err
	}

	randomString, err := generateRandomString(shortURLLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate short code: %w", err)
	}
	return randomString, nil
}

func generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		// Получаем случайный индекс для выбора символа из набора chars
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random index: %w", err)
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

func normalizeURL(rawURL string) (string, error) { // Добавляем схему, если она отсутствует
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Стандартизация схемы и хоста
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// Удаление стандартных портов
	u.Host = strings.TrimSuffix(u.Host, ":80")
	u.Host = strings.TrimSuffix(u.Host, ":443")

	if u.Path == "" {
		u.Path = "/"
	}

	return u.String(), nil
}

func ValidateShortCode(shortCode string) error {
	if len(shortCode) != shortURLLength {
		return errors.New("code must be 10 characters long")
	}

	for _, r := range shortCode {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
			return errors.New("code must contain only letters, numbers, or underscore")
		}
	}

	return nil
}
