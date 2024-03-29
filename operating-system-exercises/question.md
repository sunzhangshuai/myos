# 第四章

**程序  `process-run.py` ，允许查看进程状态在 CPU 上运行时的变化情况**。

1. 用以下标志运行程序：`-l 5:100,5:100`。CPU 利用率（CPU 使用时间的百分比）应该是多少？为什么你知道这一点？利用 `-c` 标记查看你的答案是否正确。

   > 答：CPU 利用率为 100%,程序没有进行 IO。
   >
   > 验证：`python ./python/process-run.py -l 5:100,5:100 -c`。

2. 现在用这些标志运行： `-l 4:100,1:0` 这些标志指定了一个包含 4 条指令的进程(都要使用 CPU)，并且只是简单地发出 IO 并等待它完成。完成这两个进程需要多长时间?利用 `-c` 检查你的答案是否正确。

   > 答：需要 10 个时钟周期，因为io需要等待5。
   >
   > 验证：python ./python/process-run.py。

3. 现在交换进程的顺序： `python ./python/process-run.py -l 1:0,4:100` 现在发生了什么？交换顺序是否重要？为什么？同样，用 `-c` 看看你的答案是否正确。

   > 答：进程 1 （PID1）执行 IO 时,进程 2 使用 CPU，运行时间由 10 变为 6 个时钟周期。
   >
   > 验证：`python ./python/process-run.py -l 1:0,4:100 -c`。

4. 现在探索另一些标志。一个重要的标志是 -S，它决定了当进程发出 IO 时系统如何反应。将标志设置为 SWITCH_ON_END，在进程进行 I/O 操作时,系统将不会切换到另一个进程,而是等待进程完成。当你运行以下两个进程时，会发生什么情况？一个执行 I/O，另一个执行 CPU 工作。（`-l 1:0,4:100 -c -S SWITCH_ON_END`）

   > 答：进程 1 执行 IO 操作时，进程 2 等待。
   >
   > 验证：`python ./python/process-run.py -l 1:0,4:100 -c -S SWITCH_ON_END -c`。

5. 现在,运行相同的进程，但切换行为设置，在等待 IO 时切换到另一个进程（`-l 1:0,4:100 -c -S SWITCH_ON_IO`）现在会发生什么？利用 -c 来确认你的答案是否正确。

   > 答：进程 1 执行 IO 操作时，进程 2 执行,充分利用 CPU 资源。
   >
   > 验证：`python ./python/process-run.py -l 1:0,4:100 -c -S SWITCH_ON_IO`。

6. 另一个重要的行为是 IO 完成时要做什么。利用 `-I IO_RUN_LATER`，当 IO 完成时,发出它的进程不一定马上运行。相反，当时运行的进程一直运行。当你运行这个进程组合时会发生什么？（`-l 3:0,5:100,5:100 -S SWITCH_ON_IO -I IO_RUN_LATER -c -p`）系统资源是否被有效利用？

   > 答：系统资源没有被有效利用，cpu 使用率59.09%。
   >
   > 验证：`python ./python/process-run.py -l 3:0,5:100,5:100 -S SWITCH_ON_IO -I IO_RUN_LATER -c -p`。

7. 现在运行相同的进程,但使用 `-I IO_RUN_IMMEDIATE` 设置,该设置立即运行发出 IO 的进程。这种行为有何不同？为什么运行一个刚刚完成 IO 的进程会是一个好主意？

   > 答： 系统资源被有效利用，当 IO 密集进程 IO 完毕时，立即切换到该进程发出 IO 操作请求，再让出 cpu，使得系统资源被充分利用。
   >
   > 验证：`python ./python/process-run.py -l 3:0,5:100,5:100 -S SWITCH_ON_IO -I IO_RUN_IMMEDIATE -c -p`。

8. 现在运行一些随机生成的进程,例如 `-s 1 -l 3:50,3:50，-s 2 -l 3:50,3:50，-s 3 -l 3:50,3:50` 看看你是否能预测追踪记录会如何变化？当你使用 `-I IO_RUN_IMMEDIATE` 与 `-I IO_RUN_LATER` 时会发生什么？当你使用 `-S SWITCH_ON_IO` 与 `-S SWITCH_ON_END` 时会发生什么?

   > ```shell
   > python ./python/process-run.py -s 1 -l 3:50,3:50 -c -I IO_RUN_IMMEDIATE -S SWITCH_ON_IO
   > python ./python/process-run.py -s 2 -l 3:50,3:50 -c -I IO_RUN_IMMEDIATE -S SWITCH_ON_IO
   > python ./python/process-run.py -s 3 -l 3:50,3:50 -c -I IO_RUN_IMMEDIATE -S SWITCH_ON_IO
   > ```

# 第五章

1. 编写一个调用 fork()的程序。在调用之前,让主进程访问一个变量(例如 x)并将其值设置为某个值(例如 100)。子进程中的变量有什么值?当子进程和父进程都改变 x 的值时,变量会发生什么?

   > 答: 子进程父进程各自一份 x 变量,修改互不影响。
   >
   > 验证：`./c/task1`

2. 编写一个打开文件的程序(使用 open 系统调用),然后调用 fork 创建一个新进程。子进程和父进程都可以访问 open()返回的文件描述符吗?当它们并发(即同时)写入文件时,会发生什么?

   > 答: 子进程和父进程都能访问 fd。存在竞争条件，无法同时使用 fd，但最终都会写入成功。
   >
   > 验证：`./ch5/task2`

