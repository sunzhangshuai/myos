BEGIN {
    FS = "!!#"
}
{
    for (i = 2; i <= NF; i++)
        printf("%s\n", $i)
    print "\n"
}