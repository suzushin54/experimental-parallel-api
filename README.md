# Experimental Parallel API

## Project Overview
This project is designed to experiment with and evaluate the performance differences between serial and parallel processing within a payment service architecture.

## Infrastructure Layer Simulations
The `/infra` layer uses mock implementations without actual external API calls or database access.

## Example

```sh
    make bench
    
    cd internal/service && go test -bench . -benchmem -benchtime=10s -cpuprofile=cpu.prof -memprofile=mem.prof -blockprofile=block.prof
    goos: darwin
    goarch: arm64
    pkg: github.com/suzushin54/experimental-parallel-api/internal/service
    BenchmarkSerialPaymentService-10      	       6	1803190250 ns/op	   21313 B/op	     257 allocs/op
    BenchmarkParallelPaymentService-10    	       8	1302305573 ns/op	   19355 B/op	     264 allocs/op
    PASS
    ok  	github.com/suzushin54/experimental-parallel-api/internal/service	24.601s
    
    Process finished with exit code 0
```