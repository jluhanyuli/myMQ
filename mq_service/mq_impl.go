package mq_service

import (
	"errors"
	"sync"
	"time"
)

type BrokerImpl struct {
	exit chan bool
	capacity int
	topics map[string][]chan interface{} // key： topic  value ： queue
	sync.RWMutex // 同步锁
}

func (b *BrokerImpl)publish(topic string,msg interface{})error{
	select {
		case <-b.exit:
			return errors.New("broker closed")
	default:
	}
	b.RLock()
	subscribers,ok:=b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}
	b.broadcast(msg, subscribers)//广播
	return nil
}

func (b *BrokerImpl)broadcast(msg interface{},subscribers [] chan interface{}){
	count :=len(subscribers)
	concurrency:=1
	switch  {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}
	pub:=func(start int){
		for i:=start;i<count;i+=concurrency{
			select {
			case subscribers[i]<- msg:
			case <-time.After(time.Millisecond * 5):   //等待参数duration时间后，向返回的chan里面写入当前时间
			case <-b.exit:
				return
			default:
				return
			}
		}
	}
	for i:=0;i<concurrency;i++{
		go pub(i)
	}
}

func (b *BrokerImpl)subscribe(topic string)(<-chan interface{},error){
	select {
	case <-b.exit:
		return nil,errors.New("broker closed")
	default:
	}
	ch :=make(chan interface{},b.capacity)
	b.Lock()
	b.topics[topic]=append(b.topics[topic],ch)
	b.Unlock()
	return ch,nil
}

func (b *BrokerImpl) unsubscribe(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}
	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}
	// delete subscriber
	var newSubs []chan interface{}
	for _, subscriber := range subscribers {
		if subscriber == sub {
			continue
		}
		newSubs = append(newSubs, subscriber)
	}
	b.Lock()
	b.topics[topic] = newSubs
	b.Unlock()
	return nil
}

func (b *BrokerImpl) close()  {
	select {
	case <-b.exit:
		return
	default:
		close(b.exit)
		b.Lock()
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}
	return
}

func (b *BrokerImpl)setConditions(capacity int)  {
	b.capacity = capacity
}

