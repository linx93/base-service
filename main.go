package main

import (
	"fmt"
	"github.com/linx93/base-service/cache"
)

type name1 struct {
	id int
}
type name2 struct {
	Age int
}
type name3 struct {
	Id1 int
}

func main() {

	redisCache := cache.NewRedisCache[name2]("redis-缓存")
	fmt.Println(redisCache)
}
