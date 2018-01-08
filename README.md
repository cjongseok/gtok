gtok
====
Get Google Cloud Platform access token by your GCP service account key.

Usage
----
You can use it as a part of your program,
```go
account, err := gtok.NewGcpServiceAccountFromFile("key.json") // Read your service account key file
tok, err := account.AccessToken()   // You can get a token in string
```
Or you can get a token instantly,
```sh
cd main
go run main.go <your_json_key_file> // Its output would be 'Bearer xxxxx...'
```
