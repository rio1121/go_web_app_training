package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Page webページのタイトルと本文
type Page struct {
	Title string
	Body  []byte
}

// StartWebApp Webサーバーの起動
func StartWebApp() {
	// /static/に対してハンドラーを登録
	// http.Dirの引数でディレクトリを指定.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static2"))))
	// nil = Default Handler
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func main() {
	// WebサーバーをGoroutineで起動
	go StartWebApp()

	fmt.Println("\n#####################################################################")
	fmt.Println("# Welcome to simple web application for golang by Ryunosuke Yamada. #")
	fmt.Println("# Please access to 'localhost:5555/static/'.                        #")
	fmt.Println("#                                                                   #")
	fmt.Println("# Press any key to exit.                                            #")
	fmt.Println("#####################################################################")

	scanner := bufio.NewScanner(os.Stdin)
	// 何かキーを押すとサーバーを終了.
	scanner.Scan()
}
