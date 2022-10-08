.var mutex
.var count

.main
.top

.acquire
mov  $1, %ax
xchg %ax, mutex     # 原子操作:交换ax寄存器与内存mutex空间的值(mutex设为1)
test $0, %ax        #
jne  .acquire       # 如果(%ax)!=0则自旋等待,即原mutex值不为0

# critical section
mov  count, %ax     #
add  $1, %ax        #
mov  %ax, count     # count地址的值+1

# release lock
mov  $0, mutex #  mutex设为0(释放锁)

# see if we're still looping
sub  $1, %bx
test $0, %bx  # 多次循环,直到bx值小于等于0
jgt .top

halt