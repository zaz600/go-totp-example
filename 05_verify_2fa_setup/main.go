package main

import (
	"bytes"
	"html/template"
	"image/png"
	"net/http"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

//тут будем хранить TOTP для одного пользователя
var key *otp.Key

func main() {
	//Настраиваем TOTP
	//для каждого пользователя TOTP ключ должен быть уникальным
	//В нашей программе ключ будет разный с каждым запуском (!)
	var err error
	key, err = totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "zaz600@example.com",
	})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", indexHandlerFunc)
	http.HandleFunc("/login/", loginHandlerFunc)
	http.HandleFunc("/2fa/", setup2FAHandlerFunc)
	http.HandleFunc("/qr.png", genQRCodeHandlerFunc)
	http.HandleFunc("/verify2fa/", verifi2faHandlerFunc)
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

//Отображает страницу с QR-кодом
func setup2FAHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//для простоты не обрабатываем ошибки
	t, _ := template.ParseFiles("templates/2fa.html")
	t.Execute(w, nil)
}

//Генерирует QR-код для добавления аккаунта в Яндекс.Ключ/Google.Authentificator
func genQRCodeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//для простоты не обрабатываем ошибки
	png.Encode(&buf, img)
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf.Bytes())

}

func verifi2faHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//Обрабатываем только POST-запрос
	if r.Method != "POST" {
		http.NotFound(w, r)
	}

	//для простоты не обрабатываем ошибки
	r.ParseForm()
	passcode := r.FormValue("passcode")

	valid := totp.Validate(passcode, key.Secret())
	if !valid {
		http.Error(w, "Неверный код", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`2ФА успешно настроена. <a href="/">На главную </a>`))
	//далее нам надо сохранить в базе key.Secret() пользователя
	//чтобы позднее верифицировать его одноразовые коды по этому секрету
}
