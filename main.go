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
	// nil = Default Handler
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func main() {
	go StartWebApp()

	fmt.Println("\n#####################################################################")
	fmt.Println("# Welcome to simple web application for golang by Ryunosuke Yamada. #")
	fmt.Println("#####################################################################")
	fmt.Print("\nPress any key to exit.")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return
}
