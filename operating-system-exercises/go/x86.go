package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var seed8 = flag.Int64("s", 0, "随机种子")
var numthreads = flag.Int("t", 2, "线程数")
var progfile = flag.String("p", "", "源程序 (in .s)")
var intfreq = flag.Int("i", 50, "中断周期")
var procsched = flag.String("P", "", "准确控制何时运行哪个线程")
var intrand = flag.Bool("r", false, "中断周期是否随机")
var argvStr = flag.String("a", "", "逗号分隔每个线程参数(例如: ax=1,ax=2 设置线程0 ax 寄存器为1,线程1 ax 寄存器为2)，通过冒号分隔列表为每个线程指定多个寄存器(例如，ax=1:bx=2,cx=3设 置线程0 ax和bx，对于线程1只设置cx)")
var loadaddr = flag.Int("l", 1000, "加载代码的地址")
var memsize = flag.Int("m", 100, "地址空间大小(KB)")
var memtraceStr = flag.String("M", "", "以逗号分隔的要跟踪的地址列表 (例如:20000,20001)")
var regtraceStr = flag.String("R", "", "以逗号分隔的要跟踪的寄存器列表 (例如:ax,bx,cx,dx)")
var cctrace = flag.Bool("C", false, "是否跟踪条件代码(condition codes)")
var printstats = flag.Bool("S", false, "打印额外状态")
var verbose1 = flag.Bool("v", false, "打印额外信息")
var headercount = flag.Int("H", -1, "打印行标题的频率")
var solve8 = flag.Bool("c", false, "计算结果")

const (
	// COND_GT 条件
	COND_GT = iota
	COND_GTE
	COND_LT
	COND_LTE
	COND_EQ
	COND_NEQ
)

const (
	// REG_ZERO 寄存器
	REG_ZERO = iota
	REG_AX
	REG_BX
	REG_CX
	REG_DX
	REG_EX
	REG_FX
	REG_SP
	REG_BP
)

// doSpace 加空格
func doSpace(howmuch int) {
	for i := 0; i < howmuch; i++ {
		fmt.Printf("%24s", " ")
	}
}

// zassert 断言
func zassert(cond bool, str string) {
	if !cond {
		fmt.Println("出错了:", str)
		os.Exit(1)
	}
	return
}

// Cpu cpu类型
type Cpu struct {
	ccTrace    bool                // 是否监控
	compute    bool                // 是否计算
	verbose    bool                // 是否打印额外信息
	labels     map[string]int      // 标签，循环点或者函数点
	vars       map[string]int      // 静态变量
	maxMemory  int                 // 最大内存
	memory     map[int]interface{} // 内存列表
	pMemory    map[int]string      // 指令内存
	memTrace   []string            // 要监控的内存列表
	condType   []int               // 条件类型
	conditions map[int]bool        // 条件列表
	regNums    []int               // 寄存器号列表
	regNames   map[string]int      // 寄存器名称映射
	registers  map[int]int         // 寄存器列表
	regTrace   []int               // 要监控的寄存器列表
	hdrCount   int                 // 打印行标题的频率
	PC         int                 // 程序计数器
}

// NewCPU 获取新的cpu
func NewCPU(memSize int, memtrace []string, regtrace []string, cctrace, compute, verbose bool, hdrCount int) *Cpu {
	ret := Cpu{
		maxMemory: memSize * 1024,
		memTrace:  memtrace,
		regTrace:  make([]int, len(regtrace)),
		ccTrace:   cctrace,
		compute:   compute,
		verbose:   verbose,
		condType:  []int{COND_GT, COND_GTE, COND_LT, COND_LTE, COND_EQ, COND_NEQ},
		regNums:   []int{REG_ZERO, REG_AX, REG_BX, REG_CX, REG_DX, REG_EX, REG_FX, REG_SP, REG_BP},
		regNames: map[string]int{
			"zero": REG_ZERO,
			"ax":   REG_AX,
			"bx":   REG_BX,
			"cx":   REG_CX,
			"dx":   REG_DX,
			"ex":   REG_EX,
			"fx":   REG_FX,
			"sp":   REG_SP,
			"bp":   REG_BP,
		},
		hdrCount: hdrCount,
		pMemory:  make(map[int]string),
	}
	ret.conditions = make(map[int]bool, len(ret.condType))
	ret.memory = make(map[int]interface{}, ret.maxMemory)
	ret.registers = make(map[int]int, len(ret.regNums))
	ret.labels = make(map[string]int)
	ret.vars = make(map[string]int)
	for i, r := range regtrace {
		ret.regTrace[i] = ret.getRegNum(r)
	}

	ret.initMemory()
	ret.initRegisters()
	ret.initConditionCodes()

	return &ret
}

// initConditionCodes 初始化条件
func (c *Cpu) initConditionCodes() {
	for _, cond := range c.condType {
		c.conditions[cond] = false
	}
}