3. 使用 fork()编写另一个程序。子进程应打印“hello”,父进程应打印“goodbye”你应该尝试确保子进程始终先打印。你能否不在父进程调用 wait()而做到这一点呢?

   > 答：vfork能保证子进程先执行。不过已经弃用了。

4. 编写一个调用 fork()的程序,然后调用某种形式的 exec()来运行程序"/bin/ls"看看是否可以尝试 exec 的所有变体,包括 execl()、 execle()、 execlp()、 execv()、 execvp()和 execve(),为什么同样的基本调用会有这么多变种？

   > exec 多个变体提供不同的功能。
   >
   > 在 exec 函数族中，后缀 l、v、p、e 添加到 exec 后，所指定的函数将具有某种操作能力：
   >
   > - l: 希望接收以逗号分隔的参数列表,列表以 NULL 指针作为结束标志
   > - v: 希望接收一个以 NULL 结尾字符串数组的指针
   > - p: 是一个以 NULL 结尾的字符串数组指针,函数可以利用 DOS 的 PATH 变量查找自程序文件
   > - e 函数传递指定采纳数 envp(环境变量),允许改变子进程环境,无后缀 e 是,子进程使用当前程序环境
   >
   > c 语言没有默认参数语法,只能实现多个变体。

5. 现在编写一个程序，在父进程中使用 wait(),等待子进程完成。wait()返回什么？如果你在子进程中使用 wait()会发生什么？

   > 答: wait 成功返回子进程 id,执行失败返回-1
   >
   > 子进程调用 wait,执行失败,返回-1。

6. 对前一个程序稍作修改，这次使用 waitpid()而不是 wait()。什么时候 waitpid()会有用？

   > 答: waitpid 提供更多操作,比如提供非阻塞版本 wait。

7. 编写一个创建子进程的程序，然后在子进程中关闭标准输出（STDOUT_FILENO).如果子进程在关闭描述符后调用 printf()打印输出，会发生什么？

   > printf 不会打印到控制台。

8. 编写一个程序，创建两个子进程，并使用 pipe()系统调用，将一个子进程的标准输出连接到另一个子进程的标准输入。

   > 验证：`./c/task8`

# 第六章

测量系统调用和上下文切换的成本。

- **系统调用**

  > `./c/syscall-test`。

- **上下文切换**

  > `todo`。

# 第七章

1. 使用 SJF 和 FIFO 调度程序运行长度为 200 的 3 个作业时,计算响应时间和周转时间。

   > `./go/scheduler -l 200,200,200 -p FIFO -c`。
   >
   > 一样的。

2. 现在做同样的事情,但有不同长度的作业,即 100、200 和 300。

   > `./go/scheduler -l 100,200,300 -p SJF -c`。

3. 现在做同样的事情,但采用 RR 调度程序,时间片为 1。

   > `./go/scheduler -l 100,200,300 -p RR -c`。

4. 对于什么类型的工作负载,SJF 提供与 FIFO 相同的周转时间?

   > 答：运行时间相同。

5. 对于什么类型的工作负载和量子长度(时间片长度),SJF 与 RR 提供相同的响应时间?

   > 答：运行时间 <= 时间片。

6. 随着工作长度的增加,SJF 的响应时间会怎样?你能使用模拟程序来展示趋势吗?

   > `./go/scheduler -l 200,400,600 -p SJF -c`
   >
   > 答：响应时间越来越长。

7. 随着量子长度(时间片长度)的增加,RR 的响应时间会怎样?你能写出一个方程,计算给定 N 个工作时,最坏情况的响应时间吗?

   > 答：res_avg = (0 + t1 + (t1+t2) + (t1+t2+t3) + ... (t1+t2+t3 +...tN-1))/N
   >
   > 平均响应时间增加。

# 第八章

1. 只用两个工作和两个队列运行几个随机生成的问题。针对每个工作计算 MLFQ 的执行记录。限制每项作业的长度并关闭 I/O,让你的生活更轻松。

   > `./go/mlfq -j 2 -n 2 -M 0 -m 15 -s 1 -c`。

2. 如何运行调度程序来重现本章中的每个实例？

   > `./go/mlfq -l 0,200,0 -Q 10,10,0 -c`。

3. 将如何配置调度程序参数，像轮转调度程序那样工作？

   > `./go/mlfq -q 1 -c`。

4. 设计两个工作的负载和调度程序参数，以便一个工作利用较早的规则 4a 和 4b(用-S 标志打开）来“愚弄”调度程序，在特定的时间间隔内获得 99%的 CPU。

   > 设置-S 参数, 并使时间片长度 >= IO 频率即可。
   >
   > `./go/mlfq -S -q 11 -l 0,200,9:0,200,9 -c`。

5. 给定一个系统，其最高队列中的时间片长度为 10ms,你需要如何频繁地将工作推回到最高优先级级别（带有-B 标志），以保证一个长时间运行（并可能饥饿）的工作得到至少 5%的 CPU?

   > -B 参数 <= 190 即可，即至少每 200ms 运行 10ms
   >
   > `./go/mlfq -B 190 -m 200 -c`。

6. 调度中有一个问题，即刚完成 I/O 的作业添加在队列的哪一端。-I 标志改变了这个调度模拟器的这方面行为。尝试一些工作负载，看看你是否能看到这个标志的效果。

   > `./go/mlfq -I -c`。

# 第九章

1. 计算 3 个工作在随机种子为 1、2 和 3 时的模拟解。

   > `./go/lottery -s 1 -c`
   >
   > `./go/lottery -s 2 -c`
   >
   > `./go/lottery -s 3 -c`

2. 现在运行两个具体的工作：每个长度为 10,但是一个（工作 0)只有一张彩票，另一个（工作 1)有 100 张（-l 10:1,10:100)

   彩票数量如此不平衡时会发生什么？在工作 1 完成之前，工作 0 是否会运行？多久？ 一般来说，这种彩票不平衡对彩票调度的行为有什么影响？

   > `./go/lottery -l 10:1,10:100 -c`。
   >
   > 导致作业 0 响应时间与周转时间可能非常长，作业 1 完成前,作业 0 会运行,但概率小。彩票不平衡的调度导致彩票数少的作业响应时间与周转时间变长。

