BEGIN {
    FS = ":"
    n = ARGV[2]
    ARGV[2] = ""
}
{
    x[NR] = $1; y[NR] = $2
}
END   {
    for (i = 1; i <= n; i++)
        print x[randint(NR)], y[randint(NR)]
}

function randint(n) { return int(n * rand()) + 1 }