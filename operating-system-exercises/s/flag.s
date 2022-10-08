.var flag
.var count

.main
.top

.acquire
mov  flag, %ax      #
test $0, %ax        #
jne  .acquire       # 如果flag !=0,则跳转到.acquire处,反复检测flag是否为0
mov  $1, flag       # 获取锁(将flag设为1)

# critical section
mov  count, %ax     #
add  $1, %ax        #
mov  %ax, count     # count++

# release lock
mov  $0, flag       # 释放锁(flag设为0)

# see if we're still looping
sub  $1, %bx
test $0, %bx
jgt .top        # 如果bx的值大于0则循环(回到.top处执行)

halt