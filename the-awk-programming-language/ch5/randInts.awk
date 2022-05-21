BEGIN {
    n = ARGV[1]
    k = ARGV[2]
    ARGV[1] = ARGV[2] = ""
    if (k > n) {
        print "无法实现"
        exit
    }
    random(n, k)
}

function random (n, k, A, r, i) {
    for (i = n - k + 1; i <= n; i++) {
        r = randInt(i)
        r in A ? A[i] : A[r]
    }
    for (i in A) {
        print i
    }
}

# 返回一个 1 - n的随机整数
function randInt(n) {
    return int(n * rand()) + 1
}