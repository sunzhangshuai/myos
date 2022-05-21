{
    sum[$1] += $2
}
END {
    for (category in sum) {
        printf("%s\t%d\n", category, sum[category]) | "sort -t'\t' -k1"
    }
}