package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// User ユーザー情報
type User struct {
	Name  string
	Intro string
}

// Chat チャット情報
type Chat struct {
	Name      string
	Message   string
	CreatedAt string
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

// insertChat チャットをDBへ登録
func insertChat(name, message string) {
	// データベースへのアクセス開始
	DBConnection, _ := sql.Open("sqlite3", "./chat.sql")

	// Openを呼び出す場合は必ず実行する
	defer DBConnection.Close()

	cmd := `INSERT INTO chat (name, message) VALUES (?, ?)`
	_, err := DBConnection.Exec(cmd, name, message)
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

// chatHandler - チャット処理ハンドラ
func chatHandler(writer http.ResponseWriter, req *http.Request) {
	name := req.FormValue("chat_name")
	message := req.FormValue("chat_message")
	insertChat(name, message)
	log.Println("Chat OK. 名前:", name, ", メッセージ:", message)
	http.Redirect(writer, req, "/chatroom/", http.StatusFound)
}

// renderIndexTemplate - User構造体のスライスのデータをテンプレートファイルを用いて表示
func renderIndexTemplate(writer http.ResponseWriter, u []User) {
	t, _ := template.ParseFiles("index/index.html")
	t.Execute(writer, u)
}

// renderChatTemplate - Chat構造体のスライスのデータをテンプレートファイルを用いて表示
func renderChatTemplate(writer http.ResponseWriter, c []Chat) {
	t, _ := template.ParseFiles("chatroom/chat.html")
	t.Execute(writer, c)
}

// indexHandler - 会員一覧表示ハンドラ
func indexHandler(writer http.ResponseWriter, req *http.Request) {
	// データベースへのアクセス開始
	DBConnection, _ := sql.Open("sqlite3", "./users.sql")

	// Openを呼び出す場合は必ず実行する
	defer DBConnection.Close()

	// DBのデータを読み出す
	cmd := "SELECT * from users"
	// Queryで得られる結果は必ず使用後にクローズすること
	rows, _ := DBConnection.Query(cmd)
	defer rows.Close()
	var users []User

	// データのスキャン
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Name, &u.Intro)
		if err != nil {
			log.Println(err)
		}
		users = append(users, u)
	}

	// render関数の呼び出し
	renderIndexTemplate(writer, users)
}

// chatroomHandker - チャットルームハンドラ
func chatroomHandler(writer http.ResponseWriter, req *http.Request) {
	// データベースへのアクセス開始
	DBConnection, _ := sql.Open("sqlite3", "./chat.sql")

	// Openを呼び出す場合は必ず実行する
	defer DBConnection.Close()

	// DBのデータを読み出す
	cmd := "SELECT * from chat"
	// Queryで得られる結果は必ず使用後にクローズすること
	rows, _ := DBConnection.Query(cmd)
	defer rows.Close()
	var chats []Chat

	// データのスキャン
	for rows.Next() {
		var c Chat
		err := rows.Scan(&c.Name, &c.Message, &c.CreatedAt)
		if err != nil {
			log.Println(err)
		}
		chats = append(chats, c)
	}

	// render関数の呼び出し
	renderChatTemplate(writer, chats)
}

// StartWebApp Webサーバーの起動
func StartWebApp() {
	// /static/に対してハンドラーを登録
	// http.Dirの引数でディレクトリを指定.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// 会員登録ハンドラ
	http.HandleFunc("/submit/", submitHandler)
	// 会員一覧ハンドラ
	http.HandleFunc("/index/", indexHandler)
	// チャットルームハンドラ
	http.HandleFunc("/chatroom/", chatroomHandler)
	// チャットハンドラ
	http.HandleFunc("/chat/", chatHandler)
	// リソースハンドラ - これがないとtemplate使用時にcssが適用されない.
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	// nil = Default Handler
	log.Fatal(http.ListenAndServe(":5555", nil))
}

// データベースがない場合は作成.
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

	// チャットデータ
	// データベースへのアクセス開始
	DBConnection, _ = sql.Open("sqlite3", "./chat.sql")

	// テーブル作成コマンド
	cmd = `CREATE TABLE IF NOT EXISTS chat(
				name STRING,
				message TEXT,
				created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')))`

	// コマンドを実行しつつ、エラーハンドリング
	_, err = DBConnection.Exec(cmd)
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
