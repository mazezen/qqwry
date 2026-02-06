package qqwry

import (
	"io"
	"net/http"
	"os"
)

func Download(dst string) error {
	url := "https://github.com/metowolf/qqwry.dat/releases/latest/download/qqwry.dat"

	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