3. 如果运行两个长度为 100 的工作，都有 100 张彩票（-l 100:100,100:100),调度程序有多不公平？

   > `./go/lottery -l 100:100,100:100 -c`。
   >
   > 不公平性取决于一项工作比另一项工作早完成多少。时间片 <= 100 时,时间片越小,公平程度越高。

4. 随着量子规模（－q)变大，你对上一个问题的答案如何改变？

   > 时间片 <= 100 时,时间片越小,公平程度越高，>=100 时, 一项工作比另一项工作早完成的时间为 100。

5. 你可以制作类似本章中的图表吗？还有什么值得探讨的？用步长调度程序，图表看起来如何？

# 第十四章

1. 编写一个名为 null.c 的简单程序,它创建一个指向整数的指针,将其设置为 NULL,然后尝试对其进行释放内存操作。把它编译成一个名为 null 的可执行文件。当你运行这个程序时会发生什么?

   > `./c/null`。
   >
   > 答：什么都没发生。

2. 编译该程序,其中包含符号信息(使用-g 标志)。这样做可以将更多信息放入可执行文件中,使调试器可以访问有关变量名等的更多有用信息。通过输入 gdb null 在调试器下运行该程序,然后,一旦 gdb 运行,输入 run.gdb 显示什么信息?

   > ```shell
   > Type "apropos word" to search for commands related to "word"...
   > Reading symbols from a.out...
   > Reading symbols from /Users/sunchen/Documents/mystudy/myos/operating-system-exercises/ch14/a.out.dSYM/Contents/Resources/DWARF/a.out...
   > ```

3. 对这个程序使用 valgrind 工具。我们将使用属于 valgrind 的 memcheck 工具来分析发生的情况。输入以下命令来运行程序: valgrind --leak-check=yes null 当你运行它时会发生什么?你能解释工具的输出吗?

   > linux...

4. 编写一个使用 malloc()来分配内存的简单程序,但在退出之前忘记释放它。这个程序运行时会发生什么?你可以用 gdb 来查找它的任何问题吗?用 valgrind 呢(再次使用--leak-check=yes 标志)?

5. 编写一个程序,使用 malloc 创建一个名为 data、大小为 100 的整数数组。然后,将 data[100]设置为 0。当你运行这个程序时会发生什么?当你使用 valgrind 运行这个程序时会发生什么?程序是否正确?

6. 创建一个分配整数数组的程序(如上所述),释放它们,然后尝试打印数组中某个元素的值。程序会运行吗?当你使用 valgrind时会发生什么?

   > 答：程序会正常运行。
   >
   > 无效读取，使用 free 时，不会改变被释放变量本身的值，调用 free() 后它仍然会指向相同的内存空间,但是此时该内存已无效free 并不会覆盖释放的内存, 所以读取时仍然能读取到数值。

7. 现在传递一个有趣的值来释放(例如在上面分配的数组中间的一个指针)会发生什么?你是否需要工具来找到这种类型的问题?

8. 尝试一些其他接口来分配内存。例如，创建一个简单的向量似的数据结构，以及使用 realloc()来管理向量的相关函数。使用数组来存储向量元素。当用户在向量中添加条目时，请使用 realloc()为其分配更多空间。这样的向量表现如何？它与链表相比如何？使用 valgrind 来帮助你发现错误。

9. 花更多时间阅读有关使用 gdb 和 valgrind 的信息。了解你的工具至关重要，花时间学习如何成为 UNIX 和 C 环境中的调试器专家。

# 第十五章

1. 用种子 1、2 和 3 运行,并计算进程生成的每个虚拟地址是处于界限内还是界限外? 如果在界限内,请计算地址转换。

   > ```shell
   > ./go/relocation -s 1 -c
   > ./go/relocation -s 2 -c
   > ./go/relocation -s 3 -c
   > ```

2. 使用以下标志运行: `-s 0 -n 10`。为了确保所有生成的虚拟地址都处于边界内,要将(界限寄存器)设置为什么值?

   > ```shell
   > ./go/relocation -s 0 -n 10 -c -l 1k
   > ```

3. 使用以下标志运行:`-s 1 -n 10 -l 100`。可以设置基址的最大值是多少,以便地址空间仍然完全放在物理内存中?

   > ```shell
   > ./go/relocation -s 1 -n 10 -l 100 -c -b 100
   > ```
   >
   > 答：基址寄存器最大值为 16*1024 - 100 = 16284。

4. 运行和第 3 题相同的操作,但使用较大的地址空间 `-a` 和物理内存 `-p`。

