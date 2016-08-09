package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	ES "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/estore"
	"github.com/xyproto/simplebolt"
)

func main() {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	creator := simplebolt.NewCreator(db)
	eventFactory := ES.NewEventFactory()
	eventFactory.RegisterAggregate(NewPost(), NewReply())
	store := estore.NewXyprotoEStore(creator, estore.NewEncoder(eventFactory), estore.NewDecoder(eventFactory))

	//var store = ES.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	ps := NewPostService(store)
	rs := NewReplyService(store)
	eventbus := ES.NewInternalEventBus(store)

	// 注册EventHandler/读模型Handler
	eh := NewEventHandler(ps.CommandChannel())
	readRepository := ES.NewMemoryReadRepository()
	postProjector := NewPostProjector(readRepository)
	eventbus.RegisterHandlers(eh)
	eventbus.RegisterHandlers(postProjector)

	go eventbus.HandleEvents()
	go ps.HandleCommands()
	go rs.HandleCommands()

	// 执行命令
	fmt.Printf("- 创建帖子1\tOK\n")
	post1 := ps.CreatePost("subject1", "body1", "author1")
	fmt.Printf("- 创建帖子2\tOK\n")
	post2 := ps.CreatePost("subject2", "body2", "author2")
	fmt.Printf("- 更新帖子2\tOK\n")
	ps.UpdatePost(post2, "subject2-1", "body2-1")
	fmt.Printf("- 创建帖子1回复\tOK\n")
	reply1 := rs.CreateReply(post1, ES.NewGuid(), "replybody", "author3", time.Now())
	fmt.Printf("- 创建帖子1回复\tOK\n")
	reply2 := rs.CreateReply(post1, ES.NewGuid(), "replybody2", "author4", time.Now())
	fmt.Printf("- 创建帖子2回复\tOK\n")
	reply3 := rs.CreateReply(post2, ES.NewGuid(), "replybody2", "author4", time.Now())

	// 验证
	//wait and print
	go func() {
		time.Sleep(300 * time.Millisecond)
		printEvents(store.GetEvents(ES.NewGuid(), 0, 100))
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", ps.RestoreAggregate(post1))
		fmt.Printf("%v\n------------------\n", ps.RestoreAggregate(post2))
		fmt.Printf("reply1:%v\n------------------\n", rs.RestoreAggregate(reply1))
		fmt.Printf("reply2:%v\n------------------\n", rs.RestoreAggregate(reply2))
		fmt.Printf("reply3:%v\n------------------\n", rs.RestoreAggregate(reply3))

		fmt.Printf("-----------------\nRead Model:\n\n")
		if rPost1, err := readRepository.Find(post1); err == nil {
			fmt.Printf("%v\n------------------\n", rPost1)
		}
		if rPost2, err := readRepository.Find(post2); err == nil {
			fmt.Printf("%v\n------------------\n", rPost2)
		}

		wg.Done()
	}()

	wg.Wait()
}

func printEvents(events []ES.Event) {
	fmt.Printf("-----------------\nEvents after all operations:\n\n")
	for i, e := range events {
		fmt.Printf("%v: %T\n", i, e)
	}
}
