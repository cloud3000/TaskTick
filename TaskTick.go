package main

// LogTime, is the heart-beat of scheduling tasks.

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func run(tick int, t que) {
	// fmt.Printf("%v: %s Tick: %2d \n", time.Now().Format("2006-01-02 15:04:05"), t.Title, tick)
	for _, task := range t.Tasks {
		if task.Tick == tick {
			go launch(task.Script)
		}
	}
}

//  watchFile contains the clocks for each circular processQ, one for Hours, Minutes, and Seconds.
var configPath = "config"
var configFile = "config/clock.json"

func main() {
	// readConf, channels, and watcher are all in support of realtime changes to processQ's.
	processQ := readConf(configFile)
	newConf := make(chan bool)
	quit := make(chan bool)
	wg := new(sync.WaitGroup) // I want to wait for watcher to end.
	wg.Add(1)                 // I know I don't need to wait, I'm doing it just because I should.

	// After reading the config file, configuration will emmit config file changes.
	go configuration(wg, configPath, configFile, quit, newConf)

	// The ticker and timing variables are for logging time.
	clockTick := time.NewTicker(time.Second)
	min := time.Minute
	hr := time.Hour
	deadline := hr * 24
	//msg := " "

	// This process will run until 'deadline' seconds has passed
	secQ := 0
	minQ := 0
	hrQ := 0
	for sec := time.Second; sec <= deadline; sec = sec + time.Second {

		select {
		case <-clockTick.C:
			/* Every second run the Secque Tasks.
			   Secque is a 10 second circular que.
			   You can run a task every 10 seconds by select which 1 of 10 seconds to run on.
			   If you want it to run every 5 seconds, put the task in twice, on second 0 and 5
			*/
			secQ++
			if secQ > 9 {
				secQ = 0
			}
			go run(secQ, processQ.Secque)
			//msg = fmt.Sprintf("Second Que is :%2d", secQ)

			if sec%min == 0 {
				/* Every minute run the Minque Tasks
					   Minque is a 60 minute circular que.
				   	   You can run a task every 60 minutes by select which 1 of 60 minutes to run on.
				       If you want it to run every 30 minutes, put the task in twice, maybe on minute 15 and 45
				*/
				minQ++
				if minQ > 59 {
					minQ = 0
				}
				//msg = fmt.Sprintf("Minute Que is :%2d", minQ)
				go run(minQ, processQ.Minque)
			}
			if sec%hr == 0 {
				/* Every hour run the Hrque Tasks
					   Hrque is a 24 hour circular que.
				   	   You can run a task every 24 hours by select which 1 of 24 hours to run on.
				       If you want it to run every 12 hours, put the task in twice, maybe on hour 3 and 15
				*/
				hrQ++
				if hrQ > 24 {
					hrQ = 0
				}
				//msg = fmt.Sprintf("Hour Que is :%2d", hrQ)
				go run(hrQ, processQ.Hrque)
			}
			// Display msg.
			// log.Printf("%s, total time: %v \n", msg, sec)

		case <-newConf: // Fires whenever clock.json is changed
			processQ = readConf(configFile)
			fmt.Println("Configuration Re-Loaded")
		}

		// If passed deadline, tell the configuration to quit, log the event, and stop the ticker.
		if sec >= deadline {
			quit <- true
			log.Printf("Deadline of %v has passed.\n", sec)
			clockTick.Stop()
			break
		}

	}
	// Wait for the configuration to stop, it's the right thing to do ;-)
	wg.Wait()
	println("Main is shutting down.")
}
