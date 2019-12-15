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
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	defaultPort     = 9000
	defaultSubsFile = "subs.csv"
)

type versionConfig struct {
	name, target string
	traffic      int
}

func readSubsFile(file string) []versionConfig {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("could not open subs config file")
	}
	r := csv.NewReader(f)

	var vs []versionConfig
	for row, err := r.Read(); row != nil; row, err = r.Read() {
		if err != nil {
			log.Fatal("error reading subs config file")
		}
		if len(row) != 4 {
			log.Fatalf("invalid row in subs config file %s", row)
		}
		var ta string
		if row[2] != "" {
			ta = row[1] + ":" + row[2]
		} else {
			ta = row[1] + ":" + strconv.Itoa(defaultPort)
		}
		tr, err := strconv.Atoi(row[3])
		if err != nil {
			log.Fatalf("invalid traffic value in row: %s", row)
		}
		vs = append(vs, versionConfig{row[0], ta, tr})
	}

	if len(vs) == 0 {
		log.Fatal("empty subs config file")
	}

	sort.Slice(vs, func(i, j int) bool { return vs[i].traffic > vs[j].traffic })
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
func GetConn(t string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(t, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type server struct {
	services        []versionConfig
	classifications map[string]*pb.Classification
}

func (s *server) ClassifyComment(ctx context.Context, r *commonpb.Comment) (resp *pb.Classification, err error) {
	v := s.selectVersion()
	cl, cls, err := s.getSubServiceClient(v)
	if err != nil {
		return nil, err
	}
	defer cls()

	resp, err = cl.ClassifyComment(ctx, r)
	if err != nil {
		return nil, err
	}

	resp.Name = "classy.bobsknobshop.gl/classifications/" + uuid.New().String()
	resp.CreatedOn = time.Now().Unix()
	resp.ServiceVersion = v.name
	resp.Comment = r.GetName()

	// TODO: store in separate system, this breaks in multi-instance settings
	s.classifications[resp.Name] = resp

	return resp, nil
}

// Randomly selects version based on traffic weights
func (s server) selectVersion() versionConfig {
	acc := 0
	for _, v := range s.services {
		acc += v.traffic
	}
	r := rand.Intn(acc) + 1
	acc = 0
	for _, v := range s.services {
		if v.traffic+acc >= r {
			return v
		}
		acc += v.traffic
	}
	panic("unreachable code reached")
}

// Provides a sub-service  client
func (s server) getSubServiceClient(v versionConfig) (pb.ClassyClient, func(), error) {
	conn, err := GetConn(v.target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewClassyClient(conn), CloseConnFn(conn), err
}

func (s *server) ListClassifications(context.Context, *pb.ListClassificationsRequest) (*pb.ListClassificationsResponse, error) {
	var cs []*pb.Classification
	for _, c := range s.classifications {
		cs = append(cs, c)
	}
	resp := &pb.ListClassificationsResponse{Classifications: cs}
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
	pb.RegisterClassyServer(s, &server{readSubsFile(*subsFile), map[string]*pb.Classification{}})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
