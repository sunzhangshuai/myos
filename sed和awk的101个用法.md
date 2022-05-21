# Sed 基础

> **Strem Editor**：流编辑器

## 执行流程

> REPR：Read、Execute、Print、Repeat。

- **读取**：读取一行到**模式空间**【sed内部的一个临时缓存，用于存放读取到的内容】。
- **执行**：在模式空间中执行命令。如果用了 `{}` 或 `-e` 指定了多个命令，sed 将依次执行每个命令。
- **打印**：打印模式空间的内容，然后清空模式空间。
- **重复**：重复上述过程，直到文件结束。

## 文件执行

> 如果需要重复使用一组 sed 命令，可以建立 sed 脚本文件，里面包含所有要执行的 sed 命令，然后用-f 选项来使用。

- **注释**：以#开头

- **作为解释器**

  ```shell
  # 在脚本最开始添加字符串
  #!/bin/sed -f
  
  # 屏蔽默认输出，只能是nf、不可以是fn
  #!/bin/sed -nf
  ```

## 正则表达式

| 符号  |                     含义                     |
| :---: | :------------------------------------------: |
|   ^   |                   行的开头                   |
|   $   |                   行的结尾                   |
|   .   |                   单个字符                   |
|   +   |                匹配1次或多次                 |
|   ?   |                 匹配0次或1次                 |
|   \   |                   转义字符                   |
| [0-9] |          字符集，匹配方括号中的字符          |
|  \|   |         用来匹配两遍任意一个子表达式         |
|  {m}  |                   匹配m次                    |
| {m,n} |     匹配m~n次，m和n不能为负数，且 < 255      |
|  \b   | 字符边界，匹配单词开头【\bxx】或结尾【xx\b】 |
|  \n   |             回溯引用，分组引用。             |

- **用法举例**

  ```shell
  # 只匹配重复 the 两次的行
  sed -n '/\(the\)\1/ p' file/word.txt
  
  # 把 employee.txt 中每行最后两个字符替换为",Not Defined"
  sed 's/..$/,Not Defined/' file/employee.txt
  
  # 删除以 Manager 开头的行的后面的所有内容
  sed 's/^\(Manager\).*/\1/' file/employee.txt
  
  # 删除所有以#开头的行
  sed '/^#/ d' file/employee.txt
  ```

- **特殊用法**

  > sed 可以把 DOS 的换行符(CR/LF)替换为 Unix 格式。
  >
  > - 当把 DOS 格式的文件拷到 Unix 上， 你会发现，每行结尾都有 `\r\n` 。可以用sed替换。

  ```shell
  # 使用 sed 把 DOS 格式的文件转换为 Unix 格式
  sed 's/.$//' input-file
  ```

# Sed 普通命令

## 语法

1. **通常用法**

   ```shell
   sed [options] 'commands' input-file...
   ```

   > - 每次从 input-file 读取一行命令，执行所有 sed-commands ，在读取第二行，重复。
   >
   > - 可以同时处理多个文件。

2. **命令放在文件中**

   ```shell
   sed [options] -f commands-in-a-file input-file...
   ```

3. **执行多个命令**

   ```shell
   sed [options] -e 'command1' -e 'command2' -e 'command3' input-file...
   ```

   - *换行执行*

     ```shell
     sed -n \
     -e 'command1' \
     -e 'command1' \
     input-file...
     ```

4. **命令分组**

   ```shell
   sed [options] '{
   command1
   command2
   }' input-file...
   ```

## option 命令选项

- `-n`：**屏蔽 sed 默认输出**。

- `-f`：**把 sed 命令保存在脚本文件中**，然后使用 `-f` 选项来调用。

- `-e`：**执行多个命令**，和 `--expression` 等价。

- `-i`：**修改输入文件**，和 `--in-place` 等价。

  ```shell
  # 加bak会在替换前先备份
  sed -ibak 's/John/Johnny/' file/employee.txt
  sed --in-place=bak 's/John/Johnny/' file/employee.txt
  ```

- `-c`：和 `-i` 配合使用，**可以保持文件所有者不变**，和 `--copy` 等价。

