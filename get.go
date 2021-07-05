package main

import (
	"context"
	"errors"
	"io"
	"time"

	fetch "marwan.io/wasm-fetch"

	js "github.com/vugu/vugu/js"
)

func get(url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := fetch.Fetch(url, &fetch.Opts{
		Method: fetch.MethodGet,
		Signal: ctx,
	})
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	output := string(res.Body)
	if res.Status == 500 {
		js.Global().Call("alert", output)
		panic(errors.New(string(output)))
	}
	return string(output)
}

func getCode(url string) (string, int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := fetch.Fetch(url, &fetch.Opts{
		Method: fetch.MethodGet,
		Signal: ctx,
	})
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	return string(res.Body), res.Status
}

func post(url string, kind string, data io.Reader) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := fetch.Fetch(url, &fetch.Opts{
		Method: fetch.MethodPost,
		Signal: ctx,
		Body:   data,
		Headers: map[string]string{
			"Content-Type": kind,
		},
	})
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
	output := string(res.Body)
	if res.Status == 500 {
		js.Global().Call("alert", output)
		panic(errors.New(string(output)))
	}
	return string(output)
}