// initMemory 初始化内存
func (c *Cpu) initMemory() {
	for i := 0; i < c.maxMemory; i++ {
		c.memory[i] = 0
	}
}

// initRegisters 初始化寄存器
func (c *Cpu) initRegisters() {
	for _, r := range c.regNums {
		c.registers[r] = 0
	}
}

// dumpMemory 打印内存
func (c *Cpu) dumpMemory() {
	fmt.Println("内存打印：")
	for m, v := range c.memory {
		if _, ok := c.pMemory[m]; !ok && v != 0 {
			fmt.Printf("\tm[%d]\n", v)
		}
	}
}

// getRegNum 获取寄存器号
func (c *Cpu) getRegNum(name string) int {
	if _, ok := c.regNames[name]; !ok {
		zassert(false, fmt.Sprintf("寄存器 %s 不是合法的寄存器", name))
	}
	return c.regNames[name]
}

// getRegName 获取寄存器名称
func (c *Cpu) getRegName(num int) string {
	for name, n := range c.regNames {
		if n == num {
			return name
		}
	}
	return ""
}

// getRegNums 获取寄存器号列表
func (c *Cpu) getRegNums() []int {
	return c.regNums
}

// getCondList 获取条件标记列表
func (c *Cpu) getCondList() []int {
	return c.condType
}

// getReg 获取寄存器值
func (c *Cpu) getReg(reg int) int {
	return c.registers[reg]
}

// getCond 获取条件值
func (c *Cpu) getCond(cond int) bool {
	return c.conditions[cond]
}

// getPC 获取程序计数器
func (c *Cpu) getPC() int {
	return c.PC
}

// setReg 寄存器赋值
func (c *Cpu) setReg(reg, value int) {
	c.registers[reg] = value
}

// setCond 条件赋值
func (c *Cpu) setCond(cond int, value bool) {
	c.conditions[cond] = value
}

// setPC 程序计数器赋值
func (c *Cpu) setPC(pc int) {
	c.PC = pc
}

// *************************************************** 指令相关 ***************************************************

// halt 中止
func (c *Cpu) halt() int {
	return -1
}

// yield 让出
func (c *Cpu) yield() int {
	return -2
}

// nop 空指令
func (c *Cpu) nop() int {
	return 0
}

// rDump 打印寄存器
func (c *Cpu) rDump() int {
	fmt.Println("寄存器：")
	fmt.Println("ax:", c.registers[REG_AX])
	fmt.Println("bx:", c.registers[REG_BX])
	fmt.Println("cx:", c.registers[REG_CX])
	fmt.Println("dx:", c.registers[REG_DX])
	return 0
}

// mDump 打印内存
func (c *Cpu) mDump(idx int) int {
	fmt.Printf("  m[%d] %d\n", idx, c.memory[idx])
	return 0
}

// moveV2R 将立即数放到寄存器中
func (c *Cpu) moveI2R(src, dst int) int {
	c.registers[dst] = src
	return 0
}

// moveI2M 将立即数写到内存中
func (c *Cpu) moveI2M(src, value, reg1, reg2 int) int {
	tmp := value + c.registers[reg1] + c.registers[reg2]
	c.memory[tmp] = src
	return 0
}

// moveI2m 将内存中的数写入寄存器
func (c *Cpu) moveM2R(value, reg1, reg2, dst int) int {
	tmp := value + c.registers[reg1] + c.registers[reg2]
	c.registers[dst] = c.memory[tmp].(int)
	return 0
}

// moveR2M 将寄存器的值写入内存
func (c *Cpu) moveR2M(src, value, reg1, reg2 int) int {
	tmp := value + c.registers[reg1] + c.registers[reg2]
	c.memory[tmp] = c.registers[src]
	return 0
}

// moveR2R 从寄存器写到寄存器
func (c *Cpu) moveR2R(src, dst int) int {
	c.registers[dst] = c.registers[src]
	return 0
}

// leaM2R 加载有效地址（除了内存值的最终更改以外的所有内容）
func (c *Cpu) leaM2R(value, reg1, reg2, dst int) int {
	tmp := value + c.registers[reg1] + c.registers[reg2]
	c.registers[dst] = tmp
	return 0
}

// addI2R 寄存器加立即数
func (c *Cpu) addI2R(src, dst int) int {
	c.registers[dst] += src
	return 0
}

// addR2R 寄存器的值加一个寄存器值
func (c *Cpu) addR2R(src, dst int) int {
	c.registers[dst] += c.registers[src]
	return 0
}

// subR2R 寄存器减立即数
func (c *Cpu) subI2R(src, dst int) int {
	c.registers[dst] -= src
	return 0
}

