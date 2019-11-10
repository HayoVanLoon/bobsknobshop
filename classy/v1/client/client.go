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
	"strconv"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = 9000
)

func classifyComment(host string, port int, c *commonpb.Comment) {
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer func() { _ = conn.Close() }()

	cl := pb.NewClassyClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := cl.ClassifyComment(ctx, c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", resp)
}

func listClassifications(host string, port int) {
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer func() { _ = conn.Close() }()

	cl := pb.NewClassyClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := cl.ListClassifications(ctx, &pb.ListClassificationsRequest{})
	if err != nil {
		panic(err)
	}
	for _, c := range resp.Classifications {
		fmt.Println(c)
	}
}

// Meant for debugging purposes.
func main() {
	var host = flag.String("host", defaultHost, "service host")
	var port = flag.Int("port", defaultPort, "service port")
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

	classifyComment(*host, *port, question)
	classifyComment(*host, *port, complaint)
	listClassifications(*host, *port)
}
