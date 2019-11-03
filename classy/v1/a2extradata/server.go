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
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	commonpb "github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"github.com/HayoVanLoon/genproto/bobsknobshop/peddler/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"unicode"
)

const (
	defaultPort = "9000"
)

type server struct {
}

func getClient(s string) (peddler.PeddlerClient, func(), error) {
	conn, err := GetConn(s)
	if err != nil {
		return nil, nil, err
	}
	closeConnFn := CloseConnFn(conn)
	cl := peddler.NewPeddlerClient(conn)
	return cl, closeConnFn, nil
}

func (s *server) ClassifyComment(ctx context.Context, r *commonpb.Comment) (*pb.Classification, error) {
	qc := 0
	ec := 0
	emo := 0
	lst := '.'
	for _, c := range r.Text {
		if c == '?' {
			qc += 1
		} else if c == '!' {
			ec += 1
		} else if unicode.IsUpper(c) && !unicode.IsPunct(lst) {
			emo += 1
		}

		if !unicode.IsSpace(c) {
			lst = c
		}
	}

	cl, closeFn, err := getClient("s")
	if err != nil {
		// TODO: handle error
		return nil, err
	}
	defer closeFn()

	log.Print(cl)

	resp := &pb.Classification{}
	if ec > 0 && emo < 2 {
		resp.Category = "compliment"
	} else if emo > 2 {
		resp.Category = "complaint"
	} else if qc > 0 {
		resp.Category = "question"
	} else {
		resp.Category = "comment"
	}

	return resp, nil
}

func main() {
	var port = flag.String("port", defaultPort, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
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