// subR2R 寄存器的值减一个寄存器值
func (c *Cpu) subR2R(src, dst int) int {
	c.registers[dst] -= c.registers[src]
	return 0
}

// negR 寄存器值取反
func (c *Cpu) negR(src int) int {
	c.registers[src] = -c.registers[src]
	return 0
}

// *************************************************** 锁相关 ***************************************************

// atomicExchange 原子交换
func (c *Cpu) atomicExchange(src, value, reg1, reg2 int) int {
	tmp := value + c.registers[reg1] + c.registers[reg2]
	old := c.memory[tmp]
	c.memory[tmp] = c.registers[src]
	c.registers[src] = old.(int)
	return 0
}

// fetchAdd 原子增加
func (c *Cpu) fetchAdd(src, value, reg1, reg2 int) int {
	tmp := value + c.registers[reg1] + c.registers[reg2]
	old := c.memory[tmp]
	c.memory[tmp] = c.memory[tmp].(int) + c.registers[src]
	c.registers[src] = old.(int)
	return 0
}

// *************************************************** 条件测试 ***************************************************

func (c *Cpu) testAll(src, dst int) {
	c.initConditionCodes()
	if dst > src {
		c.conditions[COND_GT] = true
	}
	if dst >= src {
		c.conditions[COND_GTE] = true
	}
	if dst < src {
		c.conditions[COND_LT] = true
	}
	if dst <= src {
		c.conditions[COND_LTE] = true
	}
	if dst == src {
		c.conditions[COND_EQ] = true
	}
	if dst != src {
		c.conditions[COND_NEQ] = true
	}
}

// testIR 测试立即数和寄存器
func (c *Cpu) testIR(src, dst int) int {
	c.initConditionCodes()
	c.testAll(src, c.getReg(dst))
	return 0
}

// testRI 测试寄存器和立即数
func (c *Cpu) testRI(src, dst int) int {
	c.initConditionCodes()
	c.testAll(c.getReg(src), dst)
	return 0
}

// testRR 测试寄存器和寄存器
func (c *Cpu) testRR(src, dst int) int {
	c.initConditionCodes()
	c.testAll(c.getReg(src), c.getReg(dst))
	return 0
}

// *************************************************** 跳转 ***************************************************

// jump 跳转
func (c *Cpu) jump(targ int) int {
	c.PC = targ
	return 0
}

// jumpNotEqual 不等则跳转
func (c *Cpu) jumpNotEqual(targ int) int {
	if c.conditions[COND_NEQ] {
		c.PC = targ
	}
	return 0
}

// jumpEqual 相等则跳转
func (c *Cpu) jumpEqual(targ int) int {
	if c.conditions[COND_EQ] {
		c.PC = targ
	}
	return 0
}

// jumpLessThen 小于则跳转
func (c *Cpu) jumpLessThen(targ int) int {
	if c.conditions[COND_LT] {
		c.PC = targ
	}
	return 0
}

// jumpLessThenOrEqual 小于等于则跳转
func (c *Cpu) jumpLessThenOrEqual(targ int) int {
	if c.conditions[COND_LTE] {
		c.PC = targ
	}
	return 0
}

// jumpGreaterThen 大于则跳转
func (c *Cpu) jumpGreaterThen(targ int) int {
	if c.conditions[COND_GT] {
		c.PC = targ
	}
	return 0
}

// jumpGreaterThenOrEqual 大于等于则跳转
func (c *Cpu) jumpGreaterThenOrEqual(targ int) int {
	if c.conditions[COND_GTE] {
		c.PC = targ
	}
	return 0
}

// *************************************************** 调用或返回 ***************************************************

// call 函数调用
func (c *Cpu) call(targ int) int {
	c.registers[REG_SP] -= 4
	c.memory[c.registers[REG_SP]] = c.PC
	c.PC = targ
	return 0
}

// ret 函数返回
func (c *Cpu) ret() int {
	c.PC = c.memory[c.registers[REG_SP]].(int)
	c.registers[REG_SP] += 4
	return 0
}

// *************************************************** 栈相关 ***************************************************

// pushR 寄存器入栈
func (c *Cpu) pushR(reg int) int {
	c.registers[REG_SP] -= 4
	c.memory[c.registers[REG_SP]] = c.registers[reg]
	return 0
}

// pushM 内存入栈
func (c *Cpu) pushM(value, reg1, reg2 int) int {
	c.registers[REG_SP] -= 4
	tmp := value + c.registers[reg1] + c.registers[reg2]
	c.memory[c.registers[REG_SP]] = tmp
	return 0
}

// pop 出栈
func (c *Cpu) pop() int {
	c.registers[REG_SP] += 4
	return 0
}

// popR 出栈到寄存器
func (c *Cpu) popR(dst int) int {
	c.registers[dst] = c.registers[REG_SP]
	c.registers[REG_SP] += 4
	return 0
}

