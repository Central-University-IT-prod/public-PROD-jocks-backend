package utils

import (
	"reflect"
	"testing"
)

type Payload struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestJWT(t *testing.T) {
	secret := "very secret"
	p := Payload{Field1: "field1", Field2: 22}

	token, err := GenerateJWT(p, secret)
	if err != nil {
		t.Error(err)
	}

	var res Payload
	err = ReadJWT(&res, token, secret)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(res, p) {
		t.Errorf("res: %v not equal to payload: %v", res, p)
	}
}
