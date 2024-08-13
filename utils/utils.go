package utils

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

const DEFAULT_CONTAINER_PORT = "8888"

func GenerateRandomPassword(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)
    password := make([]byte, length)
    for i := range password {
        password[i] = charset[random.Intn(len(charset))]
    }
    return string(password)
}

func GetRandomFreePort() (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), nil
}