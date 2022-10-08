package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Job1 任务
type Job1 struct {
	JobID     int
	currPri   int  // 当前队列
	ticksLeft int  // 当前任务剩余时间片长度
	allotLeft int  // 剩余时间片次数
	startTime int  // 开始的时间周期
	endTime   int  // 结束时间
	runTime   int  // 运行时间
	timeLeft  int  // 剩余运行时间
	ioFreq    int  // io操作频率
	doingIO   bool // 是否在io
	firstRun  int  // 第一次run的时间
}

// jobStatus 任务状态
type jobStatus struct {
	jobIdx int
	status string
}

var seed1 = flag.Int64("s", 0, "指定随机种子")
var numQueues = flag.Int("n", 3, "MLFQ中的队列数（如果没有使用-Q）")
var quantum1 = flag.Int("q", 10, "时间片长度(如果没有使用-Q参数)")
var allotment = flag.Int("a", 1, "分配长度（如果不使用-A）")
var quantumList = flag.String("Q", "", "指定为x，y，z，...为每个队列级别的时间片长度，其中x是优先级最高的队列的时间片长度，y是第二高的队列的时间片长度，依此类推")
var allotmentList = flag.String("A", "", "x、 y，z，。。。其中x是最高优先级队列，y次高，依此类推")
var numJobs = flag.Int("j", 3, "系统中的作业数")
var maxlen1 = flag.Int("m", 100, "作业的最大运行时间（如果是随机的）")
var maxio = flag.Int("M", 10, "作业的最大I/O频率（如果是随机的）")
var boost = flag.Int("B", 0, "将所有作业的优先级提高到高优先级的频率（0表示从不）")
var ioTime = flag.Int("i", 5, "I/O 持续时间（固定常数)")
var stay = flag.Bool("S", false, "发出I/O时重置并保持相同的优先级")
var iobump = flag.Bool("I", false, "如果指定，完成I/O的作业将立即移动到当前队列前面")
var jlist1 = flag.String("l", "", "以逗号分隔的要运行的作业列表，格式为x1，y1，z1：x2，y2，z2：...。其中x是开始时间，y是运行时间，z是作业I/O的频率")
var solve1 = flag.Bool("c", false, "计算答案")

var queues [][]*Job1
var hiQueue int

