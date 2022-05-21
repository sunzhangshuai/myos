{
    if (!($2 in pcnt)) {
        pcnt[$2] = 0
    }
    if ($1 != 0) {
        pcnt[$1]++
        slist[$2, ++scnt[$2]] = $1
    }
}

END {
    for (node in pcnt) {
        nodecnt++
        if (pcnt[node] == 0) {
            rtsort(node)
        }
    }
    if (pncnt != nodecnt) {
        print "error: input contains a cycle"
    }
    printf("\n")

}

# 深度优先遍历
function rtsort(node,     i, s) {
    visited[node] = 1
    for (i = 1; i <= scnt[node]; i++)
        if (visited[s = slist[node, i]] == 0) {
            rtsort(s)
        } else if (visited[s] == 1) {
            printf("error: nodes %s and %s are in a cycle\n", s, node)
        }
    visited[node] = 2
    printf(" %s", node)
    pncnt++ # count nodes printed
}