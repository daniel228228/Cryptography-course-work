package main

import (
	"fmt"

	"crypto-app/ui"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered: ", r)
		}
	}()

	var cmd string

	fmt.Println("Crypto App")
	Help()

	var f func() error

	for {
		f = nil
		fmt.Print("> ")
		fmt.Scan(&cmd)

		switch cmd {
		case "gen":
			f = ui.Gen
		case "encrypt":
			f = func() error { return ui.Encrypt(false, "", "", "", "", func(float64) {}) }
		case "decrypt":
			f = func() error { return ui.Decrypt(false, "", "", "", "", "", func(float64) {}) }
		case "help":
			f = Help
		case "exit":
			return
		default:
			Help()
			continue
		}

		if err := f(); err != nil {
			fmt.Println("Ошибка:", err)
		}
	}
}

func Help() error {
	fmt.Println("Помощь:\n\tgen - сгенерировать ключи\n\tencrypt - Зашифровать файл\n\tdecrypt - Расшифровать файл\n\thelp - вывод подсказки")

	return nil
}
