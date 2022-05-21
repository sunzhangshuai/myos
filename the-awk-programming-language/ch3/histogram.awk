BEGIN {
    bucketNum = 10
    width = 400
}
{
    cap = $1 / bucketNum
    if (cap > bucketCap) {
        bucketCap = cap
    }
    num[NR] = $1
}
END {
    for (i = 1; i <= NR; i++) {
        bucket[int(num[i] / bucketCap + 1)]++
    }
    for (i = 1; i <= bucketNum; i++) {
        printf("%5.1f - %5.1f: %03d %5.2f%% %s\n", (i - 1) * bucketCap, i * bucketCap - 1, bucket[i], bucket[i] * 100 / NR, strLength(bucket[i]))
    }
    printf("%-13.1f: %03d %5.2f%% %s\n", i * bucketCap, bucket[i], bucket[i] * 100 / NR, strLength(bucket[i]))
}
function strLength(v, len, t){
    len = int(v * width / NR)
    while (len-- > 0) {
        t = t "*"
    }
    return t
}