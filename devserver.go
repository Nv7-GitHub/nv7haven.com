// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/vugu/vugu/devutil"
)

func main() {
	l := "127.0.0.1:8844"
	log.Printf("Starting HTTP Server at %q", l)

	wc := devutil.NewWasmCompiler().SetDir(".")
	mux := devutil.NewMux()
	mux.Match(devutil.NoFileExt, devutil.DefaultAutoReloadIndex.Replace(
		`<!-- styles -->`,
		`<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css">
		<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	  <script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js" integrity="sha384-9/reFTGAW83EW2RDu2S0VKaIzap3H66lZH81PoYlFhbGU+6BZp6G7niu735Sk7lN" crossorigin="anonymous"></script>
	  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js" integrity="sha384-B4gt1jrGC7Jh4AgTPSdUtOBvfO8shuf57BaghqFfPlYxofvL8/KUEfYiJOMMV+rV" crossorigin="anonymous"></script>
	  <script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script>
		<meta http-equiv='cache-control' content='no-cache'>
	<meta http-equiv='expires' content='0'>
	<meta http-equiv='pragma' content='no-cache'> `))
	mux.Exact("/main.wasm", devutil.NewMainWasmHandler(wc))
	mux.Exact("/wasm_exec.js", devutil.NewWasmExecJSHandler(wc))
	mux.Default(devutil.NewFileServer().SetDir("."))

	log.Fatal(http.ListenAndServe(l, mux))
}
