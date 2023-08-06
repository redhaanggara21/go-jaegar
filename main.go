package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 10,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831", // replace host
		},
	}

	closer, err := cfg.InitGlobalTracer(
		"main-service",
	)
	defer closer.Close()

	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}

	http.HandleFunc("/get-product", handleGetProduct)

	fmt.Println("START SERVICE IN PORT 8090")
	http.ListenAndServe(":8090", nil)
}

func handleGetProduct(w http.ResponseWriter, req *http.Request) {
	trace, ctx := opentracing.StartSpanFromContext(req.Context(), "Handle /get-product")
	time.Sleep(time.Second / 2)
	defer trace.Finish()

	// check login
	if isLogin(ctx) {
		fmt.Println("USER LOGIN")
	}

	// get product
	product := getProduct(ctx)
	fmt.Println("product: ", product)
}

func isLogin(ctx context.Context) bool {
	var status bool
	trace, ctx := opentracing.StartSpanFromContext(ctx, "func isLogin")
	time.Sleep(time.Second / 2)
	defer trace.Finish()

	// get status
	status = true

	return status
}

func getProduct(ctx context.Context) map[string]string {

	trace, ctx := opentracing.StartSpanFromContext(ctx, "func getProduct")
	time.Sleep(time.Second / 2)
	defer trace.Finish()

	// get product
	product := map[string]string{"id": "1", "name": "Laptop", "warna": "hitam"}

	return product
}
