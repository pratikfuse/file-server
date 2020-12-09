package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func main() {
	errChan := make(chan error)
	path := GetDirectory()
	portFlag := flag.String("port", "8000", "port for the file server")
	flag.Parse()
	port := fmt.Sprintf(":%s", *portFlag)

	fileServer := http.FileServer(http.Dir(path))

	router := http.NewServeMux()

	router.Handle("/",
		http.StripPrefix("/", fileServer))
	server := http.Server{
		Addr: port,
		Handler: router,
	}
	go func(errChan chan<- error) {
		err := server.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}(errChan)

	fmt.Printf("Server running at 0.0.0.0%s at %s", port, path)

	if err := <-errChan; err != nil {
		fmt.Println(err)
	}

}
