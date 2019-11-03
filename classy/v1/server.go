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
	"encoding/csv"
	"flag"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	commonpb "github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	defaultPort     = 9000
	defaultSubsFile = "subs.csv"
)

type server struct {
	services map[string]string
}

func readSubsFile(file string) map[string]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("could not open config file")
	}
	r := csv.NewReader(f)

	vs := map[string]string{}
	for row, err := r.Read(); row != nil; row, err = r.Read() {
		if err != nil {
			log.Fatal("error reading config file")
		}
		if len(row) != 3 {
			log.Fatalf("invalid row in config file %s", row)
		}
		vs[row[0]] = vs[row[1]] + ":" + row[2]
	}

	return vs
}

// Provides a connection-closing function
func CloseConnFn(conn *grpc.ClientConn) func() {
	return func() {
		if err := conn.Close(); err != nil {
			log.Printf("WARN: error closing connection: %v", err)
		}
	}
}

// Establishes a connection to a service
func GetConn(s string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(s, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Provides a sub-service  client
func (s server) getSubServiceClient(sub string) (pb.ClassyClient, func(), error) {
	conn, err := GetConn(s.services[sub])
	if err != nil {
		return nil, nil, err
	}
	return pb.NewClassyClient(conn), CloseConnFn(conn), err
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
	var subsFile = flag.String("subs-file", defaultSubsFile, "csv file with version parameters")
	var port = flag.Int("port", defaultPort, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterClassyServer(s, server{readSubsFile(*subsFile)})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
