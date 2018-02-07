package main

import (
  "github.com/Kooooya/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("Hello World - Traffic")
}

func main() {
  router := traffic.New()
  router.Get("/", rootHandler)
  router.Run()
}