// getArg 获取参数值
// 解析参数，返回 (value, type)
// type：(TYPE_REGISTER：寄存器号, TYPE_IMMEDIATE：立即数, TYPE_MEMORY：内存地址, TYPE_LABEL)
// FORMATS
//    %ax           - register
//    $10           - immediate
//    10            - direct memory
//    10(%ax)       - memory + reg indirect
//    10(%ax,%bx)   - memory + 2 reg indirect
//    10(%ax,%bx,4) - XXX (not handled)
func (c *Cpu) getArg(arg string) (interface{}, string) {
	switch arg[0] {
	case '$':
		zassert(len(arg) == 2, fmt.Sprintf("正确的形式是$number (not %s)", arg))
		value, _ := strconv.Atoi(arg[1:])
		return value, "TYPE_IMMEDIATE"
	case '%':
		return c.getRegNum(arg[1:]), "TYPE_REGISTER"
	case '(':
		register := strings.Split(strings.Split(strings.Split(arg, "(")[1], ")")[0], "%")[1]
		return fmt.Sprintf("%d,%d,%d", 0, c.getRegNum(register), c.getRegNum("zero")), "TYPE_REGISTER"
	case '.':
		return arg, "TYPE_LABEL"
	default:
		// 变量
		if match, _ := regexp.MatchString(`^[A-Za-z]+$`, string(arg[0])); match {
			value, ok := c.vars[arg]
			zassert(ok, fmt.Sprintf("变量 %s 未定义", arg))
			return fmt.Sprintf("%d,%d,%d", value, c.getRegNum("zero"), c.getRegNum("zero")), "TYPE_MEMORY"
		} else if match, _ = regexp.MatchString(`^[\d-]+$`, string(arg[0])); match {
			// 变量
			// 内存地址
			neg := 1
			if arg[0] == '-' {
				arg = arg[1:]
				neg = -1
			}
			s := strings.Split(arg, "(")
			if len(s) == 1 {
				value, _ := strconv.Atoi(arg)
				return fmt.Sprintf("%d,%d,%d", value*neg, c.getRegNum("zero"), c.getRegNum("zero")), "TYPE_MEMORY"
			} else if len(s) == 2 {
				value, _ := strconv.Atoi(s[0])
				t := strings.Split(strings.Split(s[1], ")")[0], ",")
				if len(t) == 1 {
					register := t[0][1:]
					return fmt.Sprintf("%d,%d,%d", value*neg, c.getRegNum(register), c.getRegNum("zero")), "TYPE_MEMORY"
				} else if len(t) == 2 {
					register1 := t[0][1:]
					register2 := t[1][1:]
					return fmt.Sprintf("%d,%d,%d", value*neg, c.getRegNum(register1), c.getRegNum(register2)), "TYPE_MEMORY"
				}
			}
		} else {
			zassert(false, fmt.Sprintf("指令：错误的参数[%s]", arg))
		}
	}
	zassert(false, fmt.Sprintf("指令：错误的参数[%s]", arg))
	return "", ""
}

