# URL Shortener project

After running the server you an use any http client to test the server example: 

```sh
curl -X POST -H "Content-Type: application/json" -v "http://localhost:8080/url" -d '{ "url": "https://google.com", "alias": "go2" }
```
