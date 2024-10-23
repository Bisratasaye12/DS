package main

import (
	"fmt"
	"math/rand"
	"time"
)

// signal and broadcst
func main() {
	gettingReadyForMission()
}

var ready bool

func gettingReadyForMission() {
	go getready()

	checksDone := 0

	for !ready {
		checksDone++
	}

	fmt.Printf("Ready after %v checks.\n", checksDone)
}

func getready() {
	sleep()
	ready = true
}

func sleep() {
	rand.NewSource(time.Now().UnixNano())

	someTime := rand.Intn(5)
	time.Sleep(time.Duration(someTime))
}

// func main() {
// 	var numMemPieces int

// 	memoryPool := &sync.Pool{
// 		New: func() interface{} {
// 			mem := make([]byte, 1024)
// 			numMemPieces++
// 			return &mem
// 		},
// 	}

// 	processes := 1024
// 	var wg sync.WaitGroup
// 	wg.Add(processes)

// 	for i := 0; i < processes; i++ {
// 		go func() {
// 			mem := memoryPool.Get().(*[]byte)
// 			fmt.Println("taking time on using the memory from the pool.")
// 			memoryPool.Put(mem)

// 			wg.Done()
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Printf("memory pool size: %v\n", numMemPieces)
// }

// var missionCompleted bool
// var once sync.Once

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	go func (){
// 		if findTarget(){
// 			once.Do(markMissionCompleted)
// 		}
// 		wg.Done()
// 	}()

// 	wg.Wait()
// 	checkMissionCompleted()
// }

// func markMissionCompleted() {
// 	missionCompleted = true
// }

// func checkMissionCompleted() {
// 	if missionCompleted {
// 		fmt.Println("mission completed.")
// 	} else {
// 		fmt.Println("mission failure.")
// 	}
// }

// func findTarget() bool {
// 	rand.NewSource(time.Now().UnixNano())
// 	target := 0

// 	return rand.Intn(10) == target
// }

// var (
// 	mu sync.Mutex
// 	count = 0
// )

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1000)
// 	for i := 0; i < 1000; i++ {
// 		go increment(&wg)
// 	}
// 	wg.Wait()
// 	fmt.Println(count)
// }

// func increment(wg *sync.WaitGroup) {
// 	mu.Lock()
// 	count++
// 	mu.Unlock()
// 	wg.Done()
// }

// func main(){
// 	var wg sync.WaitGroup

// 	evils := []string{"tommy", "johnny", "mani"}

// 	wg.Add(len(evils))

// 	for i := 0; i < len(evils); i++{
// 		go attack(evils[i], &wg)
// 	}

// 	wg.Wait()
// 	fmt.Println("mission accomplished!")
// }

// func attack(name string, wg *sync.WaitGroup){
// 	fmt.Printf("attacked %v\n", name)
// 	wg.Done()
// }

// func main() {
// 	ninja1 := make(chan string)
// 	ninja2 := make(chan string)

// 	go elect("ninja1", ninja1)
// 	go elect("ninja2", ninja2)

// 	select {
// 	case message := <-ninja1:
// 		fmt.Println(message)
// 	case message := <-ninja2:
// 		fmt.Println(message)
// 	}
// }

// func elect(ninja string, tracker chan string) {
// 	tracker <- fmt.Sprintf("%v elected", ninja)
// }

// func main(){
// 	start := time.Now()
// 	defer func (){
// 		fmt.Println(time.Since(start))
// 	}()

// 	// evilNinja := "ninja"
// 	signal_channel := make(chan bool, 1)
// 	signal_channel <- true
// 	// go attack(evilNinja, signal_channel)

// 	fmt.Printf("the value from channel: %v\n", <-signal_channel)
// }

// func attack(ninja string, sig_chan chan bool){
// 	time.Sleep(time.Second)
// 	fmt.Printf("..attacking ninja %v\n, the time: %v\n", ninja, time.Now())
// 	sig_chan <- true
// }

// func main(){
// 	channel := make(chan string)
// 	go score(channel)

// 	for message := range channel {
// 		fmt.Println(message)
// 	}

// }

// func score(channel chan string){
// 	rand.NewSource(time.Now().UnixNano())
// 	rounds := 10
// 	defer close(channel)

// 	for i := 0; i < rounds; i++{
// 		scr := rand.Intn(10)
// 		channel <- fmt.Sprintf("..you scored: %v", scr)
// 	}
// }
