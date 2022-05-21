BEGIN {
    # 加号
    symbol = ARGC > 1 ? ARGV[1] : "+"
    # 计算的最大值
    maxnum = ARGC > 2 ? ARGV[2] : 10
    ARGV[1] = "-"
    srand()

    do {
        n1 = randint(maxnum)
        n2 = randint(maxnum)
        val = getValue(symbol, n1, n2)

        if (ARGC > 3) {
            printf("%g %s %g = %g\n", n1, symbol, n2, val)
            exit
        }

        printf("%g %s %g = ? ", n1, symbol, n2)
        while ((input = getline) > 0) {
            if ($0 == val) {
                print "Right!"
                break
            } else if ($0 == "") {
                print n1 + n2
                break
            } else {
                printf("wrong, try again: ")
            }
        }
    } while (input > 0)
}

function randint(n) {
    return int(rand()*n)+1
}

function getValue(symbol, n1, n2) {
    switch (symbol) {
        case "+":
            return n1 + n2
        case "-":
            return n1 - n2
        case "*":
            return n1 * n2
        case "/":
            return int(n1 / n2)
        case "%":
            return n1 % n2
    }
    printf("%s symbol error", symbol)
    exit
}