- `-l`：和 `l` 命令配合使用，**指定行的长度**。和 `--line-length` 等价。

  ```shell
  # 两天命令等价
  sed -n -l 20 'l' file/employee.txt
  sed -n 'l 20' file/employee.txt
  ```

## command 命令标识

- `p`：**打印模式空间**，可控制只输出置顶内容，一般与 `-n` 配合使用。

- `数字`：**指定行地址范围**。

  ```shell
  # 只打印第 2 行
  sed -n '2 p' file/employee.txt
  
  # 打印 1 到 4 行
  sed -n '1,4 p' file/employee.txt
  
  # 打印第 2 行到最后一行
  sed -n '2,$ p' file/employee.txt
  ```

  - `+`：**指定后几行**，配合 `,` 使用。

    ```shell
    # 打印从第 1 行开始后的 4 行。
    sed -n '1,+4 p' file/employee.txt
    ```

  - `~`：**指定步进**。

    ```shell
    # 从第 1 行开始，每隔 3 行打印一行
    sed -n '1~3 p' file/employee.txt
    ```

- `//`：**模式匹配**

  > 模式匹配后加 `!` 代表不匹配。

  ```shell
  # 打印匹配模式 Jane 的行
  sed -n '/Jane/ p' file/employee.txt
  
  # 打印不匹配模式 Jane 的行
  sed -n '/Jane/! p' file/employee.txt
  
  # 打印第一次匹配 Jason 的行至第 4 行
  sed -n '/Jason/,4 p' file/employee.txt
  
  # 打印第一次匹配 Raj 的行到最后一行
  sed -n '/Raj/,$ p' file/employee.txt
  
  # 打印自匹配 Raj 的行开始到匹配 Jane 的行之间的所有内容
  sed -n '/Raj/,/Jane/ p' file/employee.txt
  
  # 打印匹配 Jason 的行和其后面的两行
  sed -n '/Jason/,+2 p' file/employee.txt
  ```

- `d`：**删除行**，只删除命令空间内容，不修改源文件。

  ```shell
  # 删除第 1 至第 4 行
  sed '1,4 d' file/employee.txt
  
  # 删除空行
  sed '/^$/ d' file/employee.txt
  
  # 删除注释
  sed '/^#/ d' file/employee.txt
  ```

- `w`：**将模式空间的内容写入文件**

  ```shell
  # 将文件内容保存到 output.txt，并输出在屏幕上。
  sed 'w file/output.txt' file/employee.txt
  
  # 将文件内容保存到 output.txt
  sed -n 'w file/output.txt' file/employee.txt
  
  # 保存自匹配 Raj 的行开始到匹配 Jane 的行之间的所有内容
  sed -n '/Raj/,/Jane/ w file/output.txt' file/employee.txt
  ```

- `a`：**追加命令**

  ```shell
  # 在第 2 行后面追加一行
  sed '2 a 203,Jack Johnson,Engineer' file/employee.txt
  
  # 在 employee.txt 文件结尾追加一行
  sed '$ a 203,Jack Johnson,Engineer' file/employee.txt
  
  # 在匹配 Jason 的行的后面追加两行
  sed '/Jason/ a 203,Jack Johnson,Engineer\n204,Mark Smith,Sales Engineer' file/employee.txt
  ```

- `i`：**插入命令**

  ```shell
  # 在第 2 行之前插入一行
  sed '2 i 203,Jack Johnson,Engineer' file/employee.txt
  
  # 在 employee.txt 文件结尾前插入一行
  sed '$ i 203,Jack Johnson,Engineer' file/employee.txt
  
  # 在匹配 Jason 的行前插入两行
  sed '/Jason/ i 203,Jack Johnson,Engineer\n204,Mark Smith,Sales Engineer' file/employee.txt
  ```

- `c`：**修改命令**，可以用新行取代旧行。

  ```shell
  # 用两行新数据取代匹配 Raj 的行
  sed '/Raj/ c 203,Jack Johnson,Engineer\n204,Mark Smith,Sales Engineer' file/employee.txt
  ```

