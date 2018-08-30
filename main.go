package main

import (
  "fmt"
  "net/http"
  "html"
  "log"
  "math/rand"
  "time"
  "context"

  "github.com/mongodb/mongo-go-driver/mongo"
)

type Resource struct {
  url string
  shortened_url string
}

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

func shortenHandler(w http.ResponseWriter, r *http.Request) {
  urlParam, ok := r.URL.Query()["url"]
  var message string
  if urlParam != nil || ok {
    rand.Seed(time.Now().UnixNano())
    message = randomString(4)
  } else {
    message = "No URL present..."
  }
  fmt.Fprintf(w, html.EscapeString(message))
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Not Yet Implemented") 
}

func main() {
  client, err := mongo.NewClient("mongodb://database:27017")
  if err != nil { log.Fatal(err) }
  err = client.Connect(context.TODO())
  if err != nil { log.Fatal(err) }
  http.HandleFunc("/shorten", shortenHandler)
  http.HandleFunc("/expand", expandHandler)

  log.Fatal(http.ListenAndServe(":8080", nil))
}