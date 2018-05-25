package util

import (
	"net"
	"net/http"
	"io/ioutil"
	"fmt"
	"compress/gzip"
	"errors"
)

const LocalHost = "127.0.0.1"

func LocalIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return LocalHost
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return LocalHost
}

func ReadResponse(resp *http.Response) (data []byte, err error) {
	defer resp.Body.Close()
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(resp.Body)
		data, err = ioutil.ReadAll(reader)
	default:
		data, err = ioutil.ReadAll(resp.Body)
	}
	if err != nil {
		uri := resp.Request.URL.String()
		err = errors.New(fmt.Sprintf("read %s body error for %s ", uri, err.Error()))
	}
	return
}
