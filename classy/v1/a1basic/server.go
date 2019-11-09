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
	"github.com/HayoVanLoon/go-commons/i18n"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
	"unicode"
)

const (
	defaultPort = 9000
)

type server struct {
	peddlerService string
}

func (s server) ClassifyComment(ctx context.Context, r *common.Comment) (*pb.Classification, error) {
	q, ex, emo := analyseText(r.GetText())

	cat := predict(q, ex, emo, len(r.GetText()))

	resp := &pb.Classification{Category: cat}
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

// Predict the comment's category
func predict(questionMarks, exclamationMarks int, emo float32, l int) string {
	log.Printf("%v; %v; %v; %v", questionMarks, exclamationMarks, emo, l)

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
			if l > 40 {
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
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterClassyServer(s, server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
