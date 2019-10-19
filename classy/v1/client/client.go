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
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	commonpb "github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
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

func createMessage(host, port string, c *commonpb.Comment) error {
	r := &pb.ClassifyCommentRequest{Comment: c}

	conn, err := getConn(host, port)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Panicf("error closing connection: %v", err)
		}
	}()

	cl := pb.NewClassyClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := cl.ClassifyComment(ctx, r)

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
	var host = flag.String("host", defaultHost, "messaging service host")
	var port = flag.String("port", defaultPort, "messaging service port")
	flag.Parse()

	question := &commonpb.Comment{
		Name:      "id-123",
		Text:      "I have a question about this product. How do I use it?",
		Author:    "Cathy",
		CreatedOn: time.Now().UnixNano() / 1000000,
	}
	complaint := &commonpb.Comment{
		Name:      "id-456",
		Text:      "The knob is too jolly. This does not please me.",
		Author:    "Alice",
		CreatedOn: time.Now().UnixNano()/1000000 + 1,
	}

	_ = createMessage(*host, *port, question)
	_ = createMessage(*host, *port, complaint)
}