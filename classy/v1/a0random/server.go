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
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
	"strconv"
)

const (
	defaultPort = 9000
)

type server struct {
}

func (s *server) ClassifyComment(ctx context.Context, r *common.Comment) (*pb.Classification, error) {
	cat := predict()

	resp := &pb.Classification{Category: cat}
	return resp, nil
}

// Predict the comment's category
func predict() string {
	opts := [5]string{"complaint", "compliment", "question", "review", "undetermined"}
	return opts[rand.Intn(len(opts))]
}

func (s *server) ListClassifications(context.Context, *pb.ListClassificationsRequest) (*pb.ListClassificationsResponse, error) {
	return nil, fmt.Errorf("not implemented (by design)")
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterClassyServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