5. 作为界限寄存器的值的函数,随机生成的虚拟地址的哪一部分是有效的?画一个图,使用不同随机种子运行,限制值从 0 到最大地址空间大小。

   > ```shell
   > python ./python/base-limit.py 
   > ```

# 第十六章

1. 先让我们用一个小地址空间来转换一些地址。这里有一组简单的参数和几个不同的随机种子。你可以转换这些地址吗?

   ```shell
   ./go/segmentation -a 128 -p 512 -b 0 -l 20 -B 512 -L 20 -s 0
   ./go/segmentation -a 128 -p 512 -b 0 -l 20 -B 512 -L 20 -s 1
   ./go/segmentation -a 128 -p 512 -b 0 -l 20 -B 512 -L 20 -s 2
   ```

2. 现在,让我们看看是否理解了这个构建的小地址空间(使用上面问题的参数)段 0 中最高的合法虚拟地址是什么?段 1 中最低的合法虚拟地址是什么?在整个地址空间中,最低和最高的非法地址是什么?最后,如何运行带有 A 标志的 segmentation.py 来测试你是否正确?

   > ```shell
   > ./go/segmentation -a 128 -p 512 -b 0 -l 20 -B 512 -L 20 -s 1 -A 19,108,20,107 -c
   > ```

3. 假设我们在一个 128 字节的物理内存中有一个很小的 16 字节地址空间。你会设置什么样的基址和界限,以便让模拟器为指定的地址流生成以下转换结果:*valid, valid, violation, ..., violation, valid, valid,即要求 0,1,14,15 有效,其余无效*?假设用以下参数：

   ```shell
   ./go/segmentation -a 16 -p 128 -A 0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15 -b ? -l ? -B ? -L ?
   ```

   > 答
   >
   > ```shell
   > ./go/segmentation -a 16 -p 128 -A 0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15 -b 14 -l 2 -B 1 -L 2 -c
   > ./go/segmentation -a 16 -p 128 -A 0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15 -b 0 -l 2 -B 15 -L 2 -c
   > ```

4. 假设我们想要生成一个问题,其中大约 90%的随机生成的虚拟地址是有效的(即不产生段异常)。你应该如何配置模拟器来做到这一点?哪些参数很重要?

   > 答
   >
   > ```shell
   > ./go/segmentation -a 100 -p 200 -b 0 -l 45 -B 100 -L 45 -n 100 -c -s 10
   > ```

5. 你可以运行模拟器,使所有虚拟地址无效吗?怎么做到?

   > 答
   >
   > ```shell
   > ./go/segmentation -b 0 -l 0 -B 0 -L 0 -n 100 -c
   > ```

# 第十七章

1. 首先运行  `-n 10 -H 0 -p BEST -s 0` 来产生一些随机分配和释放。你能预测 `malloc()` / `free()` 会返回什么吗？你可以在每次请求后猜测空闲列表的状态吗？随着时间的推移，你对空闲列表有什么发现？

   > 答：由于没有合并,空闲空间碎片越来越多。
   >
   > ```shell
   > ./go/malloc -n 10 -H 0 -p BEST -s 0 -c
   > ```

2. 使用最差匹配策略搜索空闲列表 `-p WORST` 时，结果有何不同？什么改变了？

   > 答：空闲空间碎片更多。
   >
   > ```shell
   > ./go/malloc -n 10 -H 0 -p WORST -s 0 -c
   > ```

3. 如果使用首次匹配 `-p FIRST` 会如何？使用首次匹配时，什么变快了？

   > 答：遍历时间变短
   >
   > ```shell
   > ./go/malloc -n 10 -H 0 -p FIRST -s 0 -c
   > ```

4. 对于上述问题，列表在保持有序时，可能会影响某些策略找到空闲位置所需的时间。使用不同的空闲列表排序 `-l ADDRSORT` , `-l SIZESORT+` , `-l SIZESORT-` )查看策略和列表排序如何相互影响。

   > 答：三种排序方式在 free 时会变慢,因为插入空闲块时需要遍历空闲列表,来达成某种排序方式。

   > ```shell
   > ./go/malloc -n 10 -H 0 -p BEST -s 0 -l ADDRSORT -c
   > ./go/malloc -n 10 -H 0 -p WORST -s 0 -l SIZESORT+ -c
   > ./go/malloc -n 10 -H 0 -p WORST -s 0 -l SIZESORT- -c
   > ```

5. 合并空闲列表可能非常重要。增加随机分配的数量【比如说 `-n 1000`】.随着时间的推移，大型分配请求会发生什么？在有和没有合并的情况下运行【即不用和采用 `-C` 标志】。你看到了什么结果差异？每种情况下的空闲列表有多大？在这种情况下，列表的排序是否重要？

   > 答：空闲列表没有合并,随着时间推移,大型分配请求会因为内存空间不足而返回 NULL,且空间碎片越来越多 排序后搜索可能更快。
   >
   > ```shell
   > ./go/malloc -n 1000 -c
   > ./go/malloc -n 1000 -C -c
   > ```

6. 将已分配百分比 `-P` 改为高于 50,会发生什么？它接近 100 时分配会怎样？接近 0 会怎样？

   > 答：高于 50 时,内存空间分配多于释放,可能导致内存不够，接近 0 时, 释放多余分配。
   >
   > ```shell
   > ./go/malloc -P 90 -n 30 -c
   > ./go/malloc -P 10 -n 30 -c
   > ```

