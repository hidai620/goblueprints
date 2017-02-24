package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	tmpl     *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.tmpl = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
		fmt.Println("template was compiled")
	})
	var err = t.tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal("template compile error !")
	}
}

func main() {
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte(`
	//	<h1>hello</h1>
	//	`))
	//})
	room := newRoom()

	// ルーティング設定
	http.Handle("/", &templateHandler{filename: "chat.go.html"})
	http.Handle("/room", room)

	// メッセージ受信送信の並列処理のリスナーを起動
	go room.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server start error")
		log.Fatal(err.Error())
	}
}
