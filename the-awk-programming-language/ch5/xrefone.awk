/^\.#/ {
    s[$2] = ++count[$1]; next
    if (saw[$2])
        print NR ": redefinition of", $2, "from line", saw[$2]
    saw[$2] = NR
}
{
    for (i in s) {
        gsub(i, s[i])
    }
    print
}