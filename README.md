# go-anthropic
This is a golang SDK for anthropic claude API.

# Quick Start
To start working with `go-anthropic`, instantiate the `Anthropic` object with `sdk.NewAnthropic()`:
```go
key := "your-api-key"
anthropic, err := sdk.NewAnthropic(http.DefaultClient, key)
if err != nil {
    panic(err)
}
```
`Anthropic` object is then used to interact with the api:
* `Answer` allows user to ask questions while staying oblivious to the API structure, filling all the inferrable fields
and constructing/parsing the request/response automatically.
```go
completion, err := anthropic.Answer("Why is the sky blue?":, 255)
if err != nil {
    panic(err)
}
fmt.Println(*completion)
```

* `Do` performs an explicitly described `Request` on the api. This requires a measure of knowledge of anthropic's API
  structure (see https://console.anthropic.com/docs/api/reference).
```go
completion, err := anthropic.Do(sdk.Request{
    Prompt:            "\n\nHuman: Why is the sky blue?\n\nAssistant:",
    Model:             sdk.ModelClaude__V1,
    MaxTokensToSample: 255,
})
if err != nil {
    panic(err)
}
fmt.Println(*completion)
```
