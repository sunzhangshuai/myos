{
    line = $0
}
/"/ {
    gsub(/"([^"]|\\")*"/, "", line)
} # remove strings
/\// {
    gsub(/\/([^\/]|\\\/)+\//, "", line)
} # reg exprs
/#/ {
    sub(/#.*/, "", line)
} # and comments
{
    n = split(line, x, "[^A-Za-z0-9_]+")  # into words
    for (i = 1; i <= n; i++) {
        if (x[i] ~ /^(close|system|atan2|sin|cos|rand|srand|match|sub|gsub|printf)$/) {
            warn(x[i] " is now a built-in function")
            continue
        }

        if (x[i] ~ /^(ARGC|ARGV|FNR|RSTART|RLENGTH|SUBSEP|FILENAME|BEGIN)$/) {
            warn(x[i] " is now a built-in variable")
            continue
        }

        if (x[i] ~ /^(do|delete|function|return)$/) {
            warn(x[i] " is now a keyword")
            continue
        }

        y[x[i]]++
    }
}

END {
    for (item in y) {
        if (y[item] == 1) {
            printf("%s只出现了一次\n", item)
        }
    }
}

function warn(s) {
    sub(/^[ \t]*/, "")
    printf("file %s, line %d: %s\n\t%s\n", FILENAME, FNR, s, $0)
}
