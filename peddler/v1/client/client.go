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
	"strconv"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = 9000
)

func getConn(host string, port int) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	return conn, nil
}

func createMessage(host string, port int, c string, ps []string) error {
	r := &pb.SearchOrdersRequest{
		Customer:  []string{c},
		ProductId: ps,
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

// Meant for debugging purposes.
func main() {
	var host = flag.String("host", defaultHost, "service host")
	var port = flag.Int("port", defaultPort, "service port")
	flag.Parse()

	fmt.Print("\n")
	_ = createMessage(*host, *port, "Alice", []string{})
	fmt.Print("\n")
	_ = createMessage(*host, *port, "Cathy", []string{})
	fmt.Print("\n")
	_ = createMessage(*host, *port, "", []string{"123-456-789-0-1"})
	fmt.Print("\n")
}
