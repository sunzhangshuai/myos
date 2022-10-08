# 2个整数的数组（每个大小为1字节）
# 将标志地址加载到fx寄存器
# 访问 flag[] 0(%fx,%index)
# 其中 %index 是保存0或1的寄存器
# %index：0->flag[0]，1->flag[1]

.var flag 2

# 全局turn变量
.var turn

# 全局计数器
.var count

.main

# 将标志地址放入fx
lea flag, %fx

# 假设 bx 是线程id (0 or 1)
mov %bx, %cx   # bx: self, now copies to cx
neg %cx        # cx: - self
add $1, %cx    # cx: 1 - self

.acquire
mov $1, 0(%fx,%bx)      # flag[self] = 1
mov %cx, turn           # turn       = 1 - self

.spin1
mov 0(%fx,%cx), %ax     # flag[1-self]
test $1, %ax
jne .fini               # if flag[1-self] != 1, skip past loop to .fini

.spin2                  # just labeled for fun, not needed
mov turn, %ax
test %cx, %ax           # compare 'turn' and '1 - self'
je .spin1               # if turn==1-self, go back and start spin again

# fall out of spin
.fini

# do critical section now
mov count, %ax
add $1, %ax
mov %ax, count

.release
mov $0, 0(%fx,%bx)    # flag[self] = 0


# end case: make sure it's other's turn
mov %cx, turn           # turn       = 1 - self
halt