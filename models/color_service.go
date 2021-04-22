package models

import (
	"encoding/hex"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/data"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/go-trace"
	"math/rand"
	"strconv"
	"strings"
)

type ColorServiceInterface interface {
	GetByName(name string) (res Color, err error)
	GetRandom() (res Color, err error)
}

func NewColorService() (svc *colorService) {
	var cl []Color
	cm := map[string]Color{}

	if err := json.Unmarshal([]byte(data.ColorsDataText), &cl); err != nil {
		_ = trace.TraceError(err)
	}

	for _, c := range cl {
		cm[c.Name] = c
	}

	return &colorService{
		cl: cl,
		cm: cm,
	}
}

type colorService struct {
	cl []Color
	cm map[string]Color
}

func (svc *colorService) GetByName(name string) (res Color, err error) {
	res, ok := svc.cm[name]
	if !ok {
		return res, errors.ErrorModelNotFound
	}
	return res, err
}

func (svc *colorService) GetRandom() (res Color, err error) {
	if len(svc.cl) == 0 {
		hexStr, err := svc.getRandomColorHex()
		if err != nil {
			return res, err
		}
		return Color{Hex: hexStr}, nil
	}

	idx := rand.Intn(len(svc.cl))
	return svc.cl[idx], nil
}

func (svc *colorService) getRandomColorHex() (res string, err error) {
	n := 6
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i], err = svc.getRandomHexChar()
		if err != nil {
			return res, err
		}
	}
	return strings.Join(arr, ""), nil
}

func (svc *colorService) getRandomHexChar() (res string, err error) {
	n := rand.Intn(16)
	b := []byte(strconv.Itoa(n))
	h := make([]byte, 1)
	hex.Encode(h, b)
	return string(h), nil
}

var ColorService *colorService
