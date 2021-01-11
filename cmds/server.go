package main

import "github.com/lpflpf/rpc"
import "strconv"
import "fmt"

func main() {
	serv := rpc.NewRpcServ("127.0.0.1:18080")
	serv.Impl("/conv/int2str", strconv.Itoa)
	serv.Impl("/conv/str2int", strconv.Atoi)
	serv.Impl("/math/add", func(a, b int) int { return a + b })
	serv.Impl("/fmt/sprintf", func(format string, data ...interface{}) string {
		fmt.Println(data)
		return fmt.Sprintf(format, data...)
	})
	serv.Start()
}
