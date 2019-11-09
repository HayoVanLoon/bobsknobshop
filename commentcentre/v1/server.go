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
	classypb "github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/commentcentre/v1"
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	defaultPort       = 9000
	defaultClassyHost = "classy-service"
)

type server struct {
	classyService string
	comments      []*common.Comment
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
func getClassyClient(service string) (classypb.ClassyClient, func(), error) {
	conn, err := GetConn(service)
	if err != nil {
		return nil, nil, err
	}
	return classypb.NewClassyClient(conn), CloseConnFn(conn), err
}

func (s *server) CreateComment(ctx context.Context, r *common.Comment) (*common.Comment, error) {
	cl, closeFn, err := getClassyClient(s.classyService)
	if err != nil {
		log.Print("error creating classy client")
		return nil, err
	}
	defer closeFn()

	r.Name = "commentcentre.bobsknobshop.gl/comments/" + uuid.New().String()
	r.CreatedOn = time.Now().Unix()

	// TODO: store comment
	s.comments = append(s.comments, r)

	resp, err := cl.ClassifyComment(ctx, r)
	if err != nil {
		log.Print("error calling classy service")
		return nil, err
	}
	log.Printf("%s", resp)

	return r, nil
}

func (s *server) ListComments(ctx context.Context, r *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	resp := &pb.ListCommentsResponse{Comments: s.comments}
	return resp, nil
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")

	var classyHost = flag.String("classy-host", defaultClassyHost, "classy service host")
	var classyPort = flag.Int("classy-port", defaultPort, "classy service port")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	classyService := *classyHost + ":" + strconv.Itoa(*classyPort)
	pb.RegisterCommentcentreServer(s, &server{classyService, []*common.Comment{}})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
