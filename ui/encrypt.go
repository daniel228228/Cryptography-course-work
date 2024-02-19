package ui

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"crypto-app/internal/benaloh"
	"crypto-app/internal/twofish"
)

func Encrypt(gui bool, inputFilename, publicKeyFilename, outputDirname, encryptModeName string, progress func(float64)) error {
	for {
		if !gui {
			fmt.Print("Введите путь к шифруемому файлу: ")
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
			fmt.Print("В какую папку сохранить файл и ключ: ")
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

	var encryptMode twofish.Mode

	for {
		if !gui {
			fmt.Print("Выберите режим шифрования (ECB, CBC, CFB, OFB): ")
			fmt.Scan(&encryptModeName)
		}

		if encryptModeName == "exit" {
			return nil
		}

		var err error

		if encryptMode, err = twofish.ParseMode(encryptModeName); err != nil {
			if !gui {
				fmt.Println("Такой режим отсутствует в списке. Повторите попытку или введите exit для выхода")
			} else {
				return err
			}
		} else {
			break
		}
	}

	if !gui {
		var cont string
		fmt.Print("Запустить шифрование выбранного файла? (Y/n): ")
		fmt.Scan(&cont)

		if cont == "n" {
			fmt.Println("Действие отменено")
			return nil
		}
	}

	twofishKey := make([]byte, twofish.BlockSize)
	if _, err := rand.Read(twofishKey); err != nil {
		return err
	}

	iv := make([]byte, len(twofishKey))
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	t, err := twofish.NewContext(twofishKey, iv, encryptMode)
	if err != nil {
		return err
	}

	b := benaloh.NewBenaloh()

	if err := b.ReadPublicKey(publicKeyFilename); err != nil {
		return err
	}

	k := t.GetKey(filepath.Base(inputFilename))

	encryptedKey := []string{}

	for i := 0; i < len(k); i++ {
		enc, err := b.Encrypt(new(big.Int).SetBytes([]byte{k[i]}))
		if err != nil {
			return err
		}

		encryptedKey = append(encryptedKey, enc.String())
	}

	if err := t.WriteEnc([]byte(strings.Join(encryptedKey, " ")), filepath.Dir(outputDirname)+"/"+filepath.Base(inputFilename)+".key", 0644); err != nil {
		return err
	}

	return t.EncryptFile(inputFilename, filepath.Dir(outputDirname)+"/"+filepath.Base(inputFilename)+".enc", progress)
}
