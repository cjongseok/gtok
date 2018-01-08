gtok
====
Get a Google Cloud Platform access token by your GCP service account key.

Usage
----
Use it in your Go program,
```go
account, err := gtok.NewGcpServiceAccountFromFile("key.json") // Read your service account key file
tok, err := account.AccessToken()   // You can get a token in string
```
Or you can get a token instantly,
```sh
# In this repository home,
go run main/main.go <your_json_key_file> 

# Its output would be 'Bearer xxxxx...'
```

Dependencies
---
* github.com/dgrijalva/jwt-go
