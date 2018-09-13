/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"log"
	"flag"
	"time"
    "io"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "helloworld"
)

const (
	address     = "172.30.59.27:50051"
	defaultName = "world"
)

var (
    tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
    caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
    serverAddr         = flag.String("server_addr", "127.0.0.1:50051", "The server address in the format of host:port")
    serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)


func main() {
	flag.Parse()
    // Set up a connection to the server.
	/*
    conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
    */
    notes := []*pb.HelloRequest{
        {Name: "aaa"},
        {Name: "bbb"},
        {Name: "ccc"},
        {Name: "ddd"},
        {Name: "eee"},
    }
    ctx := context.TODO()
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithInsecure())
    conn, err := grpc.Dial(*serverAddr, opts...)
    if err != nil {
        log.Fatalf("fail to dial: %v", err)
    }
    defer conn.Close()
    client := pb.NewGreeterClient(conn)

    stream, err := client.SayHello(ctx)
    if err != nil {
        log.Fatalf("%v.SayHello(_) = _, %v", client, err)
    }
    waitc := make(chan struct{})
    go func() {
        for {
            in, err := stream.Recv()
            if err == io.EOF {
                // read done.
                close(waitc)
                return
            }
            if err != nil {
                log.Fatalf("Failed to receive a note : %v", err)
            }
            log.Printf("Got message %s", in.Message)
        }
    }()
    for _, note := range notes {
        if err := stream.Send(note); err != nil {
            log.Fatalf("Failed to send a note: %v", err)
        }
    }
    time.Sleep(600 * time.Second)
    stream.CloseSend()
    <-waitc
}
