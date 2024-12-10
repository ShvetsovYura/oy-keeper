package service

import (
	"bytes"
	"image/png"
	"os"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type Totp struct {
	key *otp.Key
}

func NewTotp() *Totp {
	return &Totp{}
}

func (t *Totp) Generate(login string) error {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OyKeeper",
		AccountName: login,
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

func (t *Totp) GetKey() *otp.Key {
	return t.key
}

func (t *Totp) GenerateImage() {
	var buf bytes.Buffer

	img, err := t.key.Image(100, 100)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)
	os.WriteFile("qr-code.png", buf.Bytes(), 0644)
}

func (t *Totp) Validate(code string) bool {
	return totp.Validate(code, t.key.Secret())
}
