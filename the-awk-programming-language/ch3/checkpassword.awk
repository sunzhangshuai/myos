BEGIN {
    FS = ":"
}

$0 ~ /^#/ {
    next
}

# 密码文件由7个字段组成
NF != 7 {
    printf("line %d, does not have 7 fields: %s\n", NR, $0)
}

# 第 1 个字段是用户的登录名, 只能由字母或数字组成。
$1 ~ /[^_0-9a-zA-Z]+/ {
    printf("line %d, nonalphanumeric user id: %s\n", NR, $0)
}

# 第 2 个字段是加密后的登录密码, 如果密码是空的, 那么任何人都可以利用这个用户名来登录系统。
$2 == "" {
    printf("line %d, no password: %s\n", NR, $0)
}

# 第 3 与第 4 个字段是数字。
$3 ~ /[^-0-9]/ {
    printf("line %d, 21312312 nonnumeric user id: %s\n", NR, $0)
}
$4 ~ /[^-0-9]/ {
    printf("line %d, nonnumeric group id: %s\n", NR, $0)
}

# 第 6 个字段以 / 开始。
$6 !~ /^\// {
    printf("line %d, invalid login directory: %s\n", NR, $0)
}