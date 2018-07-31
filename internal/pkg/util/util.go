package util

import (
	"fmt"
	"strconv"
)

//ConverteInterfaceInt64 - converte tipo interface para tipo int64
func ConverteInterfaceInt64(valor interface{}) (int64, error) {
	valorConvertido, ok := valor.(int64)
	var err error
	if !ok {
		err = fmt.Errorf("Erro ao converter o valor %s para int64", valor)
	}
	return valorConvertido, err
}

// ConverteInterfaceFloat64 - converte tipo interface para tipo float64
func ConverteInterfaceFloat64(valor interface{}) (float64, error) {
	valorConvertido, ok := valor.(float64)
	var err error
	if !ok {
		err = fmt.Errorf("Erro ao converter o valor %s para float64", valor)
	}
	return valorConvertido, err
}

//ConvertStringInt64 - converte tipo string para tipo int64
func ConvertStringInt64(valor string) (int64, error) {
	valorConvertido, err := strconv.ParseInt(valor, 10, 64)
	return valorConvertido, err
}
