package iternal

import (
	"errors"
	"math/rand"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

var (
	symbols = []byte("1234567890qwertyuiopasdfghjklzxcvbnm!@#$^&*()_+QWERTYUIOPASDFGHJKLZXCVBNM")
)

type JwtWorker struct {
	secretKey   []byte
	serviceUuid string
}

func NewJwtWorker(serviceUuid string) *JwtWorker {

	w := &JwtWorker{
		secretKey:   []byte(generateRandomString(50)),
		serviceUuid: serviceUuid,
	}

	return w
}

func (w *JwtWorker) CreateJwt(playload map[string]any) (string, error) {

	p := jwtgo.MapClaims{
		"service_uuid": w.serviceUuid,
		"exp":          time.Now().Add(time.Minute * 30).Unix(),
	}

	for k, v := range playload {
		p[k] = v
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS512, p)

	result, err := token.SignedString(w.secretKey)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (w *JwtWorker) ParseJwt(token string) (jwtgo.MapClaims, error) {

	res := jwtgo.MapClaims{}

	_, err := jwtgo.ParseWithClaims(token, res, func(t *jwtgo.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return w.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func generateRandomString(length int) string {

	result := string(time.Now().Unix())

	for i := 0; i < length; i++ {
		result += string(symbols[rand.Intn(len(symbols))])
	}

	return result

}
