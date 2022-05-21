BEGIN {
    sep = "\t"

    # 字段映射初始化
    initFieldMap(fields, fieldMap)
}

FILENAME == "file/capitals.txt" && NR == 1 {
    # 关联字段和输出字段
    print "countries表连接字段：" leftKey
    print "capitals表连接字段：" rightKey
    print "输出字段：" outputField
}

FILENAME == "file/capitals.txt" {
    split($0, rowArr, "\t")
    pk = getPrimaryKey(FILENAME, rowArr, rightKey)
    capitalKeyCount[pk]++
    setKeyAndValue(capitalData, pk, capitalKeyCount[pk], FILENAME, rowArr)
    filename = FILENAME
}

FILENAME == "file/countries.txt" {
    split($0, rowArr, sep)
    key = getPrimaryKey(FILENAME, rowArr, rightKey)
    if (!(key in capitalKeyCount)) {
        next
    }
    count = split(outputField, outputFieldArr, " ")
    for (i = 1; i <= capitalKeyCount[key]; i++) {
        for (oFidx in outputFieldArr) {
            if ((FILENAME, outputFieldArr[oFidx]) in fieldMap) {
                val = rowArr[fieldMap[FILENAME, outputFieldArr[oFidx]]]
            } else if ((filename, outputFieldArr[oFidx]) in fieldMap) {
                val = capitalData[key, i, outputFieldArr[oFidx]]
            } else {
                val = "null"
            }
            printf("%-15s", val)
        }
        printf("\n")
    }
}

# 字段映射初始化
function initFieldMap(fields, fieldMap) {
    fields["file/capitals.txt"] = "COUNTRY AREA POPULATION CONTINENT"
    fieldMap["file/countries.txt", "COUNTRY"] =  1
    fieldMap["file/countries.txt", "AREA"] =  2
    fieldMap["file/countries.txt", "POPULATION"] =  3
    fieldMap["file/countries.txt", "CONTINENT"] =  4
    fields["file/capitals.txt"] = "COUNTRY CAPITAL"
    fieldMap["file/capitals.txt", "COUNTRY"] = 1
    fieldMap["file/capitals.txt", "CAPITAL"] = 2
}

# 获取主键
function getPrimaryKey(filename, rowArr, primaryKey, idx) {
    idx = fieldMap[filename, primaryKey]
    return rowArr[idx]
}

# 设置key和value
function setKeyAndValue(data, pk, idx, filename, rowArr, fieldArr, field) {
    split(fields[filename], fieldArr, " ")
    for (field in fieldArr) {
        data[pk, idx, fieldArr[field]] = rowArr[fieldMap[filename, fieldArr[field]]]
    }
}