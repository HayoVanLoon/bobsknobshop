/*
 * Copyright 2019 Hayo van Loon
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
	"context"
	"flag"
	"fmt"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/peddler/v1"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = "8080"
)

func getConn(host, port string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	return conn, nil
}

func createMessage(host, port string, c string) error {
	r := &pb.SearchOrdersRequest{
		Customer: []string{c},
	}

	conn, err := getConn(host, port)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Panicf("error closing connection: %v", err)
		}
	}()

	cl := pb.NewPeddlerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := cl.SearchOrders(ctx, r)

	log.Printf("%v\n", resp)
	if err != nil {
		log.Printf("%v\n", err)
	}

	return err
}

// Makes several calls to a messaging server.
//
// Meant for debugging purposes.
func main() {
	var host = flag.String("host", defaultHost, "service host")
	var port = flag.String("port", defaultPort, "service port")
	flag.Parse()

	_ = createMessage(*host, *port, "Alice")
	_ = createMessage(*host, *port, "Cathy")
}
