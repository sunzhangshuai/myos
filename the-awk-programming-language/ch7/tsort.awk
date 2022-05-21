{
    # pcnt 计算每个节点的入度
    if (!($1 in pcnt)) {
        pcnt[$1] = 0
    }
    if ($2 != 0) {
        pcnt[$2]++
        # 计算每个节点的出度 && 边
        slist[$1, ++scnt[$1]] = $2
    }
}

END {
    for (node in pcnt) {
        nodecnt++
        # 将入度为0的节点入队
        if (pcnt[node] == 0) {
            q[++back] = node
        }
    }

    # 遍历每个入度为0的节点
    for (front = 1; front <= back; front++) {
        printf(" %s", node = q[front])
        # 广度优先遍历
        for (i = 1; i <= scnt[node]; i++) {
            if (--pcnt[slist[node, i]] == 0) {
                q[++back] = slist[node, i]
            }
        }
    }

    if (back != nodecnt) {
        print "\nerror: input contains a cycle"
    }
    printf("\n")
}