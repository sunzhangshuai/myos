{
    # 收集数据
    print "label ", $1, $2, $3, "-", $NF | "awk -f ch6/graph.awk"
    for (i = 3; i < NF; i++) {
        print $i, test($1, $2, $i) | "awk -f ch6/graph.awk"
    }
    close("awk -f ch6/graph.awk")
}

# 排序测试
function test(sort, data, n) {
    comp = exch = 0
    if (data ~ /rand/) {
        genrand(A, n)
    } else if (data ~ /id/) {
        genid(A, n)
    } else if (data ~ /rev/) {
        genrev(A, n)
    } else if (data ~ /sort/) {
        gensort(A, n)
    } else {
        print "illegal type of data in", $0
    }

    if (sort ~ /q.*sort/){
        qsort(A, 1, n)
    } else if (sort ~ /h.*sort/) {
        hsort(A, n)
    } else if (sort ~ /i.*sort/){
        isort(A, n)
    } else {
        print "illegal type of sort in", $0
    }

    if (!check(A, n)) {
        printf("array is not sorted\n")
    }
#    print sort, data, n, comp, exch, comp+exch
    return comp+exch
}

# 插入排序
function isort(A, n,   i, j, t) {
    for (i = 2; i <= n; i++) {
        for (j = i; j > 1 && fcomp(A, j-1, j); j--) {
            swap(A, j-1, j)
        }
    }
}

# 快速排序
function qsort(A,left,right,   i,last) {
    if (left >= right) {
        return
    }
    # 选一个放前面
    swap(A, left, left + int((right - left + 1) * rand()))
    last = left
    for (i = left + 1; i <= right; i++) {
        if (fcomp(A, left, i)) {
            swap(A, ++last, i)
        }
    }
    swap(A, left, last)
    qsort(A, left, last-1)
    qsort(A, last+1, right)
}

# 堆排序
function hsort(A,n,  i) {
    for (i = int(n / 2); i >= 1; i--){
        heapify(A, i, n)
    }
    for (i = n; i > 1; i--) {
        swap(A, 1, i)
        heapify(A, 1, i-1)
    }
}

# 比较一层堆
function heapify(A,left,right,   p, c) {
    for (p = left; (c = 2 * p) <= right; p = c) {
        if (c < right && fcomp(A, c+1, c)){
            c++
        }
        if (A[p] < A[c]){
            swap(A, c, p)
        }
    }
}

# 比较
function fcomp(A,i,j) {
    comp++
    return A[i] > A[j]
}

# 交换
function swap(A,i,j,   t) {
    exch++
    t = A[i]; A[i] = A[j]; A[j] = t
}


# 检查是否有序
function check(A, n,   i) {
    for (i = 2; i <= n; i++) {
        if (A[i] < A[i - 1]) {
            return 0
        }
    }
    return 1
}

# n个随机数
function genrand(A, n,   i) {
    srand()
    for (i = 1; i <= n; i++) {
        A[i] = int(n*rand())
    }
}

# n个有序数
function gensort(A, n,   i) {
    for (i = 1; i <= n; i++) {
        A[i] = n
    }
}

# n个逆序数
function genrev(A, n,   i) {
    for (i = 1; i <= n; i++) {
        A[i] = n - i + 1
    }
}

# n个相同的数
function genid(A, n,   i) {
    for (i = 1; i <= n; i++) {
        A[i] = 1
    }
}