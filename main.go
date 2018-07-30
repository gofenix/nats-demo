package main

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
	"log"
	"strconv"
	"time"
)

func main() {
	var clusterId string = "test-cluster"
	var clientId string = "test-client"

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
		return
	}

	// 开启一个协程，不停的生产数据
	go func() {
		m := 0
		for {
			m++
			sc.Publish("foo1", []byte("hello message "+strconv.Itoa(m)))
			time.Sleep(time.Second)
		}

	}()

	// 消费数据
	i := 0
	mcb := func(msg *stan.Msg) {
		i++
		log.Println(i, "---->", msg.Subject, msg)
	}
	startOpt := stan.StartAt(pb.StartPosition_LastReceived)
	//_, err = sc.QueueSubscribe("foo1", "", mcb, startOpt)   // 也可以用queue subscribe
	_, err = sc.Subscribe("foo1", mcb, startOpt)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	// 创建一个channel，阻塞着
	signalChan := make(chan int)
	<-signalChan
}