// load 加载指令
func (c *Cpu) load(infile string, loadaddr int) {
	pc := loadaddr

	// 程序地址
	bpc := loadaddr

	// 静态数据地址
	data := 100

	// 处理标签和变量
	file, _ := os.Open(infile)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}

		line = strings.Split(line, "#")[0]
		if len(strings.Trim(line, "")) == 0 {
			continue
		}
		tmp := strings.Split(line, " ")
		if len(tmp) == 0 {
			continue
		}

		// 只注意标签和变量
		if tmp[0] == ".var" {
			c.vars[tmp[1]] = data
			if len(tmp) == 3 {
				mul, _ := strconv.Atoi(tmp[2])
				data += mul
			}
			data += 4
			zassert(data < bpc, "静态数据导致加载地址溢出")
			if c.verbose {
				fmt.Println("分配变量：", tmp[1], "-->", c.vars[tmp[1]])
			}
		} else if tmp[0][0] == '.' {
			c.labels[tmp[0]] = pc
			if c.verbose {
				fmt.Println("分配变量：", tmp[0], "-->", pc)
			}
		} else {
			pc++
		}
	}
	file.Close()
	if c.verbose {
		fmt.Println()
	}

	// 处理指令
	pc = loadaddr
	file, _ = os.Open(infile)
	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		line = strings.Split(line, "#")[0]
		if len(strings.Trim(line, "")) == 0 {
			continue
		}
		tmp := strings.SplitN(line, " ", 2)
		if len(tmp) == 0 {
			continue
		}

		if line[0] == '.' {
			continue
		}

		opcode := tmp[0]
		c.pMemory[pc] = line

		// 区分操作码
		switch opcode {
		case "mov": // load 或 store
			rtmp := strings.SplitN(tmp[1], ", ", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("mov：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			if stype == "TYPE_MEMORY" && dtype == "TYPE_MEMORY" {
				fmt.Println("mov 错误：两个内存地址")
				os.Exit(1)
			} else if stype == "TYPE_IMMEDIATE" && dtype == "TYPE_IMMEDIATE" {
				fmt.Println("mov 错误：两个立即数")
				os.Exit(1)
			} else if stype == "TYPE_IMMEDIATE" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.moveI2R(src.(int), dst.(int))
				}
			} else if stype == "TYPE_MEMORY" && dtype == "TYPE_REGISTER" {
				tmp = strings.Split(src.(string), ",")
				c.memory[pc] = func() int {
					i1, _ := strconv.Atoi(tmp[0])
					i2, _ := strconv.Atoi(tmp[1])
					i3, _ := strconv.Atoi(tmp[2])
					return c.moveM2R(i1, i2, i3, dst.(int))
				}
			} else if stype == "TYPE_REGISTER" && dtype == "TYPE_MEMORY" {
				tmp = strings.Split(dst.(string), ",")
				c.memory[pc] = func() int {
					i1, _ := strconv.Atoi(tmp[0])
					i2, _ := strconv.Atoi(tmp[1])
					i3, _ := strconv.Atoi(tmp[2])
					return c.moveR2M(src.(int), i1, i2, i3)
				}
			} else if stype == "TYPE_REGISTER" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.moveR2R(src.(int), dst.(int))
				}
			} else if stype == "TYPE_IMMEDIATE" && dtype == "TYPE_MEMORY" {
				tmp = strings.Split(dst.(string), ",")
				c.memory[pc] = func() int {
					i1, _ := strconv.Atoi(tmp[0])
					i2, _ := strconv.Atoi(tmp[1])
					i3, _ := strconv.Atoi(tmp[2])
					return c.moveI2M(src.(int), i1, i2, i3)
				}
			} else {
				zassert(false, "格式错误的mov指令")
			}
		case "lea": // 加载
			rtmp := strings.SplitN(tmp[1], ", ", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("lea：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			zassert(stype == "TYPE_MEMORY" && dtype == "TYPE_REGISTER", "格式错误的lea指令（应为内存地址源以注册目标）")
			tmp = strings.Split(src.(string), ",")
			c.memory[pc] = func() int {
				i1, _ := strconv.Atoi(tmp[0])
				i2, _ := strconv.Atoi(tmp[1])
				i3, _ := strconv.Atoi(tmp[2])
				return c.leaM2R(i1, i2, i3, dst.(int))
			}
		case "neg": // 取反
			zassert(len(tmp) == 2, fmt.Sprintf("neg：需要两个参数，用逗号分隔[%s]", line))
			dst, dtype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(dtype == "TYPE_REGISTER", "只能在寄存器中")
			c.memory[pc] = func() int {
				return c.negR(dst.(int))
			}
		case "pop":
			if len(tmp) == 1 {
				c.memory[pc] = func() int {
					return c.pop()
				}
			} else if len(tmp) == 2 {
				dst, dtype := c.getArg(strings.Trim(tmp[1], " "))
				zassert(dtype == "TYPE_REGISTER", "只能弹出到寄存器中")
				c.memory[pc] = func() int {
					return c.popR(dst.(int))
				}
			} else {
				zassert(false, "pop指令必须有零/一个参数")
			}
		case "push":
			src, stype := c.getArg(strings.Trim(tmp[1], " "))
			if stype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.pushR(src.(int))
				}
			} else if stype == "TYPE_MEMORY" {
				tmp = strings.Split(src.(string), ",")
				c.memory[pc] = func() int {
					i1, _ := strconv.Atoi(tmp[0])
					i2, _ := strconv.Atoi(tmp[1])
					i3, _ := strconv.Atoi(tmp[2])
					return c.pushM(i1, i2, i3)
				}
			} else {
				zassert(false, "push指令有误")
			}
		case "call":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			if ttype == "TYPE_LABEL" {
				c.memory[pc] = func() int {
					return c.call(c.labels[targ.(string)])
				}
			} else {
				zassert(false, "call指令必须是标签")
			}
		case "ret":
			c.memory[pc] = func() int {
				return c.ret()
			}
		case "add":
			rtmp := strings.SplitN(tmp[1], ",", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("add：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			if stype == "TYPE_IMMEDIATE" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.addI2R(src.(int), dst.(int))
				}
			} else if stype == "TYPE_REGISTER" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.addR2R(src.(int), dst.(int))
				}
			} else {
				zassert(false, "add指令的使用格式错误")
			}
		case "sub":
			rtmp := strings.SplitN(tmp[1], ",", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("sub：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			if stype == "TYPE_IMMEDIATE" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.subI2R(src.(int), dst.(int))
				}
			} else if stype == "TYPE_REGISTER" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.subR2R(src.(int), dst.(int))
				}
			} else {
				zassert(false, "sub指令的使用格式错误")
			}
		case "fetchadd":
			rtmp := strings.SplitN(tmp[1], ",", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("fetchadd：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			if stype == "TYPE_REGISTER" && dtype == "TYPE_MEMORY" {
				tmp = strings.Split(dst.(string), ",")
				c.memory[pc] = func() int {
					i1, _ := strconv.Atoi(tmp[0])
					i2, _ := strconv.Atoi(tmp[1])
					i3, _ := strconv.Atoi(tmp[2])
					return c.fetchAdd(src.(int), i1, i2, i3)
				}
			} else {
				zassert(false, "fetchadd指令的使用格式错误")
			}
		case "xchg":
			rtmp := strings.SplitN(tmp[1], ",", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("xchg：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			if stype == "TYPE_REGISTER" && dtype == "TYPE_MEMORY" {
				tmp = strings.Split(dst.(string), ",")
				c.memory[pc] = func() int {
					i1, _ := strconv.Atoi(tmp[0])
					i2, _ := strconv.Atoi(tmp[1])
					i3, _ := strconv.Atoi(tmp[2])
					return c.atomicExchange(src.(int), i1, i2, i3)
				}
			} else {
				zassert(false, "xchg指令的使用格式错误")
			}
		case "test":
			rtmp := strings.SplitN(tmp[1], ",", 2)
			zassert(len(rtmp) == 1 || len(rtmp) == 2, fmt.Sprintf("check：需要两个参数，用逗号分隔[%s]", line))
			src, stype := c.getArg(strings.Trim(rtmp[0], " "))
			dst, dtype := c.getArg(strings.Trim(rtmp[1], " "))
			if stype == "TYPE_IMMEDIATE" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.testIR(src.(int), dst.(int))
				}
			} else if stype == "TYPE_REGISTER" && dtype == "TYPE_REGISTER" {
				c.memory[pc] = func() int {
					return c.testRR(src.(int), dst.(int))
				}
			} else if stype == "TYPE_REGISTER" && dtype == "TYPE_IMMEDIATE" {
				c.memory[pc] = func() int {
					return c.testRI(src.(int), dst.(int))
				}
			} else {
				zassert(false, "check指令的使用格式错误")
			}
		case "j":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jump(c.labels[targ.(string)])
			}
		case "je":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jumpEqual(c.labels[targ.(string)])
			}
		case "jne":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jumpNotEqual(c.labels[targ.(string)])
			}
		case "jlt":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jumpLessThen(c.labels[targ.(string)])
			}
		case "jlte":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jumpLessThenOrEqual(c.labels[targ.(string)])
			}
		case "jgt":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jumpGreaterThen(c.labels[targ.(string)])
			}
		case "jgte":
			targ, ttype := c.getArg(strings.Trim(tmp[1], " "))
			zassert(ttype == "TYPE_LABEL", fmt.Sprintf("错误的跳转指令：[%s]", tmp[1]))
			c.memory[pc] = func() int {
				return c.jumpGreaterThenOrEqual(c.labels[targ.(string)])
			}
		case "nop":
			c.memory[pc] = func() int {
				return c.nop()
			}
		case "halt":
			c.memory[pc] = func() int {
				return c.halt()
			}
		case "yield":
			c.memory[pc] = func() int {
				return c.yield()
			}
		case "mdump":
			c.memory[pc] = func() int {
				idx, _ := strconv.Atoi(tmp[1])
				return c.mDump(idx)
			}
		case "rdump":
			c.memory[pc] = func() int {
				return c.rDump()
			}
		}

		if c.verbose {
			fmt.Printf("pc:%d LOADING %20s\n", pc, c.pMemory[pc])
		}

		pc++
	}

	file.Close()
	if c.verbose {
		fmt.Println()
	}
}