7. 要生成高度碎片化的空闲空间，你可以提出怎样的具体请求？使用 `-A` 标志创建碎片化的空闲列表，查看不同的策略和选项如何改变空闲列表的组织。

   > 答：使用最差适应算法申请大量空间大小为 1 的块,然后释放,且不合并即可。
   >
   > ```shell
   > ./go/malloc -A +1,-1,+1,-2,+1,-3,+1,-4,+1,-5,+1,-6,+1,-7 -p WORST -c
   > ```

# 第十八章

1. 在做地址转换之前，让我们用模拟器来研究线性页表在给定不同参数的情况下如何改变大小。在不同参数变化时，计算线性页表的大小。一些建议输入如下，通过使用 `-v` 标志，你可以看到填充了多少个页表项。

   首先，要理解线性页表大小如何随着地址空间的增长而变化：

   ```shell
   ./go/paging-linear-translate -P 1k -a 1m -p 512m -v -n 0
   ./go/paging-linear-translate -P 1k -a 2m -p 512m -v -n 0
   ./go/paging-linear-translate -P 1k -a 4m -p 512m -v -n 0
   ```

   然后，理解线性页表大小如何随页面大小的增长而变化：

   ```shell
   ./go/paging-linear-translate -P 1k -a 1m -p 512m -v -n 0
   ./go/paging-linear-translate -P 2k -a 1m -p 512m -v -n 0
   ./go/paging-linear-translate -P 4k -a 1m -p 512m -v -n 0
   ```

   > 答：页表大小 = 地址空间 / 页面大小。

2. 现在让我们做一些地址转换。从一些小例子开始，使用－u 标志更改分配给地址空间的页数。例如：

   ```shell
   ./go/paging-linear-translate -P 1k -a 16k -p 32k -v -u 0
   ./go/paging-linear-translate -P 1k -a 16k -p 32k -v -u 25
   ./go/paging-linear-translate -P 1k -a 16k -p 32k -v -u 50
   ./go/paging-linear-translate -P 1k -a 16k -p 32k -v -u 75
   ./go/paging-linear-translate -p 1k -a 16k -p 32k -v -u 100
   ```

   如果增加每个地址空间中的页的百分比，会发生什么？

   > 答：有效页比例增大。

3. 现在让我们尝试一些不同的随机种子，以及一些不同的（有时相当疯狂的）地址空间参数：

   ```shell
   ./go/paging-linear-translate -P 8  -a  32    -p 1024 -v -s 1
   ./go/paging-linear-translate -P 8k -a  32k   -p 1m   -v -s 2
   ./go/paging-linear-translate -P 1m -a  256m  -p 512m -v -s 3
   ```

   哪些参数组合是不现实的？为什么？

   > 答：第三个页太大(Linux 页的大小为 4k),导致太多空间被浪费。

4. 利用该程序尝试其他一些问题。你能找到让程序无法工作的限制吗？例如，如果地址空间大小大于物理内存，会发生什么情况？

   ```shell
   ./go/paging-linear-translate -P 0
   ./go/paging-linear-translate -P 1m -a  512m  -p 256m -v -s 3
   ```


# 第十九章

1. 为了计时，可能需要一个计时器，（例如 `gettimeofday()` ）。这种计时器的精度如何？操作要花多少时间，才能让你对它精确计时？（这有助于确定需要循环多少次,反复访问内存页,才能对它成功计时。）

   > `gettimeofday` 精度为微秒级,成本也为微秒级。

2. 写一个程序，命名为 tlb.c，大体测算一下每个页的平均访问时间。程序的输入参数有:页的数目和尝试的次数。

3. 用你喜欢的脚本语言（ `csh`、 `Python` 等）写一段脚本来运行这个程序，当访问页面从 1 增长到几千，也许每次迭代都乘 2。在不同的机器上运行这段脚本，同时收集相应数据。需要试多少次才能获得可信的测量结果？

4. 接下来,将结果绘图,类似于上图。可以用 ploticus 这样的好工具画图。可视化使数据更容易理解,你认为是什么原因？

5. 要注意编译器优化带来的影响。编译器做各种聪明的事情，包括优化掉循环，如果循环中增加的变量后续没有使用。如何确保编译器不优化掉你写的 TLB 大小测算程序的主循环？

   > gcc 选项-O 启用不同级别的优化。使用-O0(默认)禁用它们。-O3 是最高级别的优化。

6. 还有一个需要注意的地方，今天的计算机系统大多有多个 CPU，每个 CPU 当然有自己的 TLB 结构。为了得到准确的测量数据我们需要只在一个 CPU 上运行程序，避免调度器把进程从一个 CPU 调度到另一个去运行。如何做到？【提示:在 Google 上搜索「pinningthread」相关的信息】如果没有这样做，代码从一个 CPU 移到了另一个，会发生什么情况？

   > 切换 cpu 时成本增加,但是使用新的小 TLB 使速度变快。

7. 另一个可能发生的问题与初始化有关。如果在访问数组 a 之前没有初始化，第一次访问将非常耗时，由于初始访问开销，比如要求置 0。这会影响你的代码及其计时吗？如何抵消这些潜在的开销？

   > 不会，没有初始化，且计时器不会记录初始化的时间。

# 第二十六章

