# assumes %bx has loop count in it

.main
.top
# critical section
mov 2000, %ax  # 将内存 2000 中的值存到 ax 中
add $1, %ax    # 自增1，$1为立即数1
mov %ax, 2000  # 将内存值写回去

# see if we're still looping
sub  $1, %bx   # bx 减1
test $0, %bx   # bx 和 0
jgt .top       # bx 大于 0 的话，跳转到top

halt