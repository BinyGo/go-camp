package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx         = context.Background()
	ip          = "localhost"
	port uint16 = 6379
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%v", ip, port),
		WriteTimeout: time.Second * 10,
	})
	analysis(client, "10000个10字节", 10000, 10)
	analysis(client, "500000个10字节", 500000, 10)

	analysis(client, "10000个200字节", 10000, 20)
	analysis(client, "500000个20字节", 500000, 20)

	analysis(client, "10000个1k字节", 10000, 1024)
	analysis(client, "500000个1k字节", 500000, 1024)

	analysis(client, "10000个5k字节", 10000, 5120)
	analysis(client, "500000个5k字节", 500000, 5120)

	// fmt.Println("-------------------Printf_Redis_Memory_Info--------------------")
	// Printf_Redis_Memory_Info(client)
}

//insert 通过pipeline向redis写入指定count的kv数据
func insert(redisClient *redis.Client, key string, value string, count int) {
	fmt.Printf("--------------------插入:%s--------------------\n", key)
	batchCount := 1 * 10000

	pipe := redisClient.Pipeline()
	for i := 0; i < count; i++ {
		newKey := fmt.Sprintf("%s:%v", key, i)
		pipe.Set(ctx, newKey, value, -1)
		if i%batchCount == 0 {
			_, err := pipe.Exec(ctx)
			if err != nil {
				panic(err)
			}
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
}

//generateValue 生成指定字节的value
func generateValue(dataSize int) string {
	bytes := make([]byte, dataSize)
	for i, _ := range bytes {
		bytes[i] = 'a'
	}
	return string(bytes)
}

//analysis 分析redis的内存,并将内存信息输出到控制台
func analysis(redisClient *redis.Client, key string, count int, size int) {
	redisClient.FlushDB(ctx)

	val, _ := redisClient.Info(ctx, "memory").Result()

	Regex := regexp.MustCompile(`used_memory:([0-9]+)\r\n`)
	Regex2 := regexp.MustCompile(`[0-9]+`)

	beforeStr := fmt.Sprintf("%s", Regex2.Find(Regex.Find([]byte(val))))
	beforeMemory, _ := strconv.ParseInt(beforeStr, 10, 64)

	value := generateValue(size)
	insert(redisClient, key, value, count)

	val2, err := redisClient.Info(ctx, "memory").Result()
	afterStr := fmt.Sprintf("%s", Regex2.Find(Regex.Find([]byte(val2))))
	afterMemory, _ := strconv.ParseInt(afterStr, 10, 64)

	averageSize := (afterMemory - beforeMemory) / int64(count)
	fmt.Printf("插入前used_memory:%d(byte)\n", beforeMemory)
	fmt.Printf("插入后used_memory:%d(byte)\n", afterMemory)
	fmt.Printf("平均每个 key 的占用内存空间:%d(byte)\n", averageSize)

	if err != nil {
		panic(err)
	}
}

func Printf_Redis_Memory_Info(redisClient *redis.Client) {
	redisClient.FlushDB(ctx)
	val, _ := redisClient.Info(ctx, "memory").Result()
	fmt.Println(val)
}
