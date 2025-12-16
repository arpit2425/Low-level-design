package main

import (
	"fmt"
	"time"
)

type Storage struct{
	collection map[string]*Value
	deletedKeys chan string
}

type Value struct{
	val string
	ttl time.Time
}

func newCache() *Storage{
	return &Storage{
		collection: make(map[string]*Value),
		deletedKeys: make(chan string),
	}
}
func (s *Storage) Cleanup(){
	for {
		select{
		case val,ok:=<-s.deletedKeys:
			if !ok{
				return
			}
			fmt.Println("cleannup done")
		delete(s.collection,val)
		}
	}
}
func (s *Storage) Set(k,v string) (bool,error){
	if _,ok:=s.collection[k];ok{
		value:=&Value{val: v}
		s.collection[k]=value
		return false,nil
	}
	value:=&Value{val: v}
		s.collection[k]=value
		return true,nil

}
func (s *Storage) SetTtl(k,v string, ttl int) (bool,error){
	
		now:=time.Now()
		t:=now.Add(time.Duration(ttl) * time.Second)
		value:=&Value{val: v,ttl: t}
		s.collection[k]=value
		return false,nil
}
func (s *Storage) Get(k string) (string,bool){
	if _,ok:=s.collection[k];!ok{
		return "",false
	}
	val:=s.collection[k]
	
	now:=time.Now()
	if  !val.ttl.IsZero() && now.After(val.ttl) {
		s.deletedKeys<-k
		return "",false
	}
	return val.val ,true
}
func main(){
	cache:=newCache()
	
	go cache.Cleanup()
	cache.Set("a","b")
	cache.Set("b","c")
	cache.Set("a","d")
	val,_:=cache.Get("a")
	fmt.Println(val)
	val,_=cache.Get("b")
	fmt.Println(val)
	cache.SetTtl("p","q",2)
	val,_=cache.Get("p")
	fmt.Println(val)
	time.Sleep(3 *time.Second)
	val,_=cache.Get("p")
	fmt.Println(val)
	// time.Sleep(5*time.Second)
	// close(cache.deletedKeys)
}