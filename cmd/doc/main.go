package main

import (
	"flag"
	"log"
	"net/http"
)

func serveDocs() {
	portDoc := flag.String("port-doc", "8080", "Port to run the documentation server on")
	flag.Parse()

	fs := http.FileServer(http.Dir("./doc"))
	http.Handle("/doc/", http.StripPrefix("/doc/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	log.Println("Serving documentation at http://localhost:" + *portDoc + "/doc")
	log.Fatal(http.ListenAndServe(":"+*portDoc, nil))
}

func main() {
	serveDocs()
}
