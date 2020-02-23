// ctask
package ctask

import (
	"log"
	"reflect"
	"time"
)

const (
	Initial   = "Initial"
	Started   = "Started"
	Stopped   = "Stopped"
	Suspend   = "Suspend"
	Warning   = "Warning"
	Exception = "Exception"
)

type Task struct {
	taskId       string
	taskFunc     interface{}
	taskArgs     []string
	taskch       chan string
	taskState    string
	taskInitTime int
	taskWarnTime int
	taskTimer    int
}

type TaskManager struct {
	TaskQueue []*Task
	TaskCount int
	TaskState string
}
type Tasker interface {
	findTaskInQueue(taskId string) int
	AddTask(taskInitTime int, MaxTimeout int, taskId string, taskFunc interface{}, args ...string)
	RemoveTask(taskId string)
	taskStart()
	RunTask()
	taskTimer(ch chan string)
	Dog()
	taskMessage(taskIndex int, message string)
	GetTask()
	// DecTimer()
	// GetTaskStatus()
	// FeedDog(task string)
	// Run()
}

func NewTasker() Tasker {
	tasker := new(TaskManager)
	tasker.TaskCount = 0
	tasker.TaskState = Initial
	tasker.AddTask(5, 10, "taskTimer", tasker.taskTimer)
	return tasker
}

func (t *TaskManager) findTaskInQueue(taskId string) int {
	for i, task := range t.TaskQueue {
		if task.taskId == taskId {
			return i
		}
	}
	return -1
}

func (t *TaskManager) AddTask(taskWarnTime int, taskInitTime int, taskId string, taskFunc interface{}, args ...string) {
	index := t.findTaskInQueue(taskId)
	if index >= 0 {
		log.Printf("AddTask Failed:%s already in TaskQueue", taskId)
	} else {
		ch := make(chan string)
		var argsList []string
		argsList = append(argsList, args...)
		if len(t.TaskQueue) == t.TaskCount {
			task := Task{taskId, taskFunc, argsList, ch, Initial, taskInitTime, taskWarnTime, taskInitTime}
			t.TaskQueue = append(t.TaskQueue, &task)
			t.TaskCount++
		}
		// log.Println("add t:", *t)

	}

}

func (t *TaskManager) RemoveTask(taskId string) {
	index := t.findTaskInQueue(taskId)
	if index < 0 {
		log.Printf("RemoveTask Failed:%s not in TaskQueue", taskId)
	} else {
		close(t.TaskQueue[index].taskch)
		for i := index; i < t.TaskCount-1; i++ {
			t.TaskQueue[i] = t.TaskQueue[i+1]
		}
		t.TaskCount--
		t.TaskQueue[t.TaskCount] = nil
	}
}
func (t *TaskManager) taskMessage(taskIndex int, message string) {
	t.TaskQueue[taskIndex].taskch <- message
}
func (t *TaskManager) taskStart() {
	for _, task := range t.TaskQueue {
		if task.taskState != Started {
			// log.Printf("TaskStart:%s", task.taskId)
			taskFunc := task.taskFunc
			fv := reflect.ValueOf(taskFunc)
			params := make([]reflect.Value, len(task.taskArgs)+1)
			params[0] = reflect.ValueOf(task.taskch)
			for i, arg := range task.taskArgs {
				params[i+1] = reflect.ValueOf(arg)
			}
			go fv.Call(params)
			task.taskState = Started
			t.TaskState = Started
		}
	}
}
func (t *TaskManager) taskRestart(taskId string) {
	index := t.findTaskInQueue(taskId)
	if t.TaskState == Started {
		taskFunc := t.TaskQueue[index].taskFunc
		fv := reflect.ValueOf(taskFunc)
		params := make([]reflect.Value, len(t.TaskQueue[index].taskArgs)+1)
		params[0] = reflect.ValueOf(t.TaskQueue[index].taskch)
		for i, arg := range t.TaskQueue[index].taskArgs {
			params[i+1] = reflect.ValueOf(arg)
		}
		go fv.Call(params)
		t.TaskQueue[index].taskState = Started
	}

}
func (t *TaskManager) RunTask() {
	t.taskStart()
	for {
		cases := make([]reflect.SelectCase, t.TaskCount)
		for i := 0; i < t.TaskCount; i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.TaskQueue[i].taskch)}
		}
		if t.TaskCount > 0 {
			taskIndex, value, _ := reflect.Select(cases)
			msg := value.String()
			if msg == "goodby" {
				t.RemoveTask(t.TaskQueue[taskIndex].taskId)
				// log.Println("goodby", taskIndex)
			} else if msg == "feeddog" {
				t.TaskQueue[taskIndex].taskTimer = t.TaskQueue[taskIndex].taskInitTime
				// log.Println("feeddog", t.TaskQueue[taskIndex].taskId)
				if t.TaskQueue[taskIndex].taskId == "taskTimer" {
					t.Dog()
				}
			}

		} else {
			time.Sleep(time.Second)
			log.Println("TaskQueue is empty")
		}

	}
}

func (t *TaskManager) taskTimer(ch chan string) {
	for {
		time.Sleep(time.Second * 1)
		ch <- "feeddog"
	}
}

func (t *TaskManager) Dog() {
	for i := 0; i < t.TaskCount; i++ {
		t.TaskQueue[i].taskTimer--
		if t.TaskQueue[i].taskTimer <= 0 {
			t.TaskQueue[i].taskState = Exception
			// t.taskMessage(i, "You are Exception")
			log.Println(t.TaskQueue[i].taskId, t.TaskQueue[i].taskState, t.TaskQueue[i].taskTimer)
		} else if t.TaskQueue[i].taskTimer < t.TaskQueue[i].taskWarnTime {
			t.TaskQueue[i].taskState = Warning
			log.Println(t.TaskQueue[i].taskId, t.TaskQueue[i].taskState, t.TaskQueue[i].taskTimer)
		}
	}
}
func (t *TaskManager) GetTask() {
	log.Println(t)
}
func TaskContor(ch chan string) {
	select {
	case e1 := <-ch:
		log.Println(e1)
	default:
	}
}
func TaskFeedDog(ch chan string) {
	ch <- "feeddog"
}
