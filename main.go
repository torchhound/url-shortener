package main

import (
  "fmt"
  "net/http"
  "html"
  "log"
  "math/rand"
  "time"
)

func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

func randomString(len int) string {
    bytes := make([]byte, len)
    for i := 0; i < len; i++ {
        if(randomInt(0, 2) == 1) {
          bytes[i] = byte(randomInt(97, 122))
        } else {
          bytes[i] = byte(randomInt(65, 90))
        }
    }
    return string(bytes)
}

func main() {
  http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
    urlParam, ok := r.URL.Query()["url"]
    var message string
    if urlParam != nil || ok {
      rand.Seed(time.Now().UnixNano())
      message = randomString(4)
    } else {
      message = "No URL present..."
    }
    fmt.Fprintf(w, html.EscapeString(message))
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}