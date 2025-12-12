package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type TaskSchedular struct{
	ctx context.Context
	cancel context.CancelFunc
	jobs chan Task
	rateLimit *time.Ticker
	workerCount int
	wg sync.WaitGroup
	successResult []Result
	FailedResult []Failed


}
type Result struct{
	msg string
}
type Failed struct{
	err string
}
type Job func()(interface{},error)
func successResp()(interface{},error){
	fmt.Println("executed")
	return "executed",nil
}
func failedResp()(interface{},error){
	fmt.Println("failed")
	return "",errors.New("failed")
}


type Task struct{
	id int
	retry int
	job Job
}
func (t *Task) ExecuteTask() (string, error){
	for i:=0;i<t.retry;i++{
		v,err:=t.job()
		if err==nil{
			return v.(string),nil
		}
	}
	return "",errors.New("failed")
}
func (t *TaskSchedular) Worker(ctx context.Context){
	defer t.wg.Done()
	for {
		
		select{
		case <-t.ctx.Done():
			return
		case <-t.rateLimit.C:
			select{
			case job,ok:=<-t.jobs:
			if !ok {
				return
			}
			val,err:=job.ExecuteTask()
			if err !=nil{
				t.FailedResult = append(t.FailedResult, Failed{err:err.Error()} )
			} else{
				t.successResult=append(t.successResult, Result{msg: val})
			}
		}
	}
	}
}

func (t *TaskSchedular) Start(){
	for i:=0;i<t.workerCount;i++{
		t.wg.Add(1)
		go t.Worker(t.ctx)
	}
	// t.wg.Wait()
}
func (t *TaskSchedular) Submit(task Task){
	select{
	case t.jobs<-task:
	return 
	}
}
func (t *TaskSchedular) Close(){
	close(t.jobs)  // important
    t.cancel()
    t.wg.Wait()
	if t.rateLimit != nil {
		t.rateLimit.Stop()
	}
}
func newTaskSchedular(ctx context.Context, rate int64,worker int) *TaskSchedular{
	// limit:=1/rate
	rateLimit := time.NewTicker(time.Second / time.Duration(rate))
	fmt.Println( time.Duration(rate))

	c,cancel:=context.WithCancel(ctx)
	return &TaskSchedular{
		ctx: c,
		cancel: cancel,
		jobs: make(chan Task),
		workerCount: worker,
		successResult: make([]Result,0),
		FailedResult: make([]Failed, 0),
		rateLimit: rateLimit,
	}
}
func main(){
ctx:=context.TODO()
ts:=newTaskSchedular(ctx,1,2)
go ts.Start()
t1:=Task{
	id:1,
	retry:2,
	job: failedResp,
}
ts.Submit(t1)
t2:=Task{
	id:1,
	retry:2,
	job: successResp,
}
ts.Submit(t2)
ts.Close()
ts.wg.Wait()


}