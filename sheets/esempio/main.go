package main

import (
	"net/http"
	"sheets"
)


func main() {
	http.HandleFunc("/",sheets.SheetAppend)
	http.ListenAndServe(":8000",nil)
}