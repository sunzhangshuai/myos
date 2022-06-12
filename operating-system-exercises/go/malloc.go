package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Malloc 内存分配类
type Malloc struct {
	size         int
	headerSize   int
	freelist     [][2]int
	sizemap      map[int]int
	policy       string
	returnPolicy string // 排序策略
	coalesce     bool
	align        int
}

// addToMap 地址和大小的映射
func (m *Malloc) addToMap(addr int, size int) {
	m.sizemap[addr] = size
}

// Malloc 分配空间
func (m *Malloc) malloc(size int) (int, int, bool) {
	// 地址对齐
	if m.align != -1 {
		left := size % m.align
		diff := 0
		if left != 0 {
			diff = m.align - left
		}
		left += diff
	}

	// 算上头
	size += m.headerSize

	// 查找最佳空闲块儿
	bestIdx := -1
	bestSize := -1
	if m.policy == "BEST" {
		bestSize = m.size + 1
	}
	bestAddr := -1

	for i, free := range m.freelist {
		eaddr, esize := free[0], free[1]
		if esize >= size && ((m.policy == "BEST" && esize < bestSize) ||
			(m.policy == "WORST" && esize > bestSize) ||
			m.policy == "FIRST") {
			bestAddr = eaddr
			bestSize = esize
			bestIdx = i
			if m.policy == "FIRST" {
				break
			}
		}
	}
	// 说明找到了
	if bestIdx != -1 {
		if bestSize > size {
			// 更新空闲块儿地址
			m.freelist[bestIdx] = [2]int{bestAddr + size, bestSize - size}
			m.addToMap(bestAddr, size)
		} else {
			m.freelist = append(m.freelist[:bestIdx], m.freelist[bestIdx+1:]...)
			m.addToMap(bestAddr, size)
		}
		return bestAddr, bestIdx + 1, true
	}
	return -1, 0, false
}

// free 释放空间
func (m *Malloc) free(addr int) bool {
	size, ok := m.sizemap[addr]
	if !ok {
		return false
	}

	switch m.returnPolicy {
	case "INSERT-BACK":
		m.freelist = append(m.freelist, [2]int{addr, size})
	case "INSERT-FRONT":
		m.freelist = append([][2]int{{addr, size}}, m.freelist...)
	case "ADDRSORT":
		var findFlag bool
		for i, free := range m.freelist {
			if free[0] > addr {
				findFlag = true
				m.freelist = append(m.freelist[:i], append([][2]int{{addr, size}}, m.freelist[i:]...)...)
				break
			}
		}
		if !findFlag {
			m.freelist = append(m.freelist, [2]int{addr, size})
		}
	case "SIZESORT+":
		var findFlag bool
		for i, free := range m.freelist {
			if free[1] > size {
				findFlag = true
				m.freelist = append(m.freelist[:i], append([][2]int{{addr, size}}, m.freelist[i:]...)...)
				break
			}
		}
		if !findFlag {
			m.freelist = append(m.freelist, [2]int{addr, size})
		}
	case "SIZESORT-":
		var findFlag bool
		for i, free := range m.freelist {
			if free[1] < size {
				findFlag = true
				m.freelist = append(m.freelist[:i], append([][2]int{{addr, size}}, m.freelist[i:]...)...)
				break
			}
		}
		if !findFlag {
			m.freelist = append(m.freelist, [2]int{addr, size})
		}
	}

	// 合并，按地址排有意义
	if m.coalesce && m.returnPolicy == "ADDRSORT" {
		newList := make([][2]int, 0, len(m.freelist))
		curr := m.freelist[0]

		for i := 1; i < len(m.freelist); i++ {
			eaddr, esize := m.freelist[i][0], m.freelist[i][1]
			if eaddr == curr[0]+curr[1] {
				curr = [2]int{curr[0], curr[1] + esize}
			} else {
				newList = append(newList, curr)
				curr = m.freelist[i]
			}
		}
		newList = append(newList, curr)
		m.freelist = newList
	}

	// 删掉索引
	delete(m.sizemap, addr)
	return true
}

// dump 打印
func (m *Malloc) dump() {
	fmt.Printf("可用列表 [ 长度 %d ]: ", len(m.freelist))
	for _, free := range m.freelist {
		fmt.Printf("[ 地址:%d 长度:%d ]'", free[0], free[1])
	}
	fmt.Println()
}

