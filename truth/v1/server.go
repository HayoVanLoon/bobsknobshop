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
	"fmt"
	"github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	ccpb "github.com/HayoVanLoon/genproto/bobsknobshop/commentcentre/v1"
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/truth/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort       = 9000
	defaultCcHost     = "commentcentre-service"
	defaultClassyHost = "classy-service"
)

func readCommentsFile(fileName string) [][]string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("could not open comments file")
	}
	r := csv.NewReader(f)

	var rows [][]string
	for row, err := r.Read(); row != nil; row, err = r.Read() {
		if err != nil {
			log.Fatalf("error reading comments file %s", err)
		}
		if len(row) != 4 {
			log.Fatalf("invalid row in comments file %s", row)
		}
		rows = append(rows, row)
	}

	return rows
}

type server struct {
	commentCentreService string
	classyService        string
	exampleComments      [][]string
}

func (s server) GetServiceKpi(ctx context.Context, r *pb.GetServiceKpiRequest) (*pb.GetServiceKpiResponse, error) {
	if r.GetName() == "classy" {
		return s.calculateClassyKpi(ctx, r)
	} else {
		return nil, fmt.Errorf("unsupported service %s", r.GetName())
	}
}

// Calculates the KPI
func (s server) calculateClassyKpi(ctx context.Context, r *pb.GetServiceKpiRequest) (*pb.GetServiceKpiResponse, error) {
	cls, err := s.getClassifications()
	if err != nil {
		return nil, err
	}

	cos, err := s.getComments()
	if err != nil {
		return nil, err
	}
	lbls := getLabels(cos, s.exampleComments)

	ws := map[string]float32{
		// complaints:3, questions:2, support:2, reviews:7 ~= 1 : 1 : 1 : 3
		"complaint": 1, "question": 1, "support": 1, "review": 3,
	}

	startTs, endTs := getBoundaries(r)

	acc := map[string]float32{}
	cnt := map[string]int{}
	fst := map[string]int64{}
	lst := map[string]int64{}
	for _, cl := range cls {
		if startTs <= cl.GetCreatedOn() && cl.GetCreatedOn() < endTs {
			if lbl, ok := lbls[cl.GetComment()]; ok {
				if cl.GetCategory() != lbl {
					acc[cl.GetServiceVersion()] += ws[cl.GetCategory()]
				}
				cnt[cl.GetServiceVersion()] += 1
			}
			if cl.GetCreatedOn() < fst[cl.GetServiceVersion()] {
				fst[cl.GetServiceVersion()] = cl.GetCreatedOn()
			}
			if lst[cl.GetServiceVersion()] < cl.GetCreatedOn() {
				lst[cl.GetServiceVersion()] = cl.GetCreatedOn()
			}
		}
	}

	var vs []*pb.GetServiceKpiResponse_Version
	for version, acc := range acc {
		vs = append(vs, &pb.GetServiceKpiResponse_Version{
			Name:           version,
			StartTimestamp: fst[version],
			EndTimestamp:   lst[version],
			Unit:           "weighted errors",
			Value:          acc / float32(cnt[version]),
		})
	}

	resp := &pb.GetServiceKpiResponse{Versions: vs}
	return resp, nil
}

// Retrieves all classifications.
// Probably not a good idea on production servers with a big load.
func (s server) getClassifications() ([]*classy.Classification, error) {
	conn, err := grpc.Dial(s.classyService, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	cl := classy.NewClassyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := cl.ListClassifications(ctx, &classy.ListClassificationsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetClassifications(), nil
}

// Retrieves all comments.
// Probably not a good idea on production servers with a big load.
func (s server) getComments() (map[string]*common.Comment, error) {
	conn, err := grpc.Dial(s.commentCentreService, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	cl := ccpb.NewCommentcentreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := cl.ListComments(ctx, &ccpb.ListCommentsRequest{})
	if err != nil {
		return nil, err
	}
	xs := map[string]*common.Comment{}
	for _, c := range resp.GetComments() {
		xs[c.GetName()] = c
	}
	return xs, nil
}

// Links Comments to their labels from the comments file
// This is a quick solution for not having actual comment category updates.
func getLabels(cos map[string]*common.Comment, rows [][]string) map[string]string {
	lbls := map[string]string{}
	for _, row := range rows {
		for _, co := range cos {
			if co.GetAuthor() == row[0] && co.GetText() == row[2] {
				lbls[co.GetName()] = row[3]
			}
		}
	}
	return lbls
}

// Returns period, taking into account default zero values.
func getBoundaries(r *pb.GetServiceKpiRequest) (int64, int64) {
	startTs := r.GetStartTimestamp()
	endTs := r.GetEndTimestamp()
	if endTs == 0 {
		endTs = time.Date(2999, 12, 31, 23, 59, 59, 999, time.UTC).Unix()
	}
	return startTs, endTs
}

func main() {
	var port = flag.Int("port", defaultPort, "port to listen on")

	var ccHost = flag.String("cc-host", defaultCcHost, "commentcentre service host")
	var ccPort = flag.Int("cc-port", defaultPort, "commentcentre service port")

	var classyHost = flag.String("classy-host", defaultClassyHost, "classy service host")
	var classyPort = flag.Int("classy-port", defaultPort, "classy service port")

	var fileName = flag.String("file", "comments.csv", "comment file")

	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	ccService := *ccHost + ":" + strconv.Itoa(*ccPort)
	classyService := *classyHost + ":" + strconv.Itoa(*classyPort)
	exampleComments := readCommentsFile(*fileName)
	pb.RegisterTruthServer(s, server{ccService, classyService, exampleComments})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
