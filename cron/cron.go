package cron

import (
	"context"
	"encoding/json"

	cronV3 "github.com/robfig/cron/v3"
	"helpers.zhaowenming.cn/logs"
)

//创建秒级定时任务，错误纪录
//
//添加任务 cron.AddXX，返回任务ID并可以取消
//
//添加好后运行 cron.Start()
//
//资料 https://mp.weixin.qq.com/s/Ak7RBv1NuS-VBeDNo8_fww
//
//```
//  添加函数任务
//	cron.AddFunc("[秒0-59] 分0-59 时0-23 天1-31 月1-12 周0-6",func() {
//		......
//	})
//	添加满足Job接口的任务
//	cron.AddJob("各种时间格式",helper.NewContextJob(ctx,func(ctx context.Context) {
//		......
//	}))
//```
func NewCron(skip bool) *cronV3.Cron {
	cl := cronLog{}
	//cron.With* 修改默认行为
	//WithChain指定要应用于添加到此cron的所有作业的作业包装器
	return cronV3.New(
		cronV3.WithChain(
			//cron内置了 3 个用得比较多的JobWrapper：
			//Recover：捕获内部Job产生的 panic；
			//DelayIfStillRunning：触发时，如果上一次任务还未执行完成（耗时太长），则等待上一次任务完成之后再执行；
			//SkipIfStillRunning：触发时，如果上一次任务还未完成，则跳过此次执行。
			cronV3.Recover(&cl), //发生宕机恢复并记录
			func() cronV3.JobWrapper {
				//上一个任务超时后，怎么处理
				//TODO 上一个任务本身，最好做好超时的预案
				if skip {
					//跳过本任务
					return cronV3.SkipIfStillRunning(&cl)
				}
				//继续等待执行
				return cronV3.DelayIfStillRunning(&cl)
			}(),
		),
		cronV3.WithSeconds(),
	)
}

type cronLog struct{}

//普通消息
func (cl *cronLog) Info(msg string, keysAndValues ...interface{}) {
	kvs, _ := json.Marshal(keysAndValues)
	logs.FastNameLog("cron", "info", msg+",kvs:"+string(kvs))
}

//错误消息
func (cl *cronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	kvs, _ := json.Marshal(keysAndValues)
	logs.FastNameLog("cron", "info", "err:"+err.Error()+",msg:"+msg+",kvs:"+string(kvs))
}

//使用上下文来控制任务
type ContextJob struct {
	//上下文控制
	ctx    context.Context
	//回调函数的参数
	params []interface{}
	//回调函数
	fun    func(ctx context.Context, params ...interface{})
}

//创建新的可以通过上下文来控制的任务
func NewContextJob(ctx context.Context, fun func(ctx context.Context, params ...interface{}), params ...interface{}) *ContextJob {
	return &ContextJob{
		ctx:    ctx,
		params: params,
		fun:    fun,
	}
}

func (cj *ContextJob) Run() {
	select {
	case <-cj.ctx.Done():
		//已经完成，结束运行
		return
	default:
		//运行，同时任务本身去通过上下文判断
		cj.fun(cj.ctx, cj.params...)
		return
	}
}
