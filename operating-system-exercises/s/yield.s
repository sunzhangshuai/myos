.var mutex
.var count

.main
.top

.acquire
mov  $1, %ax
xchg %ax, mutex     # 原子操作:交换ax寄存器与内存mutex空间的值(将mutex设为1)
test $0, %ax
je .acquire_done    # 如果ax为0则跳转(即 如果原mutex为0)
yield               # 如果ax不为0(即没有获取到锁),yield
j .acquire          # 从yield返回时,跳转到 .acquire 重新执行
.acquire_done

# critical section
mov  count, %ax
add  $1, %ax
mov  %ax, count     # count++

# release lock
mov  $0, mutex     # mutex设为0,释放锁

sub  $1, %bx
test $0, %bx       # 循环到bx为0
jgt .top

halt