package main

import (
	"fmt"
	"sync"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
type UrlShortner struct{
	urlMap map[string]string
	reverseMap map[string]string
	currentId int64
	mu sync.Mutex
}
func newUrlShortner() *UrlShortner{
	return &UrlShortner{
		urlMap: make(map[string]string),
		reverseMap: make(map[string]string),
		currentId: 1000,
	}
}
func generateKey(i int64) string{
	s:=""
	for i>0 {
		k:=i%62
		s=string(charset[k]) + s
		i=i/62
	}
	return s
}
func (u *UrlShortner) Encode(url string) string{
	u.mu.Lock()
	defer u.mu.Unlock()
	key:=generateKey(u.currentId)
	fmt.Println("key",key)
	u.currentId=u.currentId+1
	s:=fmt.Sprintf("bit.ly/%s",key)
	u.urlMap[s]=url
	u.reverseMap[url]=s
	return s
}
func (u *UrlShortner) Decode(shortUrl string) string{
	u.mu.Lock()
	defer u.mu.Unlock()
	if val,ok:=u.urlMap[shortUrl]; ok{

		return val
	}
	return ""

}
func main(){
	urlS:=newUrlShortner()
	fmt.Println("data")
	s:=urlS.Encode("www.google.com/jhjedhfberhbrecbcebgjcbeghjbcerdghjbcredj")
	fmt.Println("shorturl",s)
	os:=urlS.Decode(s)
	fmt.Println("longurl",os)
}