var seed6 = flag.Int64("s", 0, "指定随机种子")
var heapSize = flag.Int("S", 100, "堆大小")
var baseAddr = flag.Int("b", 1000, "堆的开始地址")
var headerSize = flag.Int("H", 0, "header块大小")
var alignment = flag.Int("a", -1, "分配对齐单元；-1->不对齐")
var policy1 = flag.String("p", "BEST", "空闲空间搜索算法 (BEST【最优】, WORST【最差】, FIRST【首次】)")
var order = flag.String("l", "ADDRSORT", "空闲列表排序 (ADDRSORT【按地址排】, SIZESORT+【大小升序】, SIZESORT-【大小降序】, INSERT-FRONT【头插】, INSERT-BACK【尾插】)")
var coalesce = flag.Bool("C", false, "合并空闲列表")
var opsNum = flag.Int("n", 10, "要生成的随机操作的数量")
var opsRange = flag.Int("r", 10, "最大分配空间大小")
var opsPAlloc = flag.Int("P", 50, "分配空间操作的百分比，其他是释放")
var opsList = flag.String("A", "", "不随机分配操作, 指定操作列表(+10,-0,etc)")
var solve6 = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()

	rand.Seed(*seed6)

	fmt.Println("随机种子", *seed6)
	fmt.Println("堆大小", *heapSize)
	fmt.Println("堆的开始地址", *baseAddr)
	fmt.Println("header块大小", *headerSize)
	fmt.Println("分配对齐单元；-1->不对齐", *alignment)
	fmt.Println("空闲空间搜索算法", *policy1)
	fmt.Println("空闲列表排序方式", *order)
	fmt.Println("是否合并空闲列表", *coalesce)
	fmt.Println("生成的随机操作的数量", *opsNum)
	fmt.Println("最大分配空间大小", *opsRange)
	fmt.Println("分配空间操作的百分比，其他是释放", *opsPAlloc)
	fmt.Println("操作列表", *opsList)
	fmt.Println("计算答案", *solve6)

	m := &Malloc{
		size:       *heapSize,
		headerSize: *headerSize,
		freelist: [][2]int{
			{*baseAddr, *heapSize},
		},
		sizemap:      make(map[int]int),
		policy:       *policy1,
		returnPolicy: *order,
		coalesce:     *coalesce,
		align:        *alignment,
	}

	p := make([][2]int, 0, *opsNum)

	if *opsList == "" {
		var c int
		for i := 0; i < *opsNum; i++ {
			pr := false
			if rand.Intn(100) < *opsPAlloc {
				size := rand.Intn(*opsRange) + 1
				ptr, cnt, ok := m.malloc(size)
				if ok {
					p = append(p, [2]int{c, ptr})
				}
				fmt.Printf("第[%d]题：分配(%d)空间：", c, size)
				if *solve6 {
					fmt.Printf("分配地址 %d， (找到第 %d 块儿)\n", ptr+*headerSize, cnt)
				} else {
					fmt.Printf("分配地址 ?\n")
				}
				c += 1
				pr = true
			} else {
				pLen := len(p)
				if pLen > 0 {
					idx := rand.Intn(pLen)
					rc := m.free(p[idx][1])
					fmt.Printf("释放(第[%d]个)", p[idx][0])
					if *solve6 {
						fmt.Printf("释放结果 %t\n", rc)
					} else {
						fmt.Printf("释放结果 ?\n")
					}
					p = append(p[:idx], p[idx+1:]...)
					pr = true
				}
			}

			if pr {
				if *solve6 {
					m.dump()
				} else {
					fmt.Println("List？")
				}
			}
			fmt.Println()
		}
		return
	}

	// 自定义
	var c int
	for _, opsItem := range strings.Split(*opsList, ",") {
		if opsItem[0] == '+' {
			size, _ := strconv.Atoi(opsItem[1:])
			ptr, cnt, ok := m.malloc(size)
			if ok {
				p = append(p, [2]int{c, ptr})
			}
			fmt.Printf("第[%d]题：分配(%d)空间：", c, size)
			if *solve6 {
				fmt.Printf("分配地址 %d， (找到第 %d 块儿)\n", ptr+*headerSize, cnt)
			} else {
				fmt.Printf("分配地址 ?\n")
			}
			c += 1
		} else if opsItem[0] == '-' {
			idx, _ := strconv.Atoi(opsItem[1:])
			if idx > len(p) {
				fmt.Println("错误的释放：操作跳过")
				continue
			}
			idx--
			rc := m.free(p[idx][1])
			fmt.Printf("释放(第[%d]个)", p[idx][0])
			if *solve6 {
				fmt.Printf("释放结果 %t\n", rc)
			} else {
				fmt.Printf("释放结果 ?\n")
			}
		} else {
			fmt.Println("参数格式有误")
		}
		if *solve6 {
			m.dump()
		} else {
			fmt.Println("List？")
		}
		fmt.Println()
	}
}
