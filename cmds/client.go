package main

import "fmt"
import "github.com/lpflpf/rpc"

type Conv struct {
	Int2Str func(int) string                `rpc:"conv/int2str"`
	Str2Int func(input string) (int, error) `rpc:"conv/str2int"`
}

type Math struct {
	Add func(int, int) int `rpc:"math/add"`
}

type Fmt struct {
	Sprintf func(string, ...interface{}) string `rpc:"fmt/sprintf"`
}

func main() {
	conv := &Conv{}
	rpc.Connect("http://127.0.0.1:18080", conv)
	fmt.Println(conv.Int2Str(123), conv.Int2Str(456))
	fmt.Println(conv.Str2Int("1234"))

	math := &Math{}
	rpc.Connect("http://127.0.0.1:18080", math)
	fmt.Println(math.Add(1, 2))

	format := &Fmt{}
	rpc.Connect("http://127.0.0.1:18080", format)
	fmt.Println(format.Sprintf("%s %s", "hello", "world"))
}
