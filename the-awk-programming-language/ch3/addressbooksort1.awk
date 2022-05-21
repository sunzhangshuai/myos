BEGIN { RS = ""; FS = "\n" }
{
    printf("%s!!#", x[split($1, x, " ")])
    if ($0 ~ /!!#/) {
        next
    }
    for (i = 1; i <= NF; i++)
        printf("%s%s", $i, i < NF ? "!!#" : "\n")
}