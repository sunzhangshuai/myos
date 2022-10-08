.var ticket # ticket = 0
.var turn   # turn = 0
.var count

.main
.top

.acquire
mov $1, %ax
fetchadd %ax, ticket  # ticket+1 后,旧 ticket 值存入ax
.tryagain
mov turn, %cx         # (%cx) = turn
test %cx, %ax
jne .tryagain         # 如果 (%cx) != 旧ticket 值则自旋

# critical section
mov  count, %ax       # get the value at the address
add  $1, %ax          # increment it
mov  %ax, count       # store it back

# release lock
mov $1, %ax
fetchadd %ax, turn     # turn +1,释放锁

# see if we're still looping
sub  $1, %bx
test $0, %bx
jgt .top

halt