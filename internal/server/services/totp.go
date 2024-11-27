package services

import (
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type Totp struct {
	key *otp.Key
}

func New(accountName string) *Totp {
	return &Totp{}
}

func (t *Totp) Generate(accountName) error {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OyKeeper",
		AccountName: accountName,
	})
	if err != nil {
		panic(err)
	}

	t.key = key
	return nil
}

func (t *Totp) Load(url string) error {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return err
	}
	t.key = key
	return nil
}

func (t *Totp) GenerateImage() {
	var buf bytes.Buffer

	img, err := t.key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)
	os.WriteFile("qr-code.png", buf.Bytes(), 0644)
}

func (t *Totp) validate(code string) (bool, error) {

	valid, err := totp.Validate(code, t.key.Secret())

	if err != nil {
		fmt.Println(err.Error())
	}
	if valid {
		println("Valid passcode!")
		os.Exit(0)
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
}
