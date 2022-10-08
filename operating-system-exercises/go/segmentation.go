package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var seed5 = flag.Int64("s", 0, "指定随机种子")
var addresses1 = flag.String("A", "-1", "要访问的一组逗号分隔的页面-1表示随机生成")
var asizestr1 = flag.String("a", "1k", "地址空间大小 (e.g., 16, 64k, 32m)")
var psizestr1 = flag.String("p", "16k", "物理内存大小(e.g., 16, 64k)")
var num1 = flag.Int("n", 5, "生成虚拟地址数量")
var basestr1 = flag.String("b", "-1", "段1基址寄存器的值")
var limitstr1 = flag.String("l", "-1", "段1界限寄存器的值")
var basestr2 = flag.String("B", "-1", "段2基址寄存器的值")
var limitstr2 = flag.String("L", "-1", "段2界限寄存器的值")
var solve5 = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()

	fmt.Println()
	fmt.Println("参数 随机数种子", *seed5)
	fmt.Println("参数 虚拟地址大小", *asizestr1)
	fmt.Println("参数 物理内存大小", *psizestr1)
	fmt.Println()

	rand.Seed(*seed5)
	asize := convert1(*asizestr1)
	psize := convert1(*psizestr1)
	base1 := convert1(*basestr1)
	limit1 := convert1(*limitstr1)
	base2 := convert1(*basestr2)
	limit2 := convert1(*limitstr2)

	if limit1 == -1 {
		limit1 = asize/4 + rand.Intn(asize/4)
	}

	if base1 == -1 {
		for done := 0; done == 0; done++ {
			base1 = rand.Intn(psize)
			if base1+limit1 < psize {
				done = 1
			}
		}
	}

	if limit2 == -1 {
		limit2 = asize/4 + rand.Intn(asize/4)
	}

	if base2 == -1 {
		for done := 0; done == 0; done++ {
			base2 = rand.Intn(psize)
			if base2+limit2 < psize {
				if (base2 > (base1 + limit1)) || ((base2 + limit2) < base1) {
					done = 1
				}
			}
		}
	} else {
		base2 -= limit2
	}
	nbase2 := base2 + limit2

	if limit1 > asize/2 || limit2 > asize/2 {
		fmt.Println("错误：边界寄存器对于此地址空间太大")
		os.Exit(1)
	}

	if (limit1+base1) > base2 && base2 > base1 {
		fmt.Println("错误：段在物理内存中重叠")
		os.Exit(1)
	}

	fmt.Println("段寄存器信息：")
	fmt.Println()
	fmt.Printf("\t基址1: 0x%o (decimal %d)\n", base1, base1)
	fmt.Printf("\t边界1: %d", limit1)
	fmt.Println()
	fmt.Printf("\t基址2: 0x%o (decimal %d)\n", base2+limit2, base2+limit2)
	fmt.Printf("\t边界2: %d", limit2)
	fmt.Println()

	var addresses []int
	if *addresses1 == "-1" {
		addresses = make([]int, *num1)
		for i := 0; i < *num1; i++ {
			addresses[i] = rand.Intn(asize)
		}
	} else {
		adds := strings.Split(*addresses1, ",")
		addresses = make([]int, len(adds))
		for i, address := range adds {
			addresses[i], _ = strconv.Atoi(address)
		}
	}

	fmt.Println("虚拟地址调用栈")

	if !*solve5 {
		for i, address := range addresses {
			fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> PA或段冲突\n", i, address, address)
		}
	}

	if !*solve5 {
		fmt.Println("对于每个虚拟地址，要么写下它转换为的物理地址，要么写下它是一个越界地址（分段冲突）。对于这个问题，您应该假设一个具有两个段的简单地址空间：因此，可以使用虚拟地址的顶部位来检查虚拟地址是在段1（topbit=1）还是在段2（topbit=2）中。请注意，根据段的不同，给您的基/极限对以不同的方向增长，即段1以正方向增长，而段2以负方向增长。")
		fmt.Println()
		return
	}

	for i, address := range addresses {
		if address >= (asize / 2) {
			// 段2
			paddr := nbase2 - (asize - address)
			if paddr < base2 {
				fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> 段地址不合法\n", i, address, address)
			} else {
				fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> 合法: 0x%o (decimal: %4d)\n", i, address, address, paddr, paddr)
			}
		} else {
			// 段1
			if address >= limit1 {
				fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> 段地址不合法\n", i, address, address)
			} else {
				paddr := base1 + address
				fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> 合法: 0x%o (decimal: %4d)\n", i, address, address, paddr, paddr)
			}
		}
	}
}

// convert1 大小转化
func convert1(size string) int {
	var res int
	length := len(size)
	lastchar := size[length-1]
	switch lastchar {
	case 'k', 'K':
		res, _ = strconv.Atoi(size[:length-1])
		res *= 1024
	case 'm', 'M':
		res, _ = strconv.Atoi(size[:length-1])
		res *= 1024 * 1024
	case 'g', 'G':
		res, _ = strconv.Atoi(size[:length-1])
		res *= 1024 * 1024 * 1024
	default:
		res, _ = strconv.Atoi(size)
	}
	return res
}
