# What's this ?

Provide a client supporting generalized invocation with a service registry(nacos) 


# Registry support 
- nacos 

# How to use it?

> unary invoke 

``` go 
	c := NewGenericClient()
	c.Debug = true
	c.Registryconfig = &NacosRegistryConfig{
		NamespaceId: "local",
		Port:        os.Getenv("NACOS_PORT"),
		IpAddr:      os.Getenv("NACOS_ADDR"),
	}
	r := c.GenericUnaryInvoke(context.Background(), InvokerPayload{
		Service:   "greet.Greet",
		Method:    "Ping",
		JsonParam: `{"ping":"OK"}`,
	})
    // Block until receive response from chan 
	response := <-r
```
