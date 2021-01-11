package rpc

import "net/http"
import "reflect"
import "encoding/json"
import "io/ioutil"

type RpcServ struct {
	serv *http.Server
	mux  *http.ServeMux
}

func (rs *RpcServ) Impl(router string, f interface{}) {
	rs.mux.HandleFunc(router, func(rw http.ResponseWriter, request *http.Request) {
		rt := reflect.TypeOf(f)
		requestBody, _ := ioutil.ReadAll(request.Body)

		requestData := []string{}
		_ = json.Unmarshal(requestBody, &requestData)

		params := []reflect.Value{}

		num := rt.NumIn()
		if rt.IsVariadic() {
			num = num - 1
		}
		for i := 0; i < num; i++ {
			val := reflect.New(rt.In(i))
			json.Unmarshal([]byte(requestData[i]), val.Interface())
			params = append(params, val.Elem())
		}

		call := reflect.ValueOf(f)
		result := []reflect.Value{}

		if rt.IsVariadic() {
			val := reflect.MakeSlice(rt.In(num), 0, 0).Interface()
			json.Unmarshal([]byte(requestData[num]), &val)
			params = append(params, reflect.ValueOf(val))
			result = call.CallSlice(params)
		} else {
			result = call.Call(params)
		}

		response := []string{}
		for _, res := range result {
			val, _ := json.Marshal(res.Interface())
			response = append(response, string(val))
		}

		data, _ := json.Marshal(response)
		rw.Write(data)
	})
}

func (rs *RpcServ) Start() {
	rs.serv.Handler = rs.mux
	rs.serv.ListenAndServe()
}

func NewRpcServ(addr string) *RpcServ {
	return &RpcServ{
		serv: &http.Server{Addr: addr},
		mux:  http.NewServeMux(),
	}
}
