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
	"context"
	"flag"
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/peddler/v1"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort = 9000
	self        = "peddler"
	version     = "v1"
)

func amount(a int64, b int32) *money.Money {
	return &money.Money{CurrencyCode: "EUR", Units: a, Nanos: b}
}

func orderLine(sku string, q int32, t *money.Money) *common.Order_OrderLine {
	return &common.Order_OrderLine{
		Name:     uuid.New().String(),
		Sku:      sku,
		Quantity: q,
		Total:    t,
	}
}

func order(c string, m *money.Money, os []*common.Order_OrderLine) common.Order {
	return common.Order{
		Name:       uuid.New().String(),
		CreatedOn:  time.Now().Unix(),
		Client:     c,
		Total:      m,
		OrderLines: os,
	}
}

func createOrders() []common.Order {
	orders := []common.Order{
		order("Alice", amount(5, 0), []*common.Order_OrderLine{
			orderLine("123-456-789-0-1", 1, amount(6, 0)),
		}),
		order("Alice", amount(5, 0), []*common.Order_OrderLine{
			orderLine("123-456-789-0-1", 1, amount(6, 0)),
		}),
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Name < orders[j].Name
	})
	return orders
}

var orders = createOrders()

type server struct {
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
	for _, o := range orders {
		added := false
		oc := strings.ToUpper(o.Client)
		for _, c := range cs {
			if oc == strings.ToUpper(c) {
				added = true
				os = append(os, &o)
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
					os = append(os, &o)
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
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPeddlerServer(s, newServer())

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
