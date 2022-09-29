# protoc-gen-go-opensearch

Plugin used to generate mappings for OpenSearch indexes based on `protojson`.
WIP

## TODO

[ ] A way to generate mappings for external protos
[ ] Add options to match options from `protojson`

It is highly based of the code of [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go/tree/b92717ecb630d4a4824b372bf98c729d87311a4d/cmd/protoc-gen-go).
I just got the code, stripped out any unnecessary parts, and built this plugin.

# Usage

Install it using the `go install` command:
> go install github.com/Hellysonrp/protoc-gen-go-opensearch

Some usage examples:
> protoc --plugin protoc-gen-go-opensearch --go-opensearch_out=output example.proto

> protoc --plugin protoc-gen-go-opensearch --go-opensearch_out=paths=source_relative:output example.proto

If you have problems with `protoc` not finding the plugin in the `PATH`, I recommend passing the absolute path to the plugin:
> protoc --plugin ${HOME}/go/bin/protoc-gen-go-opensearch --go-opensearch_out=output example.proto

> protoc --plugin ${HOME}/go/bin/protoc-gen-go-opensearch --go-opensearch_out=paths=source_relative:output example.proto

Then you can use it to get the mappings to send to OpenSearch:

```go
mappings := types.OpensearchMapping{
    Properties: (&Example{}).GetOpensearchMappings(),
}

b, err := json.Marshal(mappings)
if err != nil {
    panic(err)
}

buf := bytes.NewBuffer(b)

...

// this is only an example
// you may want to get the response and error
opensearchapi.SearchRequest{
    Index: []string{index},
    Body:  buf,
}.Do(ctx, client)
```