1. 开始,我们来看一个简单的程序，「loop.s」。首先,阅读这个程序，看看你是否能理解它： `cat loop.s`。然后，用这些参数运行它：

   ```shell
   ./x86.py -p loop.s -t 1 -i 100 -R dx
   ```

   这指定了一个单线程，每 100 条指令产生一个中断，并且追踪寄存器 `%d`。你能弄清楚 `%dx` 在运行过程中的值吗? 你有答案之后，运行上面的代码并使用 `-c` 标志来检查你的答案。注意答案的左边显示了右侧指令运行后寄存器的值【或内存的值】。

   > ```shell
   > ./go/x86 -p s/loop.s -t 1 -i 100 -R dx -c
   > ```

2. 现在运行相同的代码,但使用这些标志：

   ```shell
   ./go/x86 -p s/loop.s -t 2 -i 100 -a dx=3,dx=3 -R dx -c
   ```

   这指定了两个线程，并将每个 `%dx` 寄存器初始化为 3。 `%dx` 会看到什么值？使用 `-c` 标志运行以查看答案。多个线程的存在是否会影响计算？这段代码有竞态条件吗？

   > 答：线程 1 运行结果与线程 0 相同, 多个线程不会影响计算,因为指令执行长度小于中断周期, 这段代码没有竞态条件。

3. 现在运行以下命令：

   ```shell
   ./go/x86 -p s/loop.s -t 2 -i 3 -r -a dx=3,dx=3 -R dx -c
   ```

   这使得中断间隔非常小且随机，使用不同的种子和-s 来查看不同的交替、中断频率是否会改变这个程序的行为？

   > 答：中断频率不会改变程序的行为, 两个线程没有访问共享变量。

4. 接下来我们将研究一个不同的程序 `s/looping-race-nolock.s` ，该程序访问位于内存地址 2000 的共享变量，简单起见，我们称这个变量为 `x`。使用单线程运行它，并确保你了解它的功能，如下所示：

   ```shell
   ./go/x86 -p s/looping-race-nolock.s -t 1 -M 2000 -c
   ```

   在整个运行过程中， `x`（即内存地址为 2000）的值是多少？使用-c 来检查你的答案。

   > 答：x的值为1。
   >
   > ```shell
   >  2000          Thread 0        
   >     0   
   >     0   1000 mov 2000, %ax  
   >     0   1001 add $1, %ax    
   >     1   1002 mov %ax, 2000  
   >     1   1003 sub  $1, %bx
   >     1   1004 test $0, %bx
   >     1   1005 jgt .top
   >     1   1006 halt
   > ```

5. 现在运行多个迭代和线程：

   ```shell
   ./go/x86 -p s/looping-race-nolock.s -t 2 -a bx=3 -M 2000 -c
   ```

   你明白为什么每个线程中的代码循环 3 次吗？ x 的最终值是什么？

   > 答：x 最终值为【即内存地址 2000 的值)】6，因为运行期间没有被中断，并且有 2 个线程修改内存地址 2000 的值。

6. 现在以随机中断间隔运行：

   ```shell
   ./go/x86 -p s/looping-race-nolock.s -t 2 -M 2000 -i 4 -r -s 0 -c
   ```

   然后改变随机种子，设置 `-s 1`，然后 `-s 2` 等。只看线程交替，你能说出 `x` 的最终值是什么吗？中断的确切位置是否重要？在哪里发生是安全的？中断在哪里会引起麻烦？换句话说，临界区究竟在哪里？

   > 答：`-s` 为1时，x为1；`-s` 为0时，x为2。
   >
   > 临界区是将 2000 区域内存复制到 ax 之后,将 ax 值写回 2000 之前。
   >
   > ```shell
   >  2000          Thread 0                Thread 1        
   >     0   
   >     0   1000 mov 2000, %ax  
   >     0   1001 add $1, %ax    
   >     0   ----- Interrupt ----- ----- Interrupt ----- 
   >     0                           1000 mov 2000, %ax  
   >     0                           1001 add $1, %ax    
   >     1                           1002 mov %ax, 2000  
   >     1                           1003 sub  $1, %bx   
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1   1002 mov %ax, 2000  
   >     1   1003 sub  $1, %bx   
   >     1   1004 test $0, %bx   
   >     1   1005 jgt .top       
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1                           1004 test $0, %bx   
   >     1                           1005 jgt .top       
   >     1                           1006 halt
   >     1   ----- Halt;Switch ----- ----- Halt;Switch ----- 
   >     1   1006 halt
   > ```

7. 现在使用固定的中断间隔来进一步探索程序。运行：

   ```shell
   ./go/x86 -p s/looping-race-nolock.s -a bx=1 -t 2 -M 2000 -i 1 -c
   ```

   看看你能否猜测共享变量 `x` 的最终值是什么。当你改用`-i 2`，`-i 3` 等标志呢？对于哪个中新间隔，程序会给出“正确的”最终答案？

   > 答：-i 为 1 的情况：
   >
   > ```shell
   >  2000          Thread 0                Thread 1        
   >     0   
   >     0   1000 mov 2000, %ax  
   >     0   ----- Interrupt ----- ----- Interrupt ----- 
   >     0                           1000 mov 2000, %ax  
   >     0   ----- Interrupt ----- ----- Interrupt ----- 
   >     0   1001 add $1, %ax    
   >     0   ----- Interrupt ----- ----- Interrupt ----- 
   >     0                           1001 add $1, %ax    
   >     0   ----- Interrupt ----- ----- Interrupt ----- 
   >     1   1002 mov %ax, 2000  
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1                           1002 mov %ax, 2000  
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1   1003 sub  $1, %bx   
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1                           1003 sub  $1, %bx   
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1   1004 test $0, %bx   
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1                           1004 test $0, %bx   
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1   1005 jgt .top       
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1                           1005 jgt .top       
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1   1006 halt
   >     1   ----- Halt;Switch ----- ----- Halt;Switch ----- 
   >     1   ----- Interrupt ----- ----- Interrupt ----- 
   >     1                           1006 halt
   > ```
   >
   > ax 最终值为 1。
   >
   > `-i 2`时,最终值为 1；
   >
   > `-i 3`时最终值为 2；`-i 3`的结果是正确的。

