package main

import (
	"flag"
	"github.com/helloshiki/ndb"
	"github.com/helloshiki/ndb/example/controller"
	"github.com/helloshiki/ndb/example/models"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	CpuProfile  = flag.String("cpu-profile", "", "write cpu profile to file")
	HeapProfile = flag.String("heap-profile", "", "write heap profile to file")
)

func main() {
	log.Printf("main")
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if *CpuProfile != "" {
		file, err := os.Create(*CpuProfile)
		if err != nil {
			log.Panicln(err)
		}
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	if *HeapProfile != "" {
		file, err := os.Create(*HeapProfile)
		if err != nil {
			log.Panicln(err)
		}
		defer pprof.WriteHeapProfile(file)
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//for test
	sample()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	log.Println(<-sigChan)

}

func sample() {
	go Work()

	//////////////////建表/////////////////////////////////////////
	u1 := models.User{UID: 1, GID: 0, TCC: decimal.New(99, 2), ETH: decimal.New(199, 2), NASH: decimal.New(299, 2), Worker: map[int]bool{2: true}, I1: 1000}
	ndb.CreateTable(&u1)
	ndb.Insert(&u1)
	ndb.Insert(&u1)

	u2 := models.User{UID: 2, GID: 0, TCC: decimal.New(99, 2), ETH: decimal.New(199, 2), NASH: decimal.New(299, 2), Worker: map[int]bool{1: true}}
	ndb.Insert(&u2)

	m1 := models.TchMachine{}
	ndb.CreateTable(&m1)
	//////////////////建表/////////////////////////////////////////

	///////////////////插入/////////////////////////////////////////
	ucnt := 10
	//插入ucnt个用户
	for i := 0; i <= ucnt; i++ {
		u := models.User{UID: i, GID: 0, TCC: decimal.New(99, 2), ETH: decimal.New(199, 2), NASH: decimal.New(299, 2), Worker: map[int]bool{1: true}}
		ndb.Insert(&u)
	}
	mcnt := 1000000
	start := time.Now().Unix()

	//插入mcnt台矿机
	for i := 0; i < mcnt; i++ {
		m := models.TchMachine{ID: i, GID: 0, UID: i % ucnt}
		//log.Printf("m:+%v", m)
		ndb.Load(&m)
	}
	end := time.Now().Unix()
	log.Printf("insert %d records in %d second", mcnt, end-start)
	///////////////////插入/////////////////////////////////////////

	///////////////////删除/////////////////////////////////////////
	u10 := models.User{UID: 10}
	ndb.Delete(&u10)
	///////////////////删除/////////////////////////////////////////

	///////////////////更新/////////////////////////////////////////
	start = time.Now().Unix()
	for i := 0; i < mcnt; i++ {
		u := models.User{UID: i % 10}
		ndb.UpdateField(&u, "TCC", "INC", decimal.New(123, -1), false)
	}

	end = time.Now().Unix()
	log.Printf("update %d records in %d second", mcnt, end-start)

	start = time.Now().Unix()
	for i := 0; i < mcnt; i++ {
		m := models.TchMachine{ID: i % 10, GID: 0, UID: i % ucnt}
		ndb.Update(&m)
	}
	end = time.Now().Unix()
	log.Printf("update %d records in %d second", mcnt, end-start)
	///////////////////更新/////////////////////////////////////////

	///////////////////转账/////////////////////////////////////////////
	//ndb.UpdateFunc((controller.Transfer(nil, nil, nil)).(ndb.CallBack))
	log.Printf("before transfer: user1: %+v, user2: %+v", ndb.Get(&u1), ndb.Get(&u2))

	controller.Transfer(1, 2, "TCC", decimal.New(1, 1))
	controller.Transfer(1, 2, "TCC", decimal.New(1, 100))

	controller.Transfer(1, 2, "ETH", decimal.New(1, -1))

	controller.Transfer(1, 2, "NASH", decimal.New(1, 2))

	log.Printf("after transfer: user1: %+v, user2: %+v", ndb.Get(&u1), ndb.Get(&u2))

	log.Printf("before u1 is %+v", u1)
	if b, e, err := ndb.UpdateField(&u1, "Desc", "REPLACE", "ssss", false); err == nil {
		log.Printf("DESC change from %s to %s", b, e)
	} else {
		log.Printf("err is %+v", err)
	}
	if b, e, err := ndb.UpdateField(&u1, "Worker", "REPLACE", map[int]bool{3: true}, false); err == nil {
		log.Printf("Worker change from %s to %s", b, e)
	} else {
		log.Printf("err is %+v", err)
	}

	if b, e, err := ndb.UpdateField(&u1, "I1", "ZERO", 0, true); err == nil {
		log.Printf("I1 change from %s to %s", b, e)
	} else {
		log.Printf("err is %+v", err)
	}

	if b, e, err := ndb.UpdateField(&u1, "I1", "REPLACE", 1000, true); err == nil {
		log.Printf("I1 change from %s to %s", b, e)
	} else {
		log.Printf("err is %+v", err)
	}
	if b, e, err := ndb.UpdateField(&u1, "I1", "INC", 100, true); err == nil {
		log.Printf("I1 change from %s to %s", b, e)
	} else {
		log.Printf("err is %+v", err)
	}
	if b, e, err := ndb.UpdateField(&u1, "I1", "DEC", 10, false); err == nil {
		log.Printf("I1 change from %s to %s", b, e)
	} else {
		log.Printf("err is %+v", err)
	}

	log.Printf("after u1 is %+v", u1)

	m := models.TchMachine{GID: 0, UID: 1}
	arr := ndb.GetByIndex(&m, "guid")
	log.Printf("GID:0 UID:1's TchMachine count is %d", len(arr))

	arr = ndb.GetByIndex(&m, "pk")
	log.Printf("TchMachine count is %d", len(arr))

	time.Sleep(5 * time.Second)
	stat := ndb.GetStat(&m, "guid", true)
	for i := 0; i < len(stat); i++ {
		log.Printf("TchMachine stat[%d] is %+v", i, stat[i])

	}
	stat = ndb.GetStat(&m, "guid", false)
	for i := 0; i < len(stat); i++ {
		log.Printf("TchMachine stat[%d] is %+v", i, stat[i])

	}

	///////////////////转账/////////////////////////////////////////////

}
