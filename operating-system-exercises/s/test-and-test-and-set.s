.var mutex
.var count

.main
.top

.acquire
mov  mutex, %ax
test $0, %ax
jne .acquire        # 如果 mutex 不是0则跳转(即未获得锁)

mov  $1, %ax
xchg %ax, mutex     # 原子操作:交换ax寄存器与内存mutex空间的值(将mutex设为1)
test $0, %ax
jne .acquire        # 如果 mutex 不是0则跳转(即未获得锁)

# critical section
mov  count, %ax     #
add  $1, %ax        #
mov  %ax, count     # count++

# release lock
mov  $0, mutex      # mutex 设为1, 释放锁

# see if we're still looping
sub  $1, %bx
test $0, %bx        #
jgt .top

halt