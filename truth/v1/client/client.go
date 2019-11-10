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
	"fmt"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/truth/v1"
	"google.golang.org/grpc"
	"sort"
	"strconv"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = 9000
)

func getClassyKpis(host string, port int) {
	r := &pb.GetServiceKpiRequest{Name: "classy"}
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer func() { _ = conn.Close() }()
	cl := pb.NewTruthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := cl.GetServiceKpi(ctx, r)
	if err != nil {
		fmt.Println(err)
	}

	sort.Slice(resp.Versions, func(i, j int) bool { return resp.Versions[i].GetName() < resp.Versions[j].GetName() })
	for _, v := range resp.Versions {
		fmt.Printf("%s: %v\n", v.GetName(), v.GetValue())
	}
}

// Meant for debugging purposes.
func main() {
	var host = flag.String("host", defaultHost, "service host")
	var port = flag.Int("port", defaultPort, "service port")

	flag.Parse()

	getClassyKpis(*host, *port)
}