// printHeaders 打印头信息
func (c *Cpu) printHeaders(procs *ProcessList) {
	if len(c.memTrace) > 0 {
		for _, m := range c.memTrace {
			fmt.Printf("%10s", m)
		}
		fmt.Printf(" ")
	}

	if len(c.regTrace) > 0 {
		for _, r := range c.regTrace {
			fmt.Printf("%5s", c.getRegName(r))
		}
		fmt.Printf(" ")
	}

	if c.ccTrace {
		fmt.Printf(">= >  <= <  != ==")
	}

	// 打印线程
	for i := 0; i < procs.getNum(); i++ {
		fmt.Printf("%16s%8s", fmt.Sprintf("Thread %d", i), " ")
	}
	fmt.Println()
}

// printTrace 打印堆栈信息
func (c *Cpu) printTrace(newline bool) {
	if len(c.memTrace) > 0 {
		for _, m := range c.memTrace {
			if c.compute {
				if match, _ := regexp.MatchString(`^[\d-]+$`, m); match {
					idx, _ := strconv.Atoi(m)
					fmt.Printf("%10d", c.memory[idx])
				} else {
					fmt.Printf("%10d", c.memory[c.vars[m]])
				}
			} else {
				fmt.Printf("%10s", "?")
			}
		}
		fmt.Printf(" ")
	}

	if len(c.regTrace) > 0 {
		for _, r := range c.regTrace {
			if c.compute {
				fmt.Printf("%5d", c.registers[r])
			} else {
				fmt.Printf("%5s", "?")
			}
		}
		fmt.Printf(" ")
	}

	if c.ccTrace {
		for _, cond := range c.getCondList() {
			if c.compute {
				fmt.Printf("%5t", c.conditions[cond])
			} else {
				fmt.Printf("%5t", "?")
			}
		}
	}

	if (len(c.memTrace) > 0 || len(c.regTrace) > 0 || c.ccTrace) && newline {
		fmt.Println()
	}
}

