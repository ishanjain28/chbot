package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ishanjain28/chbot/ch"
)

// PORT on which HTTP Server is started
var PORT = os.Getenv("PORT")

func main() {

	// Initalise database packge
	// db, err := db.Init()
	// if err != nil {
	// 	log.Fatalf("error in initalising datbase: %v", err)
	// }

	// defer db.Sess.Close()

	// fmt.Println("Connected to DB")

	// go bot.Start(db)

	sem := make(chan int, 10)
	result := make(chan string)
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 200; i++ {
		go ch.Scrap(i, sem, result, &wg)
	}

	for v := range result {
		fmt.Println(v)
	}

	wg.Wait()

}

// func handler(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != "POST" {
// 		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 		return
// 	}

// 	fmt.Println(r.Header)

// 	defer r.Body.Close()

// 	io.Copy(os.Stdout, r.Body)
// }
