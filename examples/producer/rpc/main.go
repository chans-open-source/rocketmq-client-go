/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

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
	"fmt"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

const (
	producerGroup = "please_rename_unique_group_name"
	topic         = "RequestTopic"
	ttl           = 3
)

func main() {
	// create a producer to send reply message
	replyProducer, err := producer.NewDefaultProducer(
		producer.WithGroupName(producerGroup),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
	)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	err = replyProducer.Start()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	msg := primitive.NewMessage(topic, []byte("Hello world"))

	begin := time.Now().UnixNano()
	var retMsg string
	var wg sync.WaitGroup
	wg.Add(1)
	err = replyProducer.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				fmt.Printf("receive message error: %s\n", err)
			} else {
				retMsg = result.String()
			}
			wg.Done()
		},
		msg,
	)
	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	}
	wg.Wait()
	cost := time.Now().UnixNano() - begin
	fmt.Printf("request to <%s> cost: %d replyMessage: %s \n", topic, cost, retMsg)
	err = replyProducer.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}