__This project is in its early stage, and the API is likely to change and not stable.__

# WebAssembly for Proxies (Go SDK) [![Build](https://github.com/tetratelabs/proxy-wasm-go-sdk/workflows/test/badge.svg)](https://github.com/tetratelabs/proxy-wasm-go-sdk/actions) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

The Go SDK for
 [Proxy-Wasm](https://github.com/proxy-wasm/spec), enabling developers to write Proxy-Wasm extensions in Go.

This SDK is powered by [TinyGo](https://tinygo.org/) and does not support the official Go compiler.

```golang
import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var counter proxywasm.MetricCounter

type metricRootContext struct { types.DefaultRootContext }

func (ctx *metricRootContext) OnVMStart(int) types.OnVMStartStatus {
	// Initialize the metric.
	counter = proxywasm.DefineCounterMetric("proxy_wasm_go.request_counter")
	return types.OnVMStartStatusOK
}

func (*metricRootContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &metricHttpContext{}
}

type metricHttpContext struct { types.DefaultHttpContext }

func (ctx *metricHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	// Increment the request counter when we receive request headers.
	counter.Increment(1)
	return types.ActionContinue
}
```

## Requirements

This SDK depends on TinyGo and leverages its [WASI](https://github.com/WebAssembly/WASI) (WebAssembly System Interface) target.

Please follow the official instruction [here](https://tinygo.org/getting-started/) for installing TinyGo.

### Compatible ABI / Envoy builds (tested on CI)

| proxy-wasm-go-sdk| proxy-wasm ABI version |istio/proxyv2| Envoy upstream|
|:-------------:|:-------------:|:-------------:|:-------------:|
| main | 0.2.0| 1.9, 1.10 | 1.18 |
| v0.2.0 | 0.2.0| 1.8, 1.9 | 1.17 |

## Run examples

build:

```bash
# Build all examples.
make build.examples

# Build a specific example.
make build.example name=helloworld
```

run:

```bash
# This requires you to have Envoy binary locally.
make run name=helloworld
``` 

## SDK development

```bash
# Run local tests without running envoy processes.
make test

# Run all e2e tests.
# This requires you to have Envoy binary locally.
make test.e2e

# Run e2e tests for a specific example.
# This requires you to have Envoy binary locally.
make test.e2e.single name=helloworld
```

## Limitations and Considerations

- Some of existing libraries are not available (importable but runtime panic / non-importable)
    - There are several reasons for this:
        1. TinyGo's WASI target does not support some of syscall: For example, we cannot import `crypto/rand` package.
        2. TinyGo does not implement all of reflect package([examples](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/reflect/value.go#L299-L305)).
        3. [proxy-wasm-cpp-host](https://github.com/proxy-wasm/proxy-wasm-cpp-host) has not supported some of WASI APIs yet 
    - These issues will be mitigated as TinyGo and proxy-wasm-cpp-host evolve.
- There's performance overhead of using Go/TinyGo due to GC
    - `runtime.GC` is called whenever the heap runs out (see [1](https://tinygo.org/lang-support/#garbage-collection),
    [2](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/runtime/gc_conservative.go#L218-L239)).
    - TinyGo allows us to disable GC, but we cannot do that since we need to use maps (implicitly causes allocation)
     for saving the plugin's [state](https://github.com/tetratelabs/proxy-wasm-go-sdk/blob/cf6ad74ed58b284d3d8ceeb8c5dba2280d5b1007/proxywasm/vmstate.go#L41-L46).
    - Theoretically, we can implement our own GC algorithms tailored for proxy-wasm through `alloc(uintptr)` [interface](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/runtime/gc_none.go#L13) 
    with `-gc=none` option. This is the future TODO.
- `recover` is [not implemented](https://github.com/tinygo-org/tinygo/issues/891) in TinyGo, and there's no way to prevent the Wasm virtual machine from aborting.
- Goroutine support
    - In TinyGo, Goroutine is implmeneted through LLVM's coroutine (see [this blog post](https://aykevl.nl/2019/02/tinygo-goroutines)).
    - In Envoy, Wasm modules are run in the event driven manner, and therefore the "scheduler" is not executed once the main function exits.
        That means you cannot have the expected behavior of Goroutine as in ordinary host environments.
        - The question "How to deal with Goroutine in a thread local Wasm VM executed in the event drive manner" has yet to be answered.
    - We strongly recommend that you implement the `OnTick` function for any asynchronous task instead of using Goroutine.
    - The scheduler can be disabled with `-scheduler=none` option of TinyGo.

## References

- [WebAssembly for Proxies (ABI specification)](https://github.com/proxy-wasm/spec)
- [WebAssembly for Proxies (C++ SDK)](https://github.com/proxy-wasm/proxy-wasm-cpp-sdk)
- [WebAssembly for Proxies (Rust SDK)](https://github.com/proxy-wasm/proxy-wasm-rust-sdk)
- [TinyGo - Go compiler for small places](https://tinygo.org/)
