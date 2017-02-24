package main

import (
	"flag"
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
	var err = t.tmpl.Execute(w, r)
	if err != nil {
		log.Fatal("template compile error !")
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのポート")
	flag.Parse()

	room := newRoom()

	// ルーティング設定
	http.Handle("/", &templateHandler{filename: "chat.go.html"})
	http.Handle("/room", room)

	// 別スレッドでメッセージ受信送信の並列処理のリスナーを起動
	go room.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Server start error")
		log.Fatal(err.Error())
	}
}