- `l`：**打印不可见字符**。

  > 可以在 l  后指定数字 n，会在 n  个字符处自动折行，`\`表示折行，也占字符。

  ```shell
  sed -n 'l' file/tabfile.txt
  
  sed -n 'l 20' file/tabfile.txt
  ```

- `=`：**打印行号**，在每行命令空间后一行显示该行的行号。

  ```shell
  # 打印1-3行的行号
  sed '1,3 =' file/employee.txt
  
  # 打印文件的总行数
  sed -n '$ =' file/employee.txt
  ```

- `y`：**转换命令**，根据对应位置转换字符。

  ```shell
  # 把所有小写字符转换为大写字符
  sed 'y/abcdefghijklmnopqrstuvwxyz/ABCDEFGHIJKLMNOPQRSTUVWXYZ/' file/employee.txt
  ```

- `q`：**退出sed**，终止正在执行的命令并退出 sed。

  ```shell
  # 打印第 1 行后退出
  sed 'q' file/employee.txt
  
  # 打印第 5 行后退出
  sed '5 q' file/employee.txt
  
  # 遇到包含关键字 Manager 的行后退出
  sed '/Manager/ q' file/employee.txt
  ```

- `r`：**从文件读取命令**。

  ```shell
  # 合并 employee.txt 和 tabfile.txt
  sed '$ r file/tabfile.txt' file/employee.txt
  
  # 将 tabfile.txt 在 Raj 的行后打印出来
  sed '/Raj/ r file/tabfile.txt\n' file/employee.txt
  ```

- `n`：**打印当前模式空间的内容**，并从输入文件中读取下一行。

  > - sed 正常流程是读取数据、执行命令、打印输出、重复循环。
  >
  > - 命令 n 改变这个流程，它打印当前模式空间的内容，然后清除模式空间，读取下一行进 来，然后继续执行后面的命令。

  > ```shell
  > sed-command-1
  > sed-command-2
  > n
  > sed-command-3
  > sed-command-4
  > ```
  >
  > sed-command-1 和 sed-command-2 会在当前模式空间中执行，然后遇到 n，它 打印当前模式空间的内容，并清空模式空间，读取下一行，然后把 sed-command-3 和 sed-command-4 应用于新的模式空间的内容。

# Sed 替换命令

## 语法

```shell
sed '[address-range|pattern-range] s/original-string/replacement-string/[substitute-flags]' {input_file}
```

> - `address-range` 或 `pattern-range`：可选，如果没有指定，在所有行上替换。
> - `s`：**substitute** ，执行替换命令。
> - `original-string`：被 sed 搜索并替换的字符串，可以是正则表达式。
> - `replacement-strin`：替换后的字符串。
> - `substitute-flags`：可选。

## substitute-flags 标识

- `g`：**全局标识**，sed 默认情况下，只替换每行第一次出现的，加 `g` 可以替换所有出现的。

  ```shell
  # 用 A 替换每行第一次出现的 a
  sed 's/a/A/' file/employee.txt
  
  # 用 A 替换所有的 a
  sed 's/a/A/g' file/employee.txt
  ```

- `数字`：**指定每行 `original-string` 出现的次数**，只有第 n 次出现的 `original-string` 才会被替换。

  ```shell
  # 把每行第二次出现的 a 替换为 A
  sed 's/a/A/2' file/employee.txt
  
  # 每行中第二次出现的 locate 替换为 find
  sed 's/locate/find/2' file/substitute-locate.txt
  ```

- `p`：【print】**打印标志**，当替换操作完成后，打印替换后的行。和 `-n` 配套使用。

- `w`：【write】**写标识**，替换操作执行成功后，把替换后的结果保存的文件中。

- `i`：【ignore】**忽略大小写标识**，以小写字符的模式匹配 `original-string`。

  ```shell
  # 把 john 或 John 替换为 Johnny
  sed 's/john/Johnny/i' file/employee.txt
  ```

- `e`：【excuate】**执行命令标识**，将模式空间中的任何内容当做 shell 命令执行，并把命令执行的结果返回到模式空间。

  ```shell
  # 在 files.txt 文件中的每行前面添加 ls -l 并打印结果
  sed 's/^/ls -l /' file/files.txt
  
  # 在 files.txt 文件中的每行前面添加 ls -l 并当成命令执行
  sed 's/^/ls -l /e' file/files.txt
  ```

## replacement-strin 替换标识

>  **和分组配合使用时，这些选项很有用**。

```shell
# 雇员 名称 都显示为大写，职位都显示为小写
sed 's/\([^,]*\),\([^,]*\),\([^,]*\)/\U\2\E,\1,\L\3/' file/employee.txt
```

- `\l`：在 `replacement-string` 中，把紧跟在其后面的字符当做小写字符来处理。

  ```shell
  # 将 John 替换为 JOhNNY
  sed -n 's/John/JO\lHNNY/p' file/employee.txt
  ```

- `\L`：在 `replacement-string` 中，把后面所有字符当做小写字符来处理。

  ```shell
  # 将 John 替换为 JOhnny
  sed -n 's/John/JO\LHNNY/p' file/employee.txt
  ```

- `\u`：在 `replacement-string` 中，把紧跟在其后面的字符当做大写字符来处理。

  ```shell
  # 将 John 替换为 joHnny
  sed -n 's/John/jo\uhnny/p' file/employee.txt
  ```

- `\U`：在 `replacement-string` 中，把后面所有字符当做大写字符来处理。

  ```shell
  # 将 John 替换为 joHNNY
  sed -n 's/John/jo\Uhnny/p' file/employee.txt
  ```

- `\E`：和 `\L` 或 `\U` 配合使用，关闭相关功能。

  ```shell
  # 将 John 替换为 JOHNNY Boy
  sed -n 's/John/\UJohnny\E Boy/p' file/employee.txt
  ```

## 替换命令分界

> sed 默认分界符是 `/`，如果 `original-string` 或 `replacement-strin` 中有 `/`，需要用 `\` 来转义。
>
> 如果替换的是路径的话，需要写的转义符会很多，如：`sed 's/\/usr\/local\/bin/\/usr\/bin/' file/path.txt`。
>
> **可以使用任何一个字符**作为 sed 替换命令的**分界符**，如 `|` 或 `^` 或 `@` 或者 `!`。

```shell
sed 's|/usr/local/bin|/usr/bin|' file/path.txt
sed 's^/usr/local/bin^/usr/bin^' file/path.txt
sed 's@/usr/local/bin@/usr/bin@' file/path.txt
sed 's!/usr/local/bin!/usr/bin!' file/path.txt
```

## &：获取匹配到的模式

> 当在 `replacement-string` 中使用&时，它会被替换成匹配到的 `original-string` 或正则表达式。

```shell
# 给雇员 ID(即第一列的 3 个数字)加上[ ],如 101 改成[101]
sed 's/[0-9]*/[&]/' file/employee.txt

