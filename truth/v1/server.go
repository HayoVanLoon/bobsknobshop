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
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/truth/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	defaultPort = 9000
	self        = "truth"
	version     = "v1"
)

type server struct {
}

func (s server) GetServiceKpi(context.Context, *pb.GetServiceKpiRequest) (*pb.GetServiceKpiResponse, error) {
	resp := &pb.GetServiceKpiResponse{
		Versions: []*pb.GetServiceKpiResponse_Version{{
			Name:           "all",
			StartTimestamp: time.Now().Unix() - 24*time.Hour.Milliseconds(),
			EndTimestamp:   time.Now().Unix(),
			Unit:           "cabages",
			Value:          2,
		}},
	}
	return resp, nil
}

func newServer() *server {
	return &server{}
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTruthServer(s, newServer())

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
