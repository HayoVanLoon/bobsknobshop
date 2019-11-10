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
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/peddler/v1"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultPort = 9000
)

func createOrders(fileName string) []common.Order {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	r := bytes.NewReader(bs)
	dec := json.NewDecoder(r)

	var os []common.Order

	for dec.More() {
		o := common.Order{}
		err = jsonpb.UnmarshalNext(dec, &o)
		if err != nil {
			panic(err.Error())
		}
		os = append(os, o)
	}

	return os
}

type server struct {
	orders []common.Order
}

func newServer() *server {
	return &server{}
}

func (s server) SearchOrders(ctx context.Context, r *pb.SearchOrdersRequest) (*pb.SearchOrdersResponse, error) {
	cs := r.GetCustomer()
	sort.Slice(cs, func(i, j int) bool { return cs[i] < cs[j] })
	ps := r.GetProductId()
	sort.Slice(ps, func(i, j int) bool { return ps[i] < ps[j] })

	var os []*common.Order
	for i, o := range s.orders {
		added := false
		oc := strings.ToUpper(o.Client)
		for _, c := range cs {
			if oc == strings.ToUpper(c) {
				added = true
				os = append(os, &s.orders[i])
				break
			}
		}
		if added {
			continue
		}
		for _, p := range ps {
			p = strings.ToUpper(p)
			for _, ol := range o.GetOrderLines() {
				if p == strings.ToUpper(ol.Sku) {
					added = true
					os = append(os, &s.orders[i])
					break
				}
			}
			if added {
				break
			}
		}
	}

	resp := &pb.SearchOrdersResponse{Orders: os}
	return resp, nil
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")

	var fileName = flag.String("file", "comments.csv", "comment file")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPeddlerServer(s, server{createOrders(*fileName)})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
