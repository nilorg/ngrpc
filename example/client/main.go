/*
 Copyright 2019 Nilorg authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

	 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nilorg/ngrpc"
	pb "github.com/nilorg/ngrpc/example/helloworld"
)

func main() {
	client := ngrpc.NewClient("127.0.0.1:5000")
	greeterClient := pb.NewGreeterClient(client.GetConn())

	go func() {
		for {
			time.Sleep(time.Second)
			r, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: "xudeyi"})
			if err != nil {
				log.Printf("could not greet: %v", err)
				continue
			}
			log.Printf("Greeting: %s", r.Message)

		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	defer client.Close()
}
