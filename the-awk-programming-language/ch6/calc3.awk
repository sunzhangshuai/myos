# 中缀表达式
NF > 0 {
    printf("%s = ", $0)
    f=1
    e = expr()
    if (f <= NF) {
        printf("error at %s\n", $f)
    } else {
        printf("%.8g\n", e)
    }
}

# 加减运算
function expr(  e) {
    e = term()
    while ($f == "+" || $f == "-")
        e = $(f++) == "+" ? e + term() : e - term()
    return e
}

# 乘除运算
function term(  e) {        # factor | factor [*/] factor
    e = factor()
    while ($f == "*" || $f == "/" || $f == "%") {
        if ($f == "*") {
            f++
            e = e * factor()
        } else if ($f == "/") {
            f++
            div = factor()
            if (div == 0) {
                printf("无意义\n")
                next
            }
            e = e / div
        } else {
            f++
            div = factor()
            if (div == 0) {
                printf("无意义\n")
                next
            }
            e = e % div
        }
    }
    return e
}

# 解析括号
function factor(  e) {
    if ($f ~ /^[+-]?([0-9]+[.]?[0-9]*|[.][0-9]+)$/) {
        return $(f++)
    } else if ($f == "(") {
        f++
        e = expr()
        if ($(f++) != ")")
            printf("error: missing ) at %s\n", $f)
        return e
    } else {
        printf("error: expected number or ( at %s\n", $f)
        return 0
    }
}