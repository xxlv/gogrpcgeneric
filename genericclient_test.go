package gogrpcgeneric

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
)

func TestGenericInvoke(t *testing.T) {
	_ = godotenv.Load(".env")
	c := NewGenericClient()
	c.Debug = true
	c.Registryconfig = &NacosRegistryConfig{
		NamespaceId: "public",
		Port:        8848,
		IpAddr:      "127.0.0.1",
	}

	r := c.GenericUnaryInvoke(context.Background(), "demo-grpc-service", InvokerPayload{
		Service:   "example.Greeter",
		Method:    "Greet",
		JsonParam: `{"name":"hello world"}`,
	})
	response := <-r
	if response.Result != "{\n  \"message\": \"Hello hello world\"\n}" {
		t.Error("not pass")
	}
}

func BenchmarkLoadServiceDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LoadServiceDefault("demo-grpc-service", "")
	}
}
