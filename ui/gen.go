package ui

import (
	"log"
	"time"

	"crypto-app/internal/benaloh"
)

func Gen() error {
	test := benaloh.TestSolovayStrassen

	log.Print("Генерация ключей...")
	b := benaloh.NewBenaloh()
	if err := b.GenKey(test); err != nil {
		return err
	}

	log.Println("ok")

	curTime := time.Now().Format("20060102T150405")

	if err := b.WritePublicKey("./keys", "public_"+curTime, 0644); err != nil {
		return err
	}

	if err := b.WritePrivateKey("./keys", "private_"+curTime, 0644); err != nil {
		return err
	}

	log.Println("Ключи сохранены в папке keys")

	return nil
}
