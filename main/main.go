package main

import (
  "fmt"
  "github.com/cjongseok/gtok"
  "os"
)

func usage() {
  fmt.Println("gtok <json_key_file>")
  fmt.Println()
  fmt.Println("Get Google Cloud Platform access token.")
  fmt.Println()
  fmt.Println("  json_key_file    means service account json key file you downloaded when you create")
  fmt.Println("                   your service account. If you don't have a service account yet, go to")
  fmt.Println("                   https://console.cloud.google.com/apis/credentials/serviceaccountkey")

}

func main() {
  os.Exit(getGtok())
}

func getGtok() int {
  if len(os.Args) != 2 {
    usage()
    return 128
  }
  account, err := gtok.NewGcpServiceAccountFromFile(os.Args[1])
  if err != nil {
    fmt.Println(err)
    return 1
  }
  tok, err := account.AccessToken()
  if err != nil {
    fmt.Println(err)
    return 1
  }
  fmt.Println(tok)
  return 0
}
