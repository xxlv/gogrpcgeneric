package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc/reflection"

	"context"
	"net"

	"test/example"

	"google.golang.org/grpc"
)

const (
	port = ":5051"
)

type server struct{}

func (s *server) Greet(ctx context.Context, in *example.HelloRequest) (*example.HelloResponse, error) {
	return &example.HelloResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	example.RegisterGreeterServer(grpcServer, &server{})

	nacosHost := "127.0.0.1"
	nacosPort := uint64(8848)
	serviceName := "demo-grpc-service"
	ip := "127.0.0.1"
	port := uint64(5051)
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosHost,
			Port:   nacosPort,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		UpdateThreadNum:     20,
		NotLoadCacheAtStart: true,
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		log.Fatal("create nacos client failed ", err)
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "host.docker.internal",
		Port:        port,
		ServiceName: serviceName,
		Healthy:     true,
		Enable:      true,
		Weight:      100,
	})

	if err != nil {
		log.Fatal("Nacos failed ", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal("can not start gRPC ", err)
	}

	go func() {
		log.Println("run gRPC ...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC error: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan

	log.Println("close ...")
	grpcServer.GracefulStop()

	_, err = client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
	})

	if err != nil {
		log.Println("Nacos failed>:", err)
	}

}
