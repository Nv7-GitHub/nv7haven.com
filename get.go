package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	js "github.com/vugu/vugu/js"
)

func get(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	if res.StatusCode == 500 {
		js.Global().Call("alert", output)
		panic(errors.New(string(output)))
	}
	return string(output)
}

func getCode(url string) (string, int) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	return string(output), res.StatusCode
}

func post(url string, kind string, data io.Reader) string {
	res, err := http.Post(url, kind, data)
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	defer res.Body.Close()
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	if res.StatusCode == 500 {
		js.Global().Call("alert", output)
		panic(errors.New(string(output)))
	}
	return string(output)
}
