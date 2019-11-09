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
	commonpb "github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"github.com/HayoVanLoon/genproto/bobsknobshop/peddler/v1"
	"github.com/HayoVanLoon/go-commons/i18n"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
	"strings"
	"unicode"
)

const (
	defaultPort        = 9000
	defaultPeddlerHost = "peddlerService-service"
)

type server struct {
	peddlerService string
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

func getClient(s string) (peddler.PeddlerClient, func(), error) {
	conn, err := GetConn(s)
	if err != nil {
		return nil, nil, err
	}
	closeConnFn := CloseConnFn(conn)
	cl := peddler.NewPeddlerClient(conn)
	return cl, closeConnFn, nil
}

func (s server) ClassifyComment(ctx context.Context, r *commonpb.Comment) (*pb.Classification, error) {
	q, ex, emo := analyseText(r.Text)

	o, err := s.hasOrdered(ctx, r.GetAuthor(), r.GetTopic())
	if err != nil {
		log.Printf("error fetching orders, assuming none exist (%s)", err.Error())
	}

	resp := &pb.Classification{Category: predict(q, ex, emo, o)}
	return resp, nil
}

// Extracts features from the given text
func analyseText(s string) (q, ex int, emo float32) {
	lst := '.'
	for _, c := range s {
		if i18n.IsQuestionMark(c) {
			q += 1
		} else if i18n.IsExclamationMark(c) {
			ex += 1
		} else if unicode.IsUpper(c) && !unicode.IsPunct(lst) {
			emo += 1
		}

		if !unicode.IsSpace(c) {
			lst = c
		}
	}
	emo = emo / float32(len(s))
	return
}

// Checks if the customer has ordered a certain product
func (s server) hasOrdered(ctx context.Context, cust, sku string) (bool, error) {
	cl, closeFn, err := getClient(s.peddlerService)
	if err != nil {
		log.Printf("error creating client for %s", s.peddlerService)
		return false, err
	}
	defer closeFn()

	resp, err := cl.SearchOrders(ctx, &peddler.SearchOrdersRequest{Customer: []string{cust}})
	if err != nil {
		log.Print("error fetching orders")
		return false, err
	}

	upperSku := strings.ToUpper(sku)
	for _, o := range resp.GetOrders() {
		for _, ol := range o.GetOrderLines() {
			if strings.ToUpper(ol.GetSku()) == upperSku {
				return true, nil
			}
		}
	}

	return false, nil
}

// Predict the comment's category
func predict(questionMarks, exclamationMarks int, emo float32, ordered bool) string {
	log.Printf("%v; %v; %v; %v", questionMarks, exclamationMarks, emo, ordered)

	if exclamationMarks > 0 {
		if emo > 0.1 {
			return "complaint"
		} else {
			return "compliment"
		}
	} else {
		if questionMarks > 0 {
			return "question"
		} else {
			if ordered {
				return "review"
			} else {
				return "undetermined"
			}
		}
	}
}

func (s server) ListClassifications(context.Context, *pb.ListClassificationsRequest) (*pb.ListClassificationsResponse, error) {
	return nil, fmt.Errorf("not implemented (by design)")
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")

	var peddlerHost = flag.String("peddler-host", defaultPeddlerHost, "peddler service host")
	var peddlerPort = flag.Int("peddler-port", defaultPort, "peddler service port")

	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	peddlerService := *peddlerHost + ":" + strconv.Itoa(*peddlerPort)
	pb.RegisterClassyServer(s, server{peddlerService})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
