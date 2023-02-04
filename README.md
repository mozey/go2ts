# go2ts

Code for `Convert` func copied from [StirlingMarketingGroup/go2ts](https://github.com/StirlingMarketingGroup/go2ts). This repo can be imported as a package for generating TypeScript from Go as part of a build workflow. Supports Go types `StructType` and `ArrayType` 


## Testing

Run all tests
```bash
git clone https://github.com/mozey/go2ts.git
cd go2ts
go test -v ./...
```


## Demo

Run demo with [deno](https://deno.land/). 

Using TypeScript types that were generated from Go code
```bash
deno run testdata/example/demo/app.ts
```

Type usage is statically checked
```bash
deno run testdata/example/demo/err.ts
```
