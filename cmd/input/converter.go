package input

import (
	"strconv"
)

var converters = make(map[string]Converter)

type Converter interface {
	Name() string
	Convert(input string) (data []byte, err error)
}


func GetInstance(name string) Converter {
	return converters[name]
}

func RegisterInstance(converter Converter) {
	converters[converter.Name()] = converter
}

type StringConverter struct {

}

func (c *StringConverter) Name() string {
	return "string"
}

func (c *StringConverter) Convert(input string) (data []byte, err error) {
	return []byte(input), nil
}


type HexadecimalConverter struct {

}


func (c *HexadecimalConverter) Name() string {
	return "hex"
}

func (c *HexadecimalConverter) Convert(input string) (data []byte, err error) {

	//fmt.Printf("input: %s\n", input)

	data = make([]byte, len(input) / 2)
	l := len(data)
	var n int64
	for i := 0; i < l; i++ {
			n, err = strconv.ParseInt(input[i*2:i*2+2], 16, 0)
		if err != nil {
			data = nil
			return
		}
		data[i] = byte(n)
	}
	return
}

func init()  {
	RegisterInstance(new(StringConverter))
	RegisterInstance(new(HexadecimalConverter))
}