// setInt 获取中断周期
func (c *Cpu) setInt(intfreq int, intrand bool) int {
	if !intrand {
		return intfreq
	}
	return rand.Intn(intfreq) + 1
}

// run 运行指令
func (c *Cpu) run(procs *ProcessList, intfreq int, intrand bool) int {
	if procs.manual {
		intfreq = 1
		intrand = false
	}

	// 中断
	interrupt := c.setInt(intfreq, intrand)
	icount := 0
	c.printHeaders(procs)
	c.printTrace(true)

	for true {
		if c.hdrCount > 0 && icount%c.hdrCount == 0 && icount > 0 {
			c.printHeaders(procs)
			c.printTrace(true)
		}

		tid := procs.getCurr().getTid()
		precPc := c.PC
		instruction := c.memory[c.PC]
		c.PC++

		// 执行指令
		rc := instruction.(func() int)()

		// 打印堆栈信息
		c.printTrace(false)

		doSpace(tid)

		fmt.Println(precPc, c.pMemory[precPc])
		icount++

		if rc == -1 {
			procs.done()
			if procs.numDone() == procs.getNum() {
				return icount
			}
			procs.next()
			procs.restore()

			c.printTrace(false)

			for i := 0; i < procs.getNum(); i++ {
				fmt.Printf("----- Halt;Switch ----- ")
			}
			fmt.Println()
		}

		interrupt--
		if interrupt == 0 || rc == -2 {
			curr := procs.getCurr()
			interrupt = c.setInt(intfreq, intrand)
			procs.save()
			procs.next()
			procs.restore()

			next := procs.getCurr()

			if !procs.isManual() || (procs.isManual() && curr != next) {
				c.printTrace(false)
				for i := 0; i < procs.getNum(); i++ {
					fmt.Printf("----- Interrupt ----- ")
				}
				fmt.Println()
			}
		}
	}
	return 0
}

// Process 线程
type Process struct {
	cpu     *Cpu
	tid     int         // 线程id
	pc      int         // 程序计数器
	regs    map[int]int // 寄存器值
	cpuCond map[int]bool
	done    bool
	stack   int
}

// NewProcess 新的线程
func NewProcess(cpu *Cpu, tid, pc, stackbottom int, reginit string) *Process {
	ret := Process{
		cpu:   cpu,
		tid:   tid,
		pc:    pc,
		stack: stackbottom,
	}

	ret.regs = make(map[int]int)
	for _, regNum := range cpu.regNums {
		ret.regs[regNum] = 0
	}
	if reginit != "" {
		for _, r := range strings.Split(reginit, ":") {
			tmp := strings.Split(r, "=")
			ret.regs[cpu.getRegNum(tmp[0])], _ = strconv.Atoi(tmp[1])
		}
	}
	ret.regs[cpu.getRegNum("sp")] = stackbottom

	ret.cpuCond = make(map[int]bool)
	for _, c := range cpu.getCondList() {
		ret.cpuCond[c] = false
	}
	return &ret
}

// getTid 获取线程id
func (p *Process) getTid() int {
	return p.tid
}

// save 保存上下文
func (p *Process) save() {
	p.pc = p.cpu.getPC()
	for _, cond := range p.cpu.getCondList() {
		p.cpuCond[cond] = p.cpu.getCond(cond)
	}

	for _, r := range p.cpu.getRegNums() {
		p.regs[r] = p.cpu.getReg(r)
	}
}

// restore 写入上下文
func (p *Process) restore() {
	p.cpu.setPC(p.pc)
	for _, cond := range p.cpu.getCondList() {
		p.cpu.setCond(cond, p.cpuCond[cond])
	}

	for _, r := range p.cpu.getRegNums() {
		p.cpu.setReg(r, p.regs[r])
	}
}

// setDone 线程结束
func (p *Process) setDone() {
	p.done = true
}

