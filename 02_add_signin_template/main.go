package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandlerFunc)
	http.HandleFunc("/login/", loginHandlerFunc)
	http.ListenAndServe(":3000", nil)
}

func indexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//для простоты не обрабатываем ошибки
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func loginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//Обрабатываем только POST-запрос
	if r.Method != "POST" {
		http.NotFound(w, r)
	}

	//для простоты не обрабатываем ошибки
	r.ParseForm()
	user := r.FormValue("user")
	password := r.FormValue("password")

	//Проверяем логин и пароль
	if !(user == "zaz600" && password == "123") {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	w.Write([]byte("hello " + user))
}
