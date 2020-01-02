package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Page webページのタイトルと本文
type Page struct {
	Title string
	Body  []byte
}

// DBConnection ..
var DBConnection *sql.DB

// insertUser ユーザーをDBへ登録
func insertUser(name, intro string) {
	// データベースへのアクセス開始
	DBConnection, _ := sql.Open("sqlite3", "./users.sql")

	// Openを呼び出す場合は必ず実行する
	defer DBConnection.Close()

	cmd := `INSERT INTO users (name, intro) VALUES (?, ?)`
	_, err := DBConnection.Exec(cmd, name, intro)
	if err != nil {
		log.Println(err)
	}
}

// submitHandler - 会員登録処理ハンドラ
func submitHandler(writer http.ResponseWriter, req *http.Request) {
	userName := req.FormValue("user_name")
	userIntroduction := req.FormValue("user_introduction")
	insertUser(userName, userIntroduction)
	log.Println("Submit OK. ユーザー名:", userName, ", 自己紹介文:", userIntroduction)
	http.Redirect(writer, req, "/static/", http.StatusFound)
}

// StartWebApp Webサーバーの起動
func StartWebApp() {
	// /static/に対してハンドラーを登録
	// http.Dirの引数でディレクトリを指定.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// /submit/に対してハンドラーを登録
	http.HandleFunc("/submit/", submitHandler)
	// nil = Default Handler
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func init() {
	// データベースへのアクセス開始
	DBConnection, _ := sql.Open("sqlite3", "./users.sql")

	// Openを呼び出す場合は必ず実行する
	defer DBConnection.Close()

	// テーブル作成コマンド
	cmd := `CREATE TABLE IF NOT EXISTS users(
				name STRING,
				intro STRING)`

	// コマンドを実行しつつ、エラーハンドリング
	_, err := DBConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}
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
