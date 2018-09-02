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
  "github.com/mongodb/mongo-go-driver/bson"
)

type Resource struct {
  URL string          `bson:"url"`
  ShortenedUrl string `bson:"shortened_url"`
}

type Connections struct {
  Collection *mongo.Collection
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

func (c *Connections) shortenHandler(w http.ResponseWriter, r *http.Request) {
  urlParam, ok := r.URL.Query()["url"]
  var message string
  if urlParam != nil || ok {
    for {
      rand.Seed(time.Now().UnixNano())
      message = randomString(4)
      res := c.Collection.FindOne(nil, bson.NewDocument(bson.EC.String("shortened_url", message)))
      resource := Resource{}
      res.Decode(resource)

      if resource.ShortenedUrl == "" {
        insertResource := &Resource{URL: urlParam[0], ShortenedUrl: message}
        _, err := c.Collection.InsertOne(nil, insertResource)

        if err != nil {
          message = err.Error()
        }
        break
      }
    }
  } else {
    message = "No URL present..."
  }
  fmt.Fprintf(w, html.EscapeString(message))
}

func (c *Connections) expandHandler(w http.ResponseWriter, r *http.Request) {
  shortParam, ok := r.URL.Query()["short"]
  var output string
  if shortParam != nil || ok {
    res := c.Collection.FindOne(nil, bson.NewDocument(bson.EC.String("url", shortParam[0])))
    resource := &Resource{}
    res.Decode(resource)

    if resource.URL != "" {
      output = resource.URL
    } else {
      output = "Failed to find a URL for that short"
    }
  } else {
    output = "No short present..."
  }
  fmt.Fprintf(w, html.EscapeString(output))
}

func main() {
  client, err := mongo.NewClient("mongodb://database:27017")
  if err != nil { log.Fatal(err) }
  err = client.Connect(context.TODO())
  if err != nil { log.Fatal(err) }
  c := Connections{Collection: client.Database("urlShortener").Collection("resources")}
  http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
    c.shortenHandler(w, r)
  })
  http.HandleFunc("/expand", func(w http.ResponseWriter, r *http.Request) {
    c.expandHandler(w, r)
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}