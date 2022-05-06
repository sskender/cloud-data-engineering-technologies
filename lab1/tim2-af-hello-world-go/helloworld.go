package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "tim2 says hello world from go")
}

func main() {
    listenAddr := ":8080"
    if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
        listenAddr = ":" + val
    }
    http.HandleFunc("/api/tim2-af-hello-world-go", helloHandler)
    log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
    log.Fatal(http.ListenAndServe(listenAddr, nil))
}
