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
	"flag"
	"fmt"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	defaultPort = "8080"
	maxRetries  = 3
)

type server struct {
	services map[string]string
}

func newServer(services map[string]string) *server {
	return &server{services}
}

//// Provides a storage client
//func (s server) getStorageClient() (storagepb.StorageClient, func(), error) {
//	conn, err := s.getConn(storageService)
//	if err != nil {
//		log.Print("ERROR: could not open connection to storage")
//		return nil, nil, err
//	}
//	return storagepb.NewStorageClient(conn), closeConnFn(conn), err
//}

// Provides a connection-closing function
func closeConnFn(conn *grpc.ClientConn) func() {
	return func() {
		if err := conn.Close(); err != nil {
			log.Printf("WARN: error closing connection: %v", err)
		}
	}
}

// Establishes a connection to a service
func (s server) getConn(service string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(s.services[service], grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	return conn, nil
}

func (s server) ClassifyComment(ctx context.Context, r *pb.ClassifyCommentRequest) (*pb.ClassifyResponse, error) {
	return &pb.ClassifyResponse{Category: int32(len(r.Comment.Text))}, nil
}

func main() {
	var port = flag.String("port", defaultPort, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterClassyServer(s, newServer(map[string]string{}))

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