// isDone 线程是否结束
func (p *Process) isDone() bool {
	return p.done
}

// ProcessList 线程列表
type ProcessList struct {
	plist     []*Process
	curr      int // 当前线程
	active    int // 活跃的线程数
	procSched []int
	manual    bool // 手动确定
}

// NewProcessList 新的线程列表
func NewProcessList() *ProcessList {
	return &ProcessList{
		plist:     make([]*Process, 0),
		procSched: make([]int, 0),
		curr:      0,
		active:    0,
	}
}

// finalize 最终执行
func (l *ProcessList) finalize(procSched string) {
	if procSched == "" {
		for i := range l.plist {
			l.procSched = append(l.procSched, i)
		}
		l.curr = 0
		l.restore()
		return
	}

	l.manual = true
	check := make(map[int]bool)

	for i := 0; i < len(procSched); i++ {
		p, _ := strconv.Atoi(string(procSched[i]))
		if p >= l.getNum() {
			zassert(false, fmt.Sprintf("错误的调度：不能包含不存在的线程【%d】", p))
		}
		l.procSched = append(l.procSched, p)
		check[p] = true
	}

	if len(check) != l.active {
		zassert(false, fmt.Sprintf("错误的调度：不包括所有进程【%s】", procSched))
	}
	l.curr = 0
	l.restore()
}

// done 线程结束
func (l *ProcessList) done() {
	l.getCurr().setDone()
	l.active--
}

// numDone 获取结束的线程数量
func (l *ProcessList) numDone() int {
	return len(l.plist) - l.active
}

// getNum 获取线程数量
func (l *ProcessList) getNum() int {
	return len(l.plist)
}

// getNum 获取线程数量
func (l *ProcessList) add(p *Process) {
	l.active++
	l.plist = append(l.plist, p)
}

// isManual 是否手动调度
func (l *ProcessList) isManual() bool {
	return l.manual
}

// getCurr 获取当前线程
func (l *ProcessList) getCurr() *Process {
	return l.plist[l.procSched[l.curr]]
}

// save 保存上下文
func (l *ProcessList) save() {
	l.getCurr().save()
}

// restore 载入上下文
func (l *ProcessList) restore() {
	l.getCurr().restore()
}

// next 下一个线程
func (l *ProcessList) next() {
	for i := l.curr + 1; i < len(l.procSched); i++ {
		if !l.plist[l.procSched[i]].isDone() {
			l.curr = i
			return
		}
	}
	for i := 0; i < l.curr; i++ {
		if !l.plist[l.procSched[i]].isDone() {
			l.curr = i
			return
		}
	}
}

func main() {
	flag.Parse()
	rand.Seed(*seed8)

	fmt.Println("选项 随机种子：", *seed8)
	fmt.Println("选项 线程数：", *numthreads)
	fmt.Println("选项 中断周期：", *intfreq)
	fmt.Println("选项 中断周期是否随机", *intrand)
	fmt.Println("选项 进程控制：", *procsched)
	fmt.Println("选项 线程参数：", *argvStr)
	fmt.Println("选项 加载代码的地址：", *loadaddr)
	fmt.Println("选项 地址空间大小(KB)：", *memsize)
	fmt.Println("选项 跟踪的地址列表：", *memtraceStr)
	fmt.Println("选项 跟踪的寄存器列表", *regtraceStr)
	fmt.Println("选项 打印额外状态", *printstats)
	fmt.Println("选项 打印额外信息", *verbose1)
	fmt.Println()

	argv := strings.Split(*argvStr, ",")

	// 跟踪的地址列表
	memtrace := make([]string, 0)
	if *memtraceStr != "" {
		memtrace = strings.Split(*memtraceStr, ",")
	}

	// 跟踪的寄存器列表
	regtrace := make([]string, 0)
	if *regtraceStr != "" {
		regtrace = strings.Split(*regtraceStr, ",")
	}

	cpu := NewCPU(*memsize, memtrace, regtrace, *cctrace, *solve8, *verbose1, *headercount)
	cpu.load(*progfile, *loadaddr)

	procs := NewProcessList()
	pid := 0
	stack := *memsize * 1000

	var arg string
	for i := 0; i < *numthreads; i++ {
		if len(argv) > 1 {
			arg = argv[pid]
		} else {
			arg = argv[0]
		}
		procs.add(NewProcess(cpu, pid, *loadaddr, stack, arg))
		stack -= 1000
		pid++
	}

	procs.finalize(*procsched)

	satrt := time.Now()
	ic := cpu.run(procs, *intfreq, *intrand)
	end := time.Now().Sub(satrt)

	if *printstats {
		fmt.Println()
		fmt.Printf("统计：指令执行数 %d\n", ic)
		fmt.Printf("统计：指令执行数 %.2f\n", float64(ic)/float64(end)/1000)
	}
}
