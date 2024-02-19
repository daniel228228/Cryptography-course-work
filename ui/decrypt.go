package ui

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"crypto-app/internal/benaloh"
	"crypto-app/internal/twofish"
)

func Decrypt(gui bool, inputFilename, keyFilename, publicKeyFilename, privateKeyFilename, outputDirname string, progress func(float64)) error {
	for {
		if !gui {
			fmt.Print("Введите путь к дешифруемому файлу: ")
			fmt.Scan(&inputFilename)
		}

		if inputFilename == "exit" {
			return nil
		}

		if stat, err := os.Stat(inputFilename); err != nil || stat.IsDir() {
			if !gui {
				fmt.Println("Файл не существует. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else {
			break
		}
	}

	for {
		if !gui {
			fmt.Print("Введите путь к файлу с ключом (.key): ")
			fmt.Scan(&keyFilename)
		}

		if keyFilename == "exit" {
			return nil
		}

		if stat, err := os.Stat(keyFilename); err != nil || stat.IsDir() {
			if !gui {
				fmt.Println("Файл не существует. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else {
			break
		}
	}

	for {
		if !gui {
			fmt.Print("Введите путь к файлу с публичным ключом: ")
			fmt.Scan(&publicKeyFilename)
		}

		if publicKeyFilename == "exit" {
			return nil
		}

		if stat, err := os.Stat(publicKeyFilename); err != nil || stat.IsDir() {
			if !gui {
				fmt.Println("Файл не существует. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else {
			break
		}
	}

	for {
		if !gui {
			fmt.Print("Введите путь к файлу с приватным ключом: ")
			fmt.Scan(&privateKeyFilename)
		}

		if privateKeyFilename == "exit" {
			return nil
		}

		if stat, err := os.Stat(privateKeyFilename); err != nil || stat.IsDir() {
			if !gui {
				fmt.Println("Файл не существует. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else {
			break
		}
	}

	for {
		if !gui {
			fmt.Print("В какую папку сохранить файл: ")
			fmt.Scan(&outputDirname)
		}

		if outputDirname == "exit" {
			return nil
		}

		if stat, err := os.Stat(outputDirname); err != nil {
			if !gui {
				fmt.Println("Папка не существует. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else if !stat.IsDir() {
			if !gui {
				fmt.Println("Это не папка. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else {
			if outputDirname[len(outputDirname)-1] != '/' {
				outputDirname += "/"
			}

			break
		}
	}

	if !gui {
		var cont string
		fmt.Print("Запустить расшифровку выбранного файла? (Y/n): ")
		fmt.Scan(&cont)

		if cont == "n" {
			fmt.Println("Действие отменено")
			return nil
		}
	}

	encryptedKey, err := twofish.ReadEnc(keyFilename)
	if err != nil {
		return err
	}

	b := benaloh.NewBenaloh()

	if err := b.ReadPublicKey(publicKeyFilename); err != nil {
		return err
	}

	if err := b.ReadPrivateKey(privateKeyFilename); err != nil {
		return err
	}

	encKey := strings.Split(string(encryptedKey), " ")
	if len(encKey) == 0 {
		return twofish.ErrIncorrectKey
	}

	decryptedKey := []byte{}

	for _, block := range encKey {
		val, ok := new(big.Int).SetString(block, 10)
		if !ok {
			return twofish.ErrIncorrectKey
		}

		dec, err := b.Decrypt(val)
		if err != nil {
			return err
		}

		bytes := dec.Bytes()
		if len(bytes) == 0 {
			bytes = []byte{0}
		}

		decryptedKey = append(decryptedKey, bytes...)
	}

	t, outputFilename, err := twofish.NewContextFromData(decryptedKey)
	if err != nil {
		return err
	}

	return t.DecryptFile(inputFilename, filepath.Dir(outputDirname)+"/"+filepath.Base(outputFilename), progress)
}
