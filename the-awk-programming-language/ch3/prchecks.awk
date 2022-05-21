# 输入：由多行组成, 每一行都包括支票编号, 金额, 收款人, 字段之间用制表符分开。
# 输出：标准的支票格式
# 8 行高
# 第 2 行与第 3 行是支票编号与日期, 都向右缩进 45 个空格
# 第 4 行是收款人, 占用 45 个字符宽的区域, 紧跟在它后 面的是 3 个空格, 再后面是金额,
# 第 5 行是金额的大写形式, 其他行都是空白
BEGIN {
    FS = "\t"
    dashes = sp45 = sprintf("%45s", " ")
    gsub(/ /, "-", dashes)
    "date" | getline date
    split(date, d, " ")
    date = d[2] " " d[3] ", " d[6]
    initnum()
}
NF != 3 || $2 > 1000000 || $2 < 0 {
    printf("\nline %d illegal:\n%s\n\nVOID\nVOID\n\n\n", NR, $0)
    next
}
{
    printf("\n")
    printf("%s%s\n", sp45, $1)
    printf("%s%s\n", sp45, date)
    amt = sprintf("%.2f", $2)
    printf("Pay to %45.45s   $%s\n", $3 dashes, addcomma(amt))
    printf("the sum of %s\n", numtowords(amt))
    if (int(amt) < 100000) {
        printf("\n")
    }
    printf("\n\n")
}

function numtowords(n, cents, dols) {
    cents = substr(n, length(n)-1, 2)
    dols = substr(n, 1, length(n)-3)
    if (dols == 0) {
        s = "zero dollars and " cents " cents exactly"
    } else {
        if (int(dols) >= 100000) {
            s = intowords(dols) " dollars\nand " cents " cents exactly"
        } else {
            s = intowords(dols) " dollars and " cents " cents exactly"
        }
    }
    sub(/^one dollars/, "one dollar", s)
    gsub(/  +/, " ", s)
    return s
}

function intowords(n) {
    n = int(n)
    if (n >= 1000)
        return intowords(n/1000) " thousand " intowords(n%1000)
    if (n >= 100)
        return intowords(n/100) " hundred " intowords(n%100)
    if (n >= 20)
        return tens[int(n/10)] "-" intowords(n%10)
    return nums[n]
}

function initnum() {
    split("one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen", nums, " ")
    split("ten twenty thirty forty fifty sixty seventy eighty ninety", tens, " ")
}

function addcomma(n) {
    while (n ~ /[0-9][0-9][0-9][0-9]/) {
        sub(/[0-9][0-9][0-9][,.]/, ",&", n)
    }
    return n
}