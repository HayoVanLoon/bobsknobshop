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
	"github.com/HayoVanLoon/bobsknobshop/common"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	commonpb "github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

const (
	defaultPort = 9000
	self        = "classy"
	version     = "v1"
)

var implementations = map[string]common.ServiceImplementation{}

type server struct {
	services map[string]string
}

func newServer(services map[string]string) *server {
	return &server{services}
}

// Provides a sub-service  client
func (s server) getSubServiceClient(sub string) (pb.ClassyClient, func(), error) {
	conn, err := s.getConn(implementations[sub].Service())
	if err != nil {
		return nil, nil, err
	}
	return pb.NewClassyClient(conn), closeConnFn(conn), err
}

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
		return nil, err
	}
	return conn, nil
}

func (s server) ClassifyComment(ctx context.Context, r *commonpb.Comment) (resp *pb.Classification, err error) {

	cl, cls, err := s.getSubServiceClient("a3nlp")
	if err != nil {
		return nil, err
	}
	defer cls()

	resp, err = cl.ClassifyComment(ctx, r)
	if err != nil {
		return nil, err
	}

	// TODO: store response
	// TODO: store metadata

	return resp, nil
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")
	flag.Parse()

	for i, n := range []string{"a1basic", "a2extradata", "a3nlp"} {
		implementations[n] = common.NewServiceImplementation(self, version, n, *port+1+i, true)
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
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
