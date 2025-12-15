package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)
type JobFunc func()
type Job struct{
	schedule *time.Ticker
	excuatable JobFunc
	
	
}
type Cron struct{
	jobs chan Job
	crons int
	wg sync.WaitGroup
	ctx context.Context
	cancel context.CancelFunc
}
func newCronJob(ctx context.Context, cancel context.CancelFunc) *Cron{
	
	return &Cron{
		jobs: make(chan Job),
		crons:0,
		ctx:ctx,
		cancel: cancel,
	}
}
func (c *Cron) Submit(job Job){
	c.jobs<-job
	c.crons++

}
func (c *Cron) worker(job Job){
	defer c.wg.Done()
	defer job.schedule.Stop()
	for {
		select{
		case <-job.schedule.C:
			job.excuatable()
		case <-c.ctx.Done():
			return
		}
	}
}
func (c *Cron) Start(){
	for job:=range c.jobs{
		c.wg.Add(1)
		go c.worker(job)
	}
}
func (c *Cron) Cancel(){
	close(c.jobs)
	c.cancel()
	
}
func execute(){
	fmt.Println("executing task")
}
func execute2(){
	fmt.Println("executing task2")
}

func main(){
	ctx,cancel:=context.WithCancel(context.TODO())
	cr:=newCronJob(ctx,cancel)
	go cr.Start()
	t1:=time.NewTicker(1 *time.Second)
	t2:=time.NewTicker(2 *time.Second)

	job1:=&Job{
		excuatable: execute,
		schedule: t1,
	}
	job2:=&Job{
		excuatable: execute2,
		schedule: t2,
	}
	cr.Submit(*job1)
	cr.Submit(*job2)
	// cr.Cancel()
    cr.wg.Wait()
	cr.Cancel()

}