8. 现在为更多循环运行相同的代码（例如 set -a bx=100）。使用-i 标志设置哪些中断间隔会导致“正确”结果？哪些间隔会导致令人惊讶的结果。

   ```shell
   ./go/x86 -p s/looping-race-nolock.s -a bx=100 -t 2 -M 2000 -i 1 -c
   ```

   > 答：结果不确定。

9. 我们来看本作业中最后一个程序【wait-for-me.s】。像这样运行代码：

   ```shell
   ./go/x86 -p s/wait-for-me.s -a ax=1,ax=0 -R ax -M 2000 -c
   ```

   这将线程 0 的 `%ax` 寄存器设置为 1，并将线程 1 的值设置为 0，在整个运行过程中观察 `%ax` 和内存位置 2000 的值。代码的行为应该如何？线程使用的 2000 位置的值如何？它的最终值是什么？

   > 答：ax 寄存器最值没有变化, 地址 2000 最终值为 1。
   >
   > ```shell
   >  2000    ax         Thread 0                Thread 1        
   >     0     1     
   >     0     1     1000 test $1, %ax     
   >     0     1     1001 je .signaller
   >     1     1     1006 mov  $1, 2000
   >     1     1     1007 halt
   >     1     0     ----- Halt;Switch ----- ----- Halt;Switch ----- 
   >     1     0                             1000 test $1, %ax     
   >     1     0                             1001 je .signaller
   >     1     0                             1002 mov  2000, %cx
   >     1     0                             1003 test $1, %cx
   >     1     0                             1004 jne .waiter
   >     1     0                             1005 halt
   > ```

10. 现在改变输入： `./go/x86 -p s/wait-for-me.s -a ax=0,ax=1 -R ax -M 2000` 线程行为如何？线程 0 在做什么？改变中断间隔【例如，`-i 1000`，或者可能使用随机间隔】会如何改变追踪结果？程序是否高效地使用了 CPU？

    > 答：线程 0 一直在循环直到中断【等待 2000 内存的值变为 1】, 没有高效利用 cpu,线程 0 一直占用 cpu 循环等待。

# 第二十八章

1. 首先用标志  `flag.s` 运行 x86.py。该代码通过一个内存标志「实现」锁。你能理解汇编代码试图做什么吗？

   > 答
   >
   > ```
   > flag.s 作用见上面的注释,
   > 这个简单的"锁有一个问题",导致它并不能保证互斥,
   > 比如线程0执行 mov  flag, %ax 完后,时钟终端,切到线程1执行,
   > 而线程1在执行 mov  %ax, count 中断,切到线程0,此时线程1是拥有锁的,
   > 线程0继续执行 test $0, %ax ,这时ax的值是0,因为线程有单独的寄存器!!,所以线程1也获得了锁
   > ```

2. 使用默认值运行时，  `flag.s` 是否按预期工作？它会产生正确的结果吗？使用-M 和-R 标志跟踪变量和寄存器【并使用 `-c` 查看它们的值】。你能预测代码运行时 flag 最终会变成什么值吗?

   > ```shell
   > ./go/x86 -p s/flag.s -c -R ax,bx -M flag,count
   > ```

   > ```shell
   >       flag     count    ax   bx         Thread 0                Thread 1        
   >          0         0     0    0 
   >          0         0     0    0 1000 mov  flag, %ax      
   >          0         0     0    0 1001 test $0, %ax        
   >          0         0     0    0 1002 jne  .acquire       
   >          1         0     0    0 1003 mov  $1, flag       
   >          1         0     0    0 1004 mov  count, %ax     
   >          1         0     1    0 1005 add  $1, %ax        
   >          1         1     1    0 1006 mov  %ax, count     
   >          0         1     1    0 1007 mov  $0, flag       
   >          0         1     1   -1 1008 sub  $1, %bx
   >          0         1     1   -1 1009 test $0, %bx
   >          0         1     1   -1 1010 jgt .top        
   >          0         1     1   -1 1011 halt
   >          0         1     0    0 ----- Halt;Switch ----- ----- Halt;Switch ----- 
   >          0         1     0    0                         1000 mov  flag, %ax      
   >          0         1     0    0                         1001 test $0, %ax        
   >          0         1     0    0                         1002 jne  .acquire       
   >          1         1     0    0                         1003 mov  $1, flag       
   >          1         1     1    0                         1004 mov  count, %ax     
   >          1         1     2    0                         1005 add  $1, %ax        
   >          1         2     2    0                         1006 mov  %ax, count     
   >          0         2     2    0                         1007 mov  $0, flag       
   >          0         2     2   -1                         1008 sub  $1, %bx
   >          0         2     2   -1                         1009 test $0, %bx
   >          0         2     2   -1                         1010 jgt .top        
   >          0         2     2   -1                         1011 halt
   > ```
   >
   > 答：flag的值最终为0。