# 把每一行放进< >中
sed 's/^.*/<&>/' file/employee.txt 
```

## 分组替换

### 单个分组

> 以\(开始，以\)结束，可以用在回溯引用中。

- 正则表达式`\([^,]*\)`匹配字符串从开头到第一个逗号之间的所有字符，将其放到第一个分组中。
- `replacement-string` 中的 `\1` 将替代匹配到的分组。

```shell
# 只输出员工的工号
sed 's/\([^,]*\).*/\1/g' file/employee.txt
```

### 多个分组

> 使用多个分组时，需要在 `replacement-string` 中使用 `\n` 来指定第 n 个分组。
>
> sed 最多能处理 **9** 个分组

```shell
# 只打印第一列(雇员 ID)和第三列(雇员职位)
sed 's/\([^,]*\),\([^,]*\),\([^,]*\)/\1,\3/' file/employee.txt

# 交换第一列和第二列
sed 's/\([^,]*\),\([^,]*\),\([^,]*\)/\2,\1,\3/' file/employee.txt

# 格式化数字，增强可读性
sed 's/\(^\|[^0-9.]\)\([0-9]\+\)\([0-9]\{3\}\)/\1\2,\3/g' file/numbers.txt 
```

# Sed 保持空间和模式空间命令

> - **模式空间**：sed 内置的一个缓冲区，用来**存放、修改**从输入文件**读取的内容**。
> - **保持空间**：另外一个缓冲区，用来存放临时数据。每次循环读取数据过程中，模式空间的内容都会被清空，然而保持空间的内容 则保持不变，不会在循环中被删除。

- `x`：【exchange】**交换保持空间和模式空间的内容**。

  ```shell
  # 搜索关键字 Manager 并打印之前的那一行
  sed -n -e '{x;n}' -e '/Manager/ {x;p}' file/empnametitle.txt
  ```

  - `{x;n}`：`x` 交换模式空间和保持空间的内容; `n` 读取下一行到模式空间。
  - `/Manager/ {x;}`：如果当前模式空间的内容包含 Manager，就交换保持空间和模式空间的内容，并打印。

- `h`：【hold】**把模式空间的内容复制到保持空间**。

  ```shell
  # 打印管理者的名字
  sed -n -e '/Manager/! h' -e '/Manager/{x;p}' file/empnametitle.txt
  ```

  - `/Manager/! h`：如果模式空间不包含关键字 Manager，复制模式空间到保持空间。
  - `/Manager/{x;p}`：如果模式空间包含关键字 Manager，就交换保持空间和模式空间的内容，并打印。

- `H`：**把模式空间的内容追加到保持空间**，追加之前保持空间的内容不会被覆盖。

  ```shell
  # 打印管理者的名称和职位
  sed -n -e '/Manager/! h' -e '/Manager/ {H;x;p}' file/empnametitle.txt 
  ```

- `g`：【get】**把保持空间内容复制到模式空间**。

  ```shell
  # 打印管理者的名字
  sed -n -e '/Manager/! h' -e '/Manager/ {g;p}' file/empnametitle.txt
  ```

- `G`：**把保持空间内容追加到模式空间**。

  ```shell
  # 以分号分隔，打印管理者的名称和职位
  sed -n -e '/Manager/! h' -e '/Manager/ {x;G;s/\n/;/p}' file/empnametitle.txt
  ```

# Sed 多行模式及循环命令

- `N`：**从输入文件中读取下一行并追加到模式空间**，而不是替换模式空间。

  ```shell
  # 以分号分隔，打印雇员名称和职位
  sed -n 'N;s/\n/;/p' file/empnametitle.txt 
  
  # 打印文件并在同一行显示行号
  sed -n '=;p' file/employee.txt | sed -n 'N;s/\n/ /p'
  ```

- `P`：**打印多行模式中的第一行**。和 `p` 的区别是遇到 `\n` 就不再打印。

  ```shell
  # 打印所有管理者的名称
  sed -n -e 'N' -e '/Manager/P' file/empnametitle.txt
  ```

- `D`：**删除多行模式中的第一行**。和 `d` 的区别是遇到 `\n` 就不再删除。

  ```shell
  # 去掉文件里的注释
  sed -e '/@/ {N;/@.*@/ {s/@.*@//;P;D}}' file/empnametitle-with-commnet.txt
  ```

- `b`、`:label`标签：**循环和分支**

  > - `:label`：定义一个标签
  > - `b label`：执行改标签后面的命令。
  > - `b`：跳到 sed 脚本的结尾。

  ```shell
  # 把 empnametitle.txt 文件中的雇员名称和职位合并到一行内，字段之间以分号分隔，并且在管理者的名称前面加上一个星号。
  h;n;H;x
  s/\n/;/
  /Manager/!b end
  s/^/*/
  :end
  p
  ```

- `t`：**循环**，如果**前面的命令**执行成功，那么就跳转到 t 指定的标签处，继续往下执行后续命令。否则，仍然继续正常的执行流程。

  ```shell
  # 把 empnametitle.txt 文件中的雇员名称和职位合并到一行内，字段之间以分号分隔，并且在管理者的名称前面加上三个星号。
  h;n;H;x
  s/\n/;/
  :repeat
  /Manager/ s/^/*/ # 前面的命令
  /\*\*\*/! t repeat
  p
  ```

# Sed 模拟 Unix 命令

## cat

```shell
sed '' file/employee.txt
sed -n 'p' file/employee.txt
sed 'n' file/employee.txt
sed 'N' file/employee.txt
```

## grep

```shell
sed -n 's/Jane/&/ p' file/employee.txt
sed -n '/Jane/ p' file/employee.txt
```

## head

```shell
sed '11,$ d' /etc/passwd
sed -n '1,10 p' /etc/passwd
sed '10 q' /etc/passwd
```

# Awk 基础

> 是一个维护和处理文本数据文件的强大语言。
>
> 在文本数据有一定的格式，即每行数据包 含多个以分界符分隔的字段时，显得尤其有用。

## 内置变量

- `FS`：**输入字段分隔符**。默认是空格。

  ```shell
  # 指定多个分隔符
  awk 'BEGIN {FS="[,:%]"} {print $2, $3}' file/employee-multiple-fs.txt
  ```

- `OFS`：**输出字段分隔符**。默认是空格。

  ```shell
  awk 'BEGIN {FS="[,:%]"; OFS=","} {print $2, $3}' file/employee-multiple-fs.txt
  ```

- `RS`：**记录分隔符**。默认是换行符。

- `ORS`：**输出记录分隔符**。默认是换行符。

- `NR`：**记录序号**。

- `FILENAME`：**当前处理的文件**。标准输出是 `-`。

- `FNR`：**文件中的 `NR`**。

## 变量

- **命名**：以字母开头，后续字符可以是数字、字母、或下划线。不可以使用关键字。
- **声明**：不需要事先声明。

### 关联数组

## 操作符

- **一元操作符**：+、-、++、--。
- **算数操作符**：+、-、*、/、%。
- **字符串操作符**：空格。
- **赋值操作符**：=、+=、-=、*=、/=、%=。
- **比较操作符**：>、>=、<、<=、==、!=、&&、||。
- **正则表达式操作符**：~、!~。

### 位操作

- **按位与**：`and(number1, number2)`
- **按位或**：`or(number1, number2)`。
- **按位异或**：`xor(number1, number2)`。
- **左移**：`lshift(number, n)`
- **右移**：`rshift(number, n)`

# Awk 命令

```shell
awk –Fs '/pattern/ {action}' input-file...
awk –Fs '{action}' input-file...
```

> - `-F`：**字段分界符**，不指定默认使用空格。
> - `/pattern/`：可选。**指定模式**，awk 只处理模式匹配的记录，默认处理全部记录。
> - `{action}`：可选。**命令**，所有命令需要放在 `{}` 里。默认打印模式匹配的记录。
> - `/pattern/` 和  `{action}` 需要用 `''` 引起来。

# Awk 程序结构

- **BEGIN 区域**

  ```shell
  BEGIN {awk-commands}
  ```

  - **执行时机**：在最开始【执行 body 区域命令之前】。**执行一次**。
  - **内容**：可选，可以有一个或多个 awk 命令。
  - **BEGIN关键字大写**。

- **body 区域**

  ```shell
  /pattern/ {action}
  ```

  - **执行时机**：每次从文件中读取一行就会执行一次。
  - **关键字**：没有关键字，只有 **表达式** 和 **命令**。

- **END 区域**

  ```shell
  END {awk-commands}
  ```

  - **执行时机**：执行完所有操作后执行，只执行一次。
  - **内容**：可选，可以有一个或多个 awk 命令。
  - **END关键字大写**。

## 分支和循环

- **if**

  ```shell
  if (conditional-expression) {
  	action
  }
  ```

- **if else**

  ```shell
  if (conditional-expression) {
  	action1
  } else {
  	action2
  }
  
  conditional-expression ? action1 : action2
  ```

- **while**

  ```shell
  while (condition) {
  	actions
  }
  ```

- **do while**

  ```shell
  do {
  	actions
  } while (condition)
  ```

- **for**

  ```shell
  for(initialization;condition;increment/decrement) {
  
  }
  ```

- **break**：**只有在循环中才能使用**，用来跳出最内层循环。

- **continue**：**只有在循环中才能使用**，立即进入下次循环。

- **exit**：**立即停止运行的脚本**，可接受状态码，默认是0。

