package mq_service

import (
	"fmt"
	"sync"
)

type Broker interface {
	publish(topic string, msg interface{}) error
	subscribe(topic string) (<-chan interface{}, error)
	unsubscribe(topic string, sub <-chan interface{}) error
	close()
	Broadcast(msg interface{}, subscribers []chan interface{})
	setConditions(capacity int)
}

type Client struct {
	Bro *BrokerImpl
}

func  NewMqClient()*Client{
	return &Client{
		Bro:NewBroker(),
	}
}

func NewBroker()*BrokerImpl{
	return &BrokerImpl{
		make(chan bool),
		0,
		make(map[string][]chan interface{}),
		sync.RWMutex{},
	}
}


func (c *Client)SetConditions(capacity int)  {
	c.Bro.setConditions(capacity)
}
func (c *Client)Publish(topic string, msg interface{}) error{
	return c.Bro.publish(topic,msg)
}
func (c *Client)Subscribe(topic string) (<-chan interface{}, error){
	return c.Bro.subscribe(topic)
}
func (c *Client)Unsubscribe(topic string, sub <-chan interface{}) error {
	return c.Bro.unsubscribe(topic,sub)
}
func (c *Client)Close()  {
	c.Bro.close()
}
func (c *Client)GetPayLoad(sub <-chan interface{})  interface{} {
	for val := range sub {
		if val != nil {
			fmt.Println("res %s",val.(string))
			return val
		}
	}
	return nil
}


