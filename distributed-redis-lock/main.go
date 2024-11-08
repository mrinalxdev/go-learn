package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func acquireLock(client *redis.Client, lockKey string, timeout time.Duration) bool {
	ctx := context.Background()

	lockAcquired, err := client.SetNX(ctx, lockKey, "1", timeout).Result()
	if err != nil {
		fmt.Println("Error acquiring lock :", err)
		return false
	}

	return lockAcquired
}

func writeFileWithLock(client *redis.Client, lockKey string, filename string, wg *sync.WaitGroup, data string){
	defer wg.Done()

	if acquireLock(client, lockKey, 10*time.Second){
		fmt.Println("Lock acquired, writing to file ...")

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			releaseLock(client, lockKey)
			return 
		}

		defer file.Close()

		if _, err := file.WriteString(data + "\n"); err != nil {
			fmt.Println("Error writing to file : ", err)
		} else {
			fmt.Println("Successfully wrote to file : ", data)
		}

		releaseLock(client, lockKey)
		fmt.Println("Lock released.")
	} else {
		fmt.Println("Failed to acquire lock. Another process is writing.")
	}
}

func releaseLock(client *redis.Client, lockKey string){
	ctx := context.Background()
	client.Del(ctx, lockKey)
}

func main(){
	// creating the redis client and giving its own localhost port to start
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

	//defining the lock and lock timeout so that only one process should acquire the resource at a time
	lockKey := "file_write_lock"
	filename := "output.txt"

	// wait group to manage goroutines
	var wg sync.WaitGroup
	
	// this is for spawning multiple goroutines for simulating concurrent file writing
	for i := 1; i <= 5; i ++ {
		wg.Add(1)
		data := fmt.Sprintf("Data from goroutine %d", i)
		go writeFileWithLock(client, lockKey, filename, &wg, data)
	}

	wg.Wait()
	fmt.Println("All goroutines finished ..")

}