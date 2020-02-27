# ctask
An esaey go task mannager

# example
package main

import (
	"fmt"
	"time"

	"github.com/wentome/ctask"
)

func main() {
	tasker := ctask.NewTasker()
	tasker.AddTask(5, 10, "Task1", Task1, "args")
	for i := 0; i < 1; i++ {
		taskId := fmt.Sprintf("Task2-%d", i)
		tasker.AddTask(5, 10, taskId, Task2)
	}
	tasker.RunTask()
}

func Task1(ch chan string, arg1 string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("异常信息为:", err)
		}
	}()
	for {
		ctask.TaskContor(ch)
		ctask.TaskFeedDog(ch)
		time.Sleep(time.Second * 3)
	}
}

func Task2(ch chan string) {
	for {
		ctask.TaskContor(ch)
		ctask.TaskFeedDog(ch)
		time.Sleep(time.Second * 3)
	}

}

