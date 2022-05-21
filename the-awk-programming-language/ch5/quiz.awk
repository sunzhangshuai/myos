BEGIN {
    FS = ":"

    if (ARGV[4] == "auto") {
        auto = 1
        ARGC = 4
        ARGV[4] = ""
    }

    if (ARGC != 4) {
        error("usage: awk -f quiz topicfile question answer")
    }
    if (getline < ARGV[1] < 0) {
        error("no such quiz as " ARGV[1])
    }
    q = a = 1
    while ($q !~ ARGV[2] && q <= NF) {
        q++
    }
    while ($a !~ ARGV[3] && a <= NF) {
        a++
    }

    if (q > NF || a > NF || q == a) {
        error("valid subjects are " $0)
    }
    while (getline <ARGV[1] > 0) {
        qa[++nq] = $0
    }

    ARGC = 2
    ARGV[1] = "-"      # now read standard input

    p = random(nq, nq)
    split(p, parr, " ")
    do {
        ++i
        split(qa[parr[i]], x)
        printf("%s? ", x[q])
        if (auto == 1) {
            printf("%s\n", x[a])
            exit
        }
        while ((input = getline) > 0)
            if ($0 ~ "^(" x[a] ")$") {
                print "Right!"
                break
            } else if ($0 == "") {
                print x[a]
                break
            } else {
                printf("wrong, try again: ")
            }
    } while (input > 0 && i < nq)
}

function error(s) {
    printf("error: %s\n", s); exit
}

function random (n, k, i, r, p) {
    srand()
    p = " "
    for (i = n - k + 1; i <= n; i++ ) {
        r = int(rand()*i + 1)
        if (p ~ " " r " ") {
            sub(" " r " ", " " r " " i " ", p)
        } else {
            p = " " r p
        }
    }
    return p
}