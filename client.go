package rpc

import "io/ioutil"

import "bytes"
import "strings"
import "reflect"
import "errors"
import "encoding/json"
import "net/http"

type RpcClient struct {
	serv *http.Server
	mux  *http.ServeMux
}

// struct BIND RPC
func Connect(addr string, iface interface{}) error {
	rv := reflect.ValueOf(iface).Elem()
	rt := reflect.TypeOf(iface).Elem()

	if rt.Kind() != reflect.Struct {
		return errors.New("")
	}

	for i := 0; i < rt.NumField(); i++ {
		if requestPath := rt.Field(i).Tag.Get("rpc"); requestPath == "" {
			continue
		} else {
			fieldType := rt.Field(i).Type
			rv.Field(i).Set(reflect.MakeFunc(fieldType, func(params []reflect.Value) []reflect.Value {
				requestBody := []string{}
				for _, param := range params {
					raw, _ := json.Marshal(param.Interface())
					requestBody = append(requestBody, string(raw))
				}
				body, _ := json.Marshal(requestBody)

				requestUri := strings.Trim(addr, "/") + "/" + strings.Trim(requestPath, "/")
				resp, _ := http.Post(requestUri, "application/json", bytes.NewReader(body))
				defer resp.Body.Close()
				data, _ := ioutil.ReadAll(resp.Body)

				ret := []reflect.Value{}
				responseStr := []string{}
				_ = json.Unmarshal(data, &responseStr)
				for i := 0; i < fieldType.NumOut(); i++ {
					val := reflect.New(fieldType.Out(i))
					_ = json.Unmarshal([]byte(responseStr[i]), val.Interface())
					ret = append(ret, val.Elem())
				}

				return ret
			}))
		}
	}

	return nil
}
