package main

import (
	"fmt"
	"log"
	"myMq/mq_service"
	"sync"
	"time"
)

func TestClient(){
	b := mq_service.NewMqClient()
	b.SetConditions(100)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		topic := fmt.Sprintf("Golang%d", i)
		payload := fmt.Sprintf("lyh%d", i)
		ch, err := b.Subscribe(topic)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go func() {
			e := b.GetPayLoad(ch)
			if e != payload {
				log.Fatal("%s expected %s but get %s", topic, payload, e)
			}
			if err := b.Unsubscribe(topic, ch); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
		if err := b.Publish(topic, payload); err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
}
//多个topic测试
func ManyTopic()  {
	m := mq_service.NewMqClient()
	defer m.Close()
	m.SetConditions(10)
	top := ""
	for i:=0;i<10;i++{
		top = fmt.Sprintf("Golang梦工厂_%02d",i)
		go Sub(m,top)
	}
	ManyPub(m)
}

func ManyPub(c *mq_service.Client)  {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for  {
		select {
		case <- t.C:
			for i:= 0;i<10;i++{
				//多个topic 推送不同的消息
				top := fmt.Sprintf("Golang_%02d",i)
				payload := fmt.Sprintf("lyh_%02d",i)
				err := c.Publish(top,payload)
				if err != nil{
					fmt.Println("pub message failed")
				}
			}
		default:
		}
	}
}
func Sub(c *mq_service.Client,top string)  {
	ch,err := c.Subscribe(top)
	if err != nil{
		fmt.Printf("sub top:%s failedn",top)
	}
	for  {
		val := c.GetPayLoad(ch)
		if val != nil{
			fmt.Printf("%s get message is %sn",top,val)
		}
	}
}
