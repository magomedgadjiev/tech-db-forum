package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/models"
	"log"
	"math/rand"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

func ParallelFactory(output interface{}, parallel int, message string, worker func(index int) interface{}) {
	log.Printf("%s - begin", message)
	start := time.Now()
	slice := reflect.ValueOf(output).Elem()
	if slice.Kind() != reflect.Slice {
		panic("Invalid output object type: expected slice")
	}

	results := make(chan []interface{}, parallel)
	var count int32 = int32(slice.Len())
	var index int32 = 0

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			batch := make([]interface{}, count)
			size := 0
			for true {
				idx := atomic.AddInt32(&index, 1) - 1
				if idx >= count {
					break
				}
				batch[size] = worker(int(idx))
				size++
			}
			results <- batch[0:size]
			wg.Done()
		}()
	}
	// wait for the workers to finish
	wg.Wait()

	// get result
	var size int32 = 0
	for i := 0; i < parallel; i++ {
		select {
		case batch := <-results:
			for _, item := range batch {
				slice.Index(int(size)).Set(reflect.ValueOf(item))
				size++
			}
		default:
			break
		}
	}
	close(results)

	if size != count {
		panic(fmt.Sprintf("Unexpected elements count: %d (expected: %d)", size, count))
	}
	elapsed := time.Since(start)
	log.Printf("%s - done: %s", message, elapsed)
}

func Fill(url *url.URL) int {

	start := time.Now()
	// Очистка
	transport := CreateTransport(url)
	c := client.New(transport, nil)
	_, err := c.Operations.Clear(nil)
	CheckNil(err)

	// Создание пользователей
	users := make([]*models.User, 10000)
	ParallelFactory(&users, 8, "Creating users (multiple threads)", func(index int) interface{} {
		return CreateUser(c, nil)
	})
	// Создание форумов
	forums := make([]*models.Forum, 25)
	ParallelFactory(&forums, 8, "Creating forums (multiple threads)", func(index int) interface{} {
		return CreateForum(c, nil, users[rand.Intn(len(users))])
	})
	// Создание веток обсуждения
	threads := make([]*models.Thread, 5000)
	ParallelFactory(&threads, 8, "Creating threads (multiple threads)", func(index int) interface{} {
		thread := RandomThread()
		if rand.Intn(100) >= 5 {
			thread.Slug = ""
		}
		return CreateThread(c, thread, forums[rand.Intn(len(forums))], users[rand.Intn(len(users))])
	})
	// Создание постов
	posts := make([]*models.Post, 1000000)
	ParallelFactory(&posts, 8, "Creating posts (multiple threads)", func(index int) interface{} {
		post := RandomPost()
		post.Author = users[rand.Intn(len(users))].Nickname
		post.Thread = threads[rand.Intn(len(threads))].ID
		return CreatePost(c, post, nil)
	})
	elapsed := time.Since(start)
	log.Printf("Done: %s", elapsed)
	return 0
}
