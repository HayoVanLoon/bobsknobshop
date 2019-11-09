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
	pb "github.com/HayoVanLoon/genproto/bobsknobshop/commentcentre/v1"
	"github.com/HayoVanLoon/genproto/bobsknobshop/common/v1"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = 9010
)

func readCommentsFile(fileName string) []common.Comment {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("could not open comments file")
	}
	r := csv.NewReader(f)

	var cs []common.Comment
	for row, err := r.Read(); row != nil; row, err = r.Read() {
		if err != nil {
			log.Fatal("error reading comments file")
		}
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

// Meant for debugging purposes.
func main() {
	var host = flag.String("host", defaultHost, "service host")
	var port = flag.Int("port", defaultPort, "service port")
	var fileName = flag.String("file", "comments.csv", "comment file")
	flag.Parse()

	for _, c := range readCommentsFile(*fileName) {
		createComment(*host, *port, &c)
	}
}