3. 使用 `-a` 标志更改寄存器 `%bx` 的值【例如,如果只运行两个线程,就用 `-a bx=2,bx=2`】。代码是做什么的？对这段代码问上面的问题，答案如何？

   > 答：count 增加 4 次。
   >
   > ```shell
   > ./go/x86 -p s/flag.s -c -R ax,bx -M flag,count -a bx=2
   > ```

4. 对每个线程将 `bx` 设置为高值,然后使用 `-i` 标志生成不同的中断频率。什么值导致产生不好的结果？什么值导致产生良好的结果？

   > 答：锁无法保证互斥。
   >
   > ```shell
   > ./go/x86 -p s/flag.s -c -R ax,bx -M flag,count -a bx=2 -i 2
   > ```

5. 现在让我们看看程序 `test-and-set.s`。首先尝试理解使用 `xchg` 指令构建简单锁原语的代码。获取锁怎么写？释放锁如何写？

   > 答：当一个线程获取锁之后 mutex 变为 1，释放锁之后 mutex 变为 0，且操作为原子操作,解决的前面的方案带来的问题。
   >
   > ```shell
   > ./go/x86 -p s/test-and-set.s -c -R ax,bx -M mutex,count -a bx=2 -i 4
   > ```

6. 现在运行代码,再次更改中断间隔 `-i` 的值,并确保循环多次。代码是否总能按预期工作？有时会导致 CPU 使用率不高吗？如何量化呢？

   > 答：单核cpu情况下，当一个线程持有锁进入临界区时被抢占,抢占的线程将会自旋一个时间片，导致cpu利用率不高。
   >
   > 量化：计算当一个线程持有锁进入临界区时被抢占，抢占线程的自旋时间长与总时间长百分比即可。
   >
   > ```shell
   > ./go/x86 -p s/test-and-set.s -a bx=2,bx=2 -M count -c -i 1
   > ./go/x86 -p s/test-and-set.s -a bx=2,bx=2 -M count -c -i 2
   > ./go/x86 -p s/test-and-set.s -a bx=2,bx=2 -M count -c -i 3
   > ```

7. 使用 `-P` 标志生成锁相关代码的特定测试。例如，执行一个测试计划，在第一个线程中获取锁，但随后尝试在第二个线程中获取锁。正确的事情发生了吗？你还应该测试什么？

   > ```shell
   > ./go/x86 -p s/test-and-set.s -M mutex,count -R ax,bx -c -a bx=2,bx=2 -P 0011111
   > ```
   >
   > 答：结果争取。

8. 现在让我们看看 peterson.s 中的代码,它实现了 Peterson 算法【在文中的补充栏中提到】研究这些代码,看看你能否理解它。

   > ```shell
   > ./go/x86 -p s/peterson.s -c
   > ```
   >
   > 答：用两个值来确定，两个值全部都设置完成才算加锁。

9. 现在用不同的 -i 值运行代码。你看到了什么样的不同行为？

   > ```shell
   > ./go/x86 -p s/peterson.s -c  -M trun,count,100,101  -R ax,bx,cx,fx -a bx=0,bx=1 -i 2
   > ```

10. 你能控制调度【带 P 标志】来「证明」代码有效吗？你应该展示哪些不同情况？考虑互斥和避免死锁。

11. 现在研究 ticket.s 中 ticket 锁的代码。它是否与本章中的代码相符？

    > ```shell
    > ./go/x86 -p s/ticket.s -c  -M ticket,trun,count
    > ```
    >
    > 答：相符。

12. 现在运行代码，使用以下标志：`-a bx=1000,bx=1000`（标志设置每个线程循环 1000 次）。看看随着时间的推移发生了什么，线程是否花了很多时间自旋等待锁？

    > ```shell
    > ./go/x86 -p s/ticket.s -c  -M ticket,turn,count -R ax,bx,cx -a bx=1000,bx=1000 -M count
    > ```
    >
    > 答：大量时间用于自旋。

13. 添加更多的线程，代码表现如何？

    > ```shell
    > ./go/x86 -p s/ticket.s -c  -M ticket,turn,count -R ax,bx,cx -a bx=1000,bx=1000,bx=1000,bx=1000 -M count -t 4
    > ```
    >
    > 答：线程变多,cpu 利用率下降。

14. 现在来看 yield.s，其中我们假设 yield 指令能够使一个线程将 CPU 的控制权交给另一个线程【实际上，这会是一个 OS 原语，但为了简化仿真，我们假设有一个指令可以完 成任务】。找到一个场景，其中 test-and-set.s 浪费周期旋转，但 yield.s 不会。节省了多少指令？这些节省在什么情况下会出现？

    > ```shell
    > ./go/x86 -p s/test-and-set.s -a bx=2 -i 13 -t 4 | wc -l
    > ./go/x86 -p s/yield.s -a bx=2 -i 13 -t 4 | wc -l
    > ```
    >
    > 答：在当前时间片没拿到锁时，避免了无效自旋。

15. 最后来看 test-and-test-and-set.s。这把锁有什么作用？与 test-and-set.s 相比，它实现 了什么样的优点？

    > 答：减少了写锁 `xchg` 操作, Pentium cpu 上, `xchg` 需要三个时钟周期。 而`mov`只需要一个时钟周期。

# 第三十七章

todo