func main() {
	flag.Parse()

	rand.New(rand.NewSource(*seed1))

	// 队列时间片长度
	var quantums []int
	if *quantumList != "" {
		quantumStr := strings.Split(*quantumList, ",")
		*numQueues = len(quantumStr)
		quantums = make([]int, *numQueues)
		for i, item := range quantumStr {
			quantums[*numQueues-1-i], _ = strconv.Atoi(item)
		}
	} else {
		quantums = make([]int, *numQueues)
		for i := 0; i < *numQueues; i++ {
			quantums[i] = *quantum1
		}
	}

	// 配给时间片长度
	allotments := make([]int, *numQueues)
	if *allotmentList != "" {
		allotmentStr := strings.Split(*allotmentList, ",")
		if *numQueues != len(allotmentStr) {
			fmt.Println("指定的分配数必须与队列数一致")
			os.Exit(1)
		}
		for i, item := range allotmentStr {
			allotments[*numQueues-1-i], _ = strconv.Atoi(item)
		}
	} else {
		for i := 0; i < *numQueues; i++ {
			allotments[i] = *allotment
		}
	}

	// 最高队列
	hiQueue = *numQueues - 1
	// jlist1 'startTime,runTime,ioFreq:startTime,runTime,ioFreq:...'
	var jobs []*Job1
	if *jlist1 != "" {
		jobStr := strings.Split(*jlist1, ":")
		*numJobs = len(jobStr)
		jobs = make([]*Job1, *numJobs)
		for i, job := range jobStr {
			jobInfo := strings.Split(job, ",")
			if len(jobInfo) != 3 {
				fmt.Println("作业字符串格式错误。应为x1、y1、z1:x2、y2、z2：。。。")
				fmt.Println("其中x是开始时间，y是运行时，z是输入/输出频率。")
				os.Exit(1)
			}
			startTime, _ := strconv.Atoi(jobInfo[0])
			runTime, _ := strconv.Atoi(jobInfo[1])
			ioFreq, _ := strconv.Atoi(jobInfo[2])
			jobs[i] = &Job1{
				JobID:     i,
				currPri:   hiQueue,
				ticksLeft: quantums[hiQueue],
				allotLeft: allotments[hiQueue],
				startTime: startTime,
				runTime:   runTime,
				timeLeft:  runTime,
				ioFreq:    ioFreq,
				doingIO:   false,
				firstRun:  -1,
			}
		}
	} else {
		jobs = make([]*Job1, *numJobs)
		for i := 0; i < *numJobs; i++ {
			runTime := rand.Intn(*maxlen1) + 1
			ioFreq := 0
			if *maxio != 0 {
				ioFreq = rand.Intn(*maxio)
			}
			jobs[i] = &Job1{
				JobID:     i,
				currPri:   hiQueue,
				ticksLeft: quantums[hiQueue],
				allotLeft: allotments[hiQueue],
				startTime: 0,
				runTime:   runTime,
				timeLeft:  runTime,
				ioFreq:    ioFreq,
				doingIO:   false,
				firstRun:  -1,
			}
		}
	}

	// ioDone
	ioDone := make(map[int][]*jobStatus)
	for i, job := range jobs {
		if _, ok := ioDone[job.startTime]; !ok {
			ioDone[job.startTime] = make([]*jobStatus, 0)
		}
		ioDone[job.startTime] = append(ioDone[job.startTime], &jobStatus{
			jobIdx: i,
			status: "JOB BEGINS",
		})
	}

	// 输出信息
	fmt.Println("以下是输入列表：")
	fmt.Println("选项 作业数：", *numJobs)
	fmt.Println("选项 队列数：", *numQueues)
	for i := *numQueues - 1; i >= 0; i-- {
		fmt.Printf("选项 队列 %3d 的时间片使用次数是： %3d\n", i, allotments[i])
		fmt.Printf("选项 队列 %3d 的时间片长度是： %3d\n", i, quantums[i])
	}
	fmt.Println("选项 提高优先级频率：", *boost)
	fmt.Println("选项 I/O持续时间：", *ioTime)
	fmt.Println("选项 发出I/O时重置并保持相同的优先级：", *stay)
	fmt.Println("选项 完成I/O的作业将立即移动到当前队列前面：", *iobump)
	fmt.Println()
	fmt.Println("对于每个作业，给出了三个定义特征：")
	fmt.Println("\t开始时间：作业何时进入系统")
	fmt.Println("\t运行时间：作业完成所需的总CPU时间")
	fmt.Println("\tio频率：每个ioFreq时间单位，作业都会发出I/O")
	fmt.Printf("\tI/O需要 %d 才能完成\n\n", *ioTime)
	fmt.Println("任务列表")
	for i, job := range jobs {
		fmt.Printf("\tJob %2d: startTime %3d - runTime %3d - ioFreq %3d\n", i, job.startTime, job.runTime, job.ioFreq)
	}
	fmt.Println()
	if !*solve1 {
		fmt.Println("计算给定工作负载的执行跟踪。")
		fmt.Println("如果您愿意，还可以计算响应和周转时间和每个作业的时间。")
		fmt.Println()
		fmt.Println("完成后，使用-c标志可以获得准确的结果。")
		fmt.Println()
		return
	}

	// queues
	queues = make([][]*Job1, *numQueues)
	for i := 0; i < *numQueues; i++ {
		queues[i] = make([]*Job1, 0)
	}
	curtime := 0
	totalJobs := len(jobs)
	fmt.Println("执行栈：")
	for totalJobs > 0 {
		if *boost > 0 && curtime != 0 {
			if curtime%*boost == 0 {
				fmt.Printf("[ time %d ] BOOST ( every %d )\n", curtime, *boost)
				for i := 0; i < *numQueues-1; i++ {
					for _, job := range queues[i] {
						if !job.doingIO {
							queues[hiQueue] = append(queues[hiQueue], job)
						}
						queues[i] = queues[i][0:0]
					}
				}

				// 提到最高队列，修改每个任务状态
				for i := 0; i < *numJobs; i++ {
					if jobs[i].timeLeft > 0 {
						jobs[i].currPri = hiQueue
						jobs[i].ticksLeft = quantums[hiQueue]
						jobs[i].allotLeft = allotments[hiQueue]
					}
				}
			}

		}

		// 校验io有没有执行完
		if jobStatus, ok := ioDone[curtime]; ok {
			for _, jobStatu := range jobStatus {
				q := jobs[jobStatu.jobIdx].currPri
				jobs[jobStatu.jobIdx].doingIO = false
				fmt.Printf("[ time %d ] 任务 %d 状态 %s\n", curtime, jobStatu.jobIdx, jobStatu.status)
				if !*iobump || jobStatu.status == "JOB BEGINS" {
					queues[q] = append(queues[q], jobs[jobStatu.jobIdx])
				} else {
					queues[q] = append([]*Job1{jobs[jobStatu.jobIdx]}, queues[q]...)
				}
			}
		}

		// 获取当前队列
		curQueue := findQueue()
		if curQueue == -1 {
			fmt.Printf("[ time %d ] IDLE\n", curtime)
			curtime++
			continue
		}

		// 获取当前job
		curJob := queues[curQueue][0]
		if curJob.currPri != curQueue {
			fmt.Printf("currPri[%d] does not match currQueue[%d]\n", curJob.currPri, curQueue)
			os.Exit(1)
		}
		curJob.timeLeft--
		curJob.ticksLeft--

		if curJob.firstRun == -1 {
			curJob.firstRun = curtime
		}

		runTime := curJob.runTime
		ioFreq := curJob.ioFreq
		ticksLeft := curJob.ticksLeft
		allotLeft := curJob.allotLeft
		timeLeft := curJob.timeLeft

		fmt.Printf("[ time %d ] 运行 JOB %d 在队列 %d [ 剩余时间片 %d 剩余次数 %d 剩余时间 %d (of %d) ]\n", curtime, curJob.JobID, curQueue, ticksLeft, allotLeft, timeLeft, runTime)

		if timeLeft < 0 {
			fmt.Println("Error: should never have less than 0 time left to run")
			os.Exit(1)
		}
		curtime++

		if timeLeft == 0 {
			fmt.Printf("[ time %d ] 任务 %d 完成\n", curtime, curJob.JobID)
			totalJobs--
			curJob.endTime = curtime
			queues[curQueue] = queues[curQueue][1:]
			continue
		}

		// io
		issuedIO := false
		if ioFreq > 0 && (((runTime - timeLeft) % ioFreq) == 0) {
			fmt.Printf("[ time %d ] JOB %d 开始执行io任务\n", curtime, curJob.JobID)
			issuedIO = true
			queues[curQueue] = queues[curQueue][1:]
			curJob.doingIO = true
			if *stay {
				// 重置时间
				curJob.ticksLeft = quantums[curQueue]
				curJob.allotLeft = allotments[curQueue]
			}
			futureTime := curtime + *ioTime
			if _, ok := ioDone[futureTime]; !ok {
				ioDone[futureTime] = make([]*jobStatus, 0)
			}
			ioDone[futureTime] = append(ioDone[futureTime], &jobStatus{
				jobIdx: curJob.JobID,
				status: "IO_DONE",
			})
		}

		// 队列时间片用完
		if ticksLeft <= 0 {
			if !issuedIO {
				queues[curQueue] = queues[curQueue][1:]
			}
			curJob.allotLeft--
			if curJob.allotLeft <= 0 {
				if curQueue > 0 {
					curJob.currPri--
					curJob.ticksLeft = quantums[curQueue-1]
					curJob.allotLeft = allotments[curQueue-1]
					if !issuedIO {
						queues[curQueue-1] = append(queues[curQueue-1], curJob)
					}
				} else {
					curJob.ticksLeft = quantums[curQueue]
					curJob.allotLeft = allotments[curQueue]
					if !issuedIO {
						queues[curQueue] = append(queues[curQueue], curJob)
					}
				}
			} else {
				curJob.ticksLeft = quantums[curQueue]
				if !issuedIO {
					queues[curQueue] = append(queues[curQueue], curJob)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("最终统计：")
	responseSum := 0.0
	turnaroundSum := 0.0
	for i, job := range jobs {
		response := job.firstRun - job.startTime
		turnaround := job.endTime - job.startTime
		fmt.Printf("\tJob %2d: startTime %3d - response %3d - turnaround %3d\n", i, job.startTime, response, turnaround)
		responseSum += float64(response)
		turnaroundSum += float64(turnaround)
	}
	fmt.Println()
	fmt.Printf("Avg %2d: startTime  n/a - response %.2f - turnaround %.2f", *numJobs, responseSum/float64(*numJobs), turnaroundSum/float64(*numJobs))
}

// findQueue 获取当前要执行队列
func findQueue() int {
	q := hiQueue
	for q >= 0 {
		if len(queues[q]) > 0 {
			return q
		}
		q--
	}
	return q
}
