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
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/HayoVanLoon/genproto/bobsknobshop/classy/v1"
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/commentcentre/v1"
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultHost       = "localhost"
	defaultPort       = 9010
	defaultClassyHost = "localhost"
	defaultClassyPort = 9000
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
			log.Fatal("error reading comments file")
		}
		if len(row) != 4 {
			log.Fatalf("invalid row in comments file %s", row)
		}
		rows = append(rows, row)
	}

	return rows
}

func toComments(rows [][]string) []common.Comment {
	var cs []common.Comment

	for _, row := range rows {
		if len(row) != 4 {
			log.Fatalf("invalid row in comments file %s", row)
		}
		cs = append(cs, common.Comment{
			Topic:  row[1],
			Author: row[0],
			Text:   row[2],
		})
	}

	return cs
}

func createComment(host string, port int, c *common.Comment) {
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer func() { _ = conn.Close() }()

	cl := pb.NewCommentcentreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := cl.CreateComment(ctx, c)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%s\n", resp)
}

func calculateKpi(host, classyHost string, port, classyPort int, rows [][]string) map[string]float32 {
	cls := getClassifications(classyHost, classyPort)
	cos := getComments(host, port)

	lbls := map[string]string{}
	for _, row := range rows {
		for _, co := range cos {
			if co.GetAuthor() == row[0] && co.GetText() == row[2] {
				lbls[co.GetName()] = row[3]
			}
		}
	}

	kpisAcc := map[string]int{}
	kpisCnt := map[string]int{}
	for _, c := range cls {
		if co, ok := cos[c.GetComment()]; ok {
			if lbl, ok := lbls[co.GetName()]; ok {
				if c.GetCategory() == lbl {
					kpisAcc[c.GetServiceVersion()] += 1
				}
				kpisCnt[c.GetServiceVersion()] += 1
			}
		}
	}

	kpis := map[string]float32{}
	for version, correct := range kpisAcc {
		kpis[version] = float32(correct) / float32(kpisCnt[version])
	}

	return kpis
}

func getClassifications(host string, port int) []*classy.Classification {
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer func() { _ = conn.Close() }()
	cl := classy.NewClassyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := cl.ListClassifications(ctx, &classy.ListClassificationsRequest{})
	if err != nil {
		panic(err)
	}
	return resp.GetClassifications()
}

func getComments(host string, port int) map[string]*common.Comment {
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer func() { _ = conn.Close() }()
	cl := pb.NewCommentcentreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := cl.ListComments(ctx, &pb.ListCommentsRequest{})
	if err != nil {
		panic(err)
	}
	xs := map[string]*common.Comment{}
	for _, c := range resp.GetComments() {
		xs[c.GetName()] = c
	}
	return xs
}

// Meant for debugging purposes.
func main() {
	var host = flag.String("host", defaultHost, "service host")
	var port = flag.Int("port", defaultPort, "service port")

	var fileName = flag.String("file", "comments.csv", "comment file")

	var classyHost = flag.String("classy-host", defaultClassyHost, "classy service host")
	var classyPort = flag.Int("classy-port", defaultClassyPort, "classy service port")

	flag.Parse()

	rows := readCommentsFile(*fileName)

	for _, c := range toComments(rows) {
		createComment(*host, *port, &c)
	}

	kpis := calculateKpi(*host, *classyHost, *port, *classyPort, rows)

	fmt.Println(kpis)
}
