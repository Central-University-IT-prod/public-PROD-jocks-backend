package utils

import (
	"encoding/json"
	"math/rand"
	"slices"

	"github.com/google/uuid"
)

func ToMap(item any) (map[string]any, error) {
	res := make(map[string]any)
	data, err := json.Marshal(item)

	if err != nil {
		return res, err
	}

	err = json.Unmarshal(data, &res)
	return res, err
}

func Map[T, R any](s []T, mod func(T) R) []R {
	res := make([]R, len(s))
	for i, item := range s {
		res[i] = mod(item)
	}

	return res
}

func RandomChoises[T any](s []T, n int) []T {
	sCopy := make([]T, len(s))
	copy(sCopy, s)
	res := []T{}

	for len(res) < n && len(sCopy) > 0 {
		i := rand.Intn(len(sCopy))
		res = append(res, sCopy[i])
		sCopy = slices.Delete(sCopy, i, i+1)
	}

	return res
}

func GenerateToken() string {
	return uuid.NewString()
}
