package main

import (
	"container/list"
	"fmt"
)
type LruCache struct{
	capacity int
	dataList *list.List
	dataMap map[string]*list.Element
}
type Pair struct{
	key string
	val string
}

func newLruCache(size int) *LruCache{
	return &LruCache{
		capacity: size,
		dataList: list.New(),
		dataMap: make(map[string]*list.Element),
	}
}
func (l *LruCache) Set(k,v string){
	if _,ok:=l.dataMap[k];!ok{
		if l.capacity==len(l.dataMap){
			front:=l.dataList.Front()
			pr:=front.Value.(Pair)
			delete(l.dataMap,pr.key)
			l.dataList.Remove(front)
		} 
		ele:=l.dataList.PushBack(Pair{key: k,val: v})
		l.dataMap[k]=ele
	} else{
		prev:=l.dataMap[k]
		l.dataList.Remove(prev)
		ele:=l.dataList.PushBack(Pair{key: k,val: v})
		l.dataMap[k]=ele
	}
}
func (l *LruCache) Get(k string) (bool,string){
	if _,ok:=l.dataMap[k];!ok{
		return false,""
	}
	prev:=l.dataMap[k]
	v:=prev.Value.(Pair).val
	l.dataList.Remove(prev)
	ele:=l.dataList.PushBack(Pair{key: k,val: v})
	l.dataMap[k]=ele
	return true ,v

}


func main(){
	lru:=newLruCache(3)
	lru.Set("a","b")
	lru.Set("c","d")
	lru.Set("a","c")
	lru.Set("d","b")
	lru.Set("f","g")
	fmt.Println(lru.Get("c"))
	fmt.Println(lru.Get("a"))
	lru.Set("l","o")
	fmt.Println(lru.Get("a"))
	fmt.Println(lru.Get("d"))
}