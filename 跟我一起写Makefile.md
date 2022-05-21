# 概述

- **性质**：一个工程文件的编译规则，描述了整个工程的编译和链接等规则。

- **包含内容**

  > 哪些文件需要编译，哪些文件不需要编译
  >
  > 哪些文件需要先编译，哪些文件需要后编译
  >
  > 哪些文件需要重建等

- **功能**：使我们项目工程的编译变得**自动化**。

- **解决问题**

  - 编译的时候需要链接库的问题。

    > gcc编译的的时候，要去链接很多的第三方库。在编译的时候命令会很长，并且在编译的时候我们可能会涉及到文件链接的顺序问题。
    >
    > make可以把要链接的库文件放在 Makefile 中，制定相应的规则和对应的链接顺序。这样只需要执行 make 命令，工程就会自动编译。

  - 编译大的工程会花费很长的时间。

    > 每次修改源文件后都要去重新编译。
    >
    > Makefile 支持多线程并发操作，并且只会编译我们修改过的文件。

# makefile介绍

```makefile
edit : main.o kbd.o command.o display.o \
        insert.o search.o files.o utils.o
    cc -o edit main.o kbd.o command.o display.o \
        insert.o search.o files.o utils.o

main.o : main.c defs.h
    cc -c main.c
kbd.o : kbd.c defs.h command.h
    cc -c kbd.c
command.o : command.c defs.h command.h
    cc -c command.c
display.o : display.c defs.h buffer.h
    cc -c display.c
insert.o : insert.c defs.h buffer.h
    cc -c insert.c
search.o : search.c defs.h buffer.h
    cc -c search.c
files.o : files.c defs.h buffer.h command.h
    cc -c files.c
utils.o : utils.c defs.h
    cc -c utils.c
clean :
    rm edit main.o kbd.o command.o display.o \
        insert.o search.o files.o utils.o
```

> - liunx下的cc实际指向gcc。
> - `\`：换行符的意思。

## 文件名称

- 文件名可以是：`GNUmakefile` 、`makefile` 、`Makefile`。通常用 `Makefile`。

- `make` 时查找的文件顺序同上。

- **使用别的文件名**

  ```shell
  make [-f | --file] filename
  ```

## 文件组成

- **显式规则**：说明如何生成一个或多个目标文件。由书写者明显指出要生成的文件、文件的依赖文件和生成的命令。

- **隐晦规则**：由于 make 有自动推导功能，所以隐晦的规则可以让我们简略书写 Makefile。

- **变量的定义**：要定义一系列的变量，变量一般都是字符串，当 make 时，所有变量都会被扩展到引用位置。

- **文件指示**：其包括了三个部分。
  - 在一个Makefile中引用另一个Makefile。
  - 根据变量指定 `Makefile` 的有效部分，就像C语言中的预编译#if一样。
  - 定义一个多行的命令。

- **注释**：只有行注释，其注释是用 `#` 字符。

### 显式规则

```makefile
targets : prerequisites
    command
# or    
targets : prerequisites; command
    command
```

> - **targets**：规则的目标，多个用空格分开，可以使用通配符。可以是以下三类
>   - `Object File`：中间文件。
>   - 可执行文件。
>   - 标签【伪目标】。
> - **prerequisites**：依赖文件，要生成 targets 需要的【文件 | 目标】。
>   - 可以是多个，也可以是没有。
> - **command**：make 需要执行的命令【任意的 shell 命令】。
>   - 可以有多条命令，每一条命令占一行。
>   - 如果不与`target:prerequisites` 在一行，需用 `tab` 开头。
> - 如果命令太长，可以使用反斜杠（ `\` ）作为换行符。

### 引用文件

```makefile
include <filename>
```

- **include**：引用的关键字。
  - `include` 前面可以有空字符，但绝不能是 `Tab` 键开始。
  - `include` 后面可以有一个或多个空格。
- **filename**：当前操作系统Shell的文件模式【包含路径和通配符】。
- **`-`**：不理会引入失败的文件。

#### 文件查找

- **指定**绝对路径或相对路径：直接查找。
- **没指定**
  - 当前目录下寻找。
  - 根据 `-I` 或 `--include-dir` 参数指定的路径中查找。
  - 如果目录 `<prefix>/include` 【一般是： `/usr/local/bin` 或 `/usr/include` 】存在，make也会去找。

#### MAKEFILES

- 做一个类似于 `include` 的动作。
- 引入的文件中的目标【**targets**】不会起作用。

### 变量

> 字符串重复使用时，需要用到变量定义。

- e.g.

  ```makefile
  edit : main.o kbd.o command.o display.o \
          insert.o search.o files.o utils.o
      cc -o edit main.o kbd.o command.o display.o \
          insert.o search.o files.o utils.o
  ```

- 定义

  ```shell
  objects = main.o kbd.o command.o display.o \
      insert.o search.o files.o utils.o
  ```

- 使用

  ```shell
  edit : $(objects)
      cc -o edit $(objects)
  clean :
      rm edit $(objects)
  ```

### 自动推导

> 对于每个 『.o』  文件，自动把『.c』文件加在依赖关系中。

## 分组依赖

```makefile
objects = main.o kbd.o command.o display.o \
    insert.o search.o files.o utils.o

edit : $(objects)
    cc -o edit $(objects)

$(objects) : defs.h
kbd.o command.o files.o : command.h
display.o insert.o search.o files.o : buffer.h

.PHONY : clean
clean :
    rm edit $(objects)
```

- **缺点**：文件依赖关系不清晰。

## 中间文件

> 执行 make 命令时，只有修改过的源文件或者是不存在的目标文件会进行重建，而没有改变的文件不用重新编译，在很大程度上节省时间，提高编程效率。

### 清除中间文件

- makefile 中增加清除规则。

  ```makefile
  .PHONY:clean
  clean:
      -rm edit $(objects)
  ```

  > - `.PHONY` ：表示「clean」 是个伪目标文件。
  >
  > - `-`：不管执行中的错误，继续向下执行。
  > - `clean`：从来都是放在文件的最后。

- 执行 `make clean`。

# make运行

1. 读入所有的Makefile。

   > 找不到文件会报错。`make: *** No targets specified and no makefile found.  Stop.`。

2. 读入被include的其它Makefile。

3. 初始化文件中的变量。

4. 推导隐晦规则，并分析所有规则。

5. 为所有的目标文件创建依赖关系链。

6. 根据依赖关系，决定哪些目标要重新生成。

   1. 找文件中的第一个目标文件『target』，在例子中，会找到「edit」文件，把这个文件作为最终的目标文件。

   2. 当 【「edit」文件不存在| 「edit」 所依赖的『.o』文件的 **文件修改时间** 晚于 「edit」 文件】

      > 执行 command 来生成 「edit」 这个文件。

   3. 当 【「edit」 有依赖的 『.o』 文件不存在】

      > 在当前文件中找 『.o』 文件的依赖性，则再根据 『.o』的规则生成 『.o』 文件。【堆栈过程】

   4. 『.c』文件和『.h』文件存在，否则出错了。

      > 生成『.o』  文件，再用『.o』 文件生成make的终极任务，也就是执行文件「edit」了。

7. 执行生成命令。

## 退出码

- 0：表示成功执行。
- 1：运行时出现任何错误。
- 2：使用了make的「-q」选项，并且make使得一些目标不需要更新。

## 指定Makefile

`-f` 、 `--file` 、`--makefile` 都可以。

将所有指定的 「makefile」 **连在一起**传递给 `make` 执行。

## 指定目标

- **默认行为**：make的最终目标是makefile中的第一个目标。

- **任何目标都可以被指定**成终极目标：

  > 除了以 `-` 打头，或是包含了 `=` 的目标，会被命令行识别成参数或变量。

- `MAKECMDGOALS`：变量中存放你所指定的终极目标的列表。没指定为空。

### 目标书写规则

- **all**：所有目标的目标，其功能是**编译所有目标**。
- **clean**：删除所有被 `make` 创建的中间文件。
- **install**：安装已编译好的程序，把目标执行文件拷贝到指定的目标中去。
- **print**：列出改变过的源文件。
- **tar**：把源程序打包备份。也就是一个tar文件。
- **dist**：创建一个压缩文件，一般是把tar文件压成Z文件。或是gz文件。
- **TAGS**：更新所有的目标，以备**完整地重编译**使用。
- **check**、**test**：用来测试 「makefile」 的流程。

## 检查规则

- `-n`**,** `--just-print`**,** `--dry-run`**,** `--recon`

  > 不执行参数，只打印命令。
  >
  > 对于调试 「makefile」 很有用处。

- `-t`**,** `--touch`

  > 把目标文件的时间更新，但不更改目标文件。

- `-q`**,** `--question`

  > 找目标
  >
  > - **目标存在**：什么也会输出。
  > - **目标不存在**：`make: *** No rule to make target 'clean12'.  Stop.`。

- `-W <file>`**,** `--what-if=<file>`**,** `--assume-new=<file>`**,** `--new-file=<file>`

  > 需要指定一个文件。一般是是源文件【或依赖文件】。
  >
  > 根据规则推导来运行依赖于这个文件的命令。

## 所有参数

- `-b`, `-m`

  > 忽略和其它版本make的兼容性。

- `-B`, `--always-make`

  > 所有目标都更新【重编译】。

- `-C <dir>` , `--directory=<dir>`

  > 指定读取makefile的目录。
  >
  > 如果有多个“-C”参数，后面的路径以前面的作为相对路径，并以最后的目录作为被指定目录。
  >
  > `make -C ~hchen/test -C prog` 等价于 `make -C ~hchen/test/prog` 。

- `-debug[=<options>]`

  > 输出make的调试信息。如果没有参数，那就是输出最简单的调试信息。
  >
  > `options` 的取值
  >
  > - a: 也就是all，输出所有的调试信息。
  > - b: 也就是basic，只输出简单的调试信息。即输出不需要重编译的目标。
  > - v: 也就是verbose，在b选项的级别之上。输出的信息包括哪个makefile被解析，不需要被重编译的依赖文件（或是依赖目标）等。
  > - i: 也就是implicit，输出所以的隐含规则。
  > - j: 也就是jobs，输出执行规则中**命令的详细信息**，如命令的PID、返回码等。
  > - m: 也就是makefile，输出make读取makefile，更新makefile，执行makefile的信息。

- `-d`

  > 相当于“–debug=a”。

- `-e`, `--environment-overrides`

  > 指明环境变量的值覆盖makefile中定义的变量的值。

- `-f=<file>`, `--file=<file>`, `--makefile=<file>`

  > 指定需要执行的makefile。

- `-h`, `--help`

  > 显示帮助信息。

- `-i` , `--ignore-errors`

  > 在执行时忽略所有的错误。

- `-I <dir>` , `--include-dir=<dir>`

  > 指定一个被包含makefile的搜索目标。可以使用多个 `-I` 参数来指定多个目录。

- `-j [<jobsnum>]` , `--jobs[=<jobsnum>]`

  > 指同时运行命令的个数。
  >
  > 如果没有这个参数，make运行命令时能运行多少就运行多少。
  >
  > 如果有一个以上的 `-j` 参数，那么仅最后一个 `-j` 才是有效的。【注意这个参数在MS-DOS中是无用的】

- `-k`, `--keep-going`

  > 出错也不停止运行。
  >
  > 如果生成一个目标失败了，那么依赖于其上的目标就不会被执行了。

- `-l <load>` , `--load-average[=<load>]`, `-max-load[=<load>]`

  > 指定make运行命令的负载。

- `-n`, `--just-print`, `--dry-run`, `--recon`

  > 仅输出执行过程中的命令序列，但并不执行。

- `-o <file>`, `--old-file=<file>`, `--assume-old=<file>`

  > 不重新生成的指定的 「file」，即使这个目标的依赖文件新于它。

- `-p`, `--print-data-base`

  > 输出makefile中的所有数据，包括所有的规则和变量。
  >
  > 这个参数会让一个简单的makefile都会输出一堆信息。如果你只是想输出信息而不想执行makefile，你可以使用“make -qp”命令。

- `-q`, `--question`

  > 不运行命令，也不输出。
  >
  > 仅仅是检查所指定的目标是否需要更新。
  >
  > - 0：说明要更新，
  > - 2：说明有错误发生。

- `-r`, `--no-builtin-rules`

  > 禁止make使用任何隐含规则。

- `-R`, `--no-builtin-variabes`

  > 禁止make使用任何作用于变量上的隐含规则。

- `-s`, `--silent`, `--quiet`

  > 在命令运行时不输出命令的输出。

- `-S`, `--no-keep-going`, `--stop`

  > 取消`-k` 选项的作用。
  >
  > 因为有些时候，make的选项是从环境变量 `MAKEFLAGS` 中继承下来的。

- `-t`, `--touch`

  > 相当于UNIX的touch命令，只把目标的修改日期变成最新的。

- `-v`, `--version`

  > 输出make程序的版本、版权等关于make的信息。

- `-w`, `--print-directory`

  > 输出运行makefile之前和之后的信息。
  >
  > 对于跟踪嵌套式调用make时很有用。
  >
  > e.g.
  >
  > ```shell
  > make: Entering directory `/home/hchen/gnu/make'.
  > 
  > make: Leaving directory `/home/hchen/gnu/make'
  > ```

- `--no-print-directory`

  > 禁止 `-w` 选项。

- `-W <file>`, `--what-if=<file>`, `--new-file=<file>`, `--assume-file=<file>`

  > 目标 `file` 需要更新
  >
  > - 如果和 `-n` 选项使用，输出该目标更新时的运行动作。
  > - 如果没有 `-n` ，更改`file`的修改时间为当前时间。

- `--warn-undefined-variables`

  > 只要make发现有未定义的变量，那么就输出警告信息。

# 规则

## 通配符

> 支持三个通配符： `*` ， `?` 和 `~` 。

- `*`：匹配任意个字符。

- `%`：匹配任意个字符。

  ```makefile
  %.o:%.c
      gcc -o $@ $^
  ```

  > 1. "%.o" 把我们需要的所有的 ".o" 文件组合成为一个列表。
  > 2. 从列表中挨个取出每一个文件，"%" 表示取出来文件的文件名【不包含后缀】，找到文件中和 "%"名称相同的 ".c" 文件。
  > 3. 执行下面的命令，直到列表中的文件全部被取出来为止。属于静态模式规则。

- `?`：自动化变量。

- `~`

  - `~/test`：表示当前用户家目录的 test 目录。
  - `~hchen/test`：表示「hchen」用户家目录的 test 目录。

### 引用变量使用

- 不能在引用变量中直接使用。除非使用 `wildcard` 函数。

  ```makefile
  OBJ=$(wildcard *.c)
  ```

## 文件搜寻

> 当make需要去找寻文件的依赖关系时，可以在文件前加上路径，但**最好是把一个路径告诉make，让make在自动去找**。

### VPATH【大写】

> 一种特殊变量

```makefile
VPATH = src:../headers
```

> 指定目录列表【目录用『:』隔开】，按照顺序进行检索。

### vpath【小写】

一个关键字

```makefile
vpath <pattern> <directories>
```

> 为符合模式『pattern』的文件指定搜索目录『directories』。
>
> - **pattern**：需要包含 `%` 字符。 `%` 的意思是匹配零或若干字符。

```makefile
vpath <pattern>
```

> 清除符合模式『pattern』的文件的搜索目录。

```makefile
vpath
```

> 清除所有已被设置好了的文件搜索目录。

## 伪目标

- **默认目标**：可以作为「默认目标」，只要放在第一个。

- **取名**：不能和文件名重名。
- **依赖**：可以作为依赖。
- **文件生成**：不会有文件生成，所以依赖一定会执行。

### 显示指明伪目标

```makefile
.PHONY : clean
```

## 多目标

> 当多个目标同时依赖一个文件时

```makefile
bigoutput littleoutput : text.g
    generate text.g -$(subst output,,$@) > $@
```

> 等价于：
>
> ```makefile
> bigoutput : text.g
>     generate text.g -big > bigoutput
> littleoutput : text.g
>     generate text.g -little > littleoutput
> ```

- `-$(subst output,,$@)`
  - `$`：表示执行一个makefile函数，函数名为`subst`，后面的为参数。
  - `$@`：表示目标的集合，就像一个数组， `$@` 依次取出目标，并执行命令。

## 静态模式

> 加容易地定义多目标规则。

### 结构

```makefile
<targets ...> : <target-pattern> : <prereq-patterns ...>
    <commands>
    ...
```

- `targets`：定义了一系列的目标文件，可以有通配符。是目标的一个集合。
- `target-pattern`：指明了targets的模式，也就是的目标集模式。
- `prereq-patterns` ：目标的依赖模式，对target-pattern形成的模式再进行一次依赖目标的定义。

### 实例

```makefile
objects = foo.o bar.o

all: $(objects)

$(objects): %.o: %.c
    $(CC) -c $(CFLAGS) $< -o $@
```

1. `%.o` 表明要从 `$object` 中获取的所有以 `.o` 结尾的目标。
2. 依赖模式 `%.c` 则取模式 `%.o` 的 `%` ，也就是 `foo bar` ，并为其加下 `.c` 的后缀。

## 自动生成依赖性

### gcc生成依赖关系

1. ```shell
   gcc -M main.c
   ```

   > 输出：**包含标准库的头文件**。
   >
   > ```makefile
   > main.o: main.c defs.h /usr/include/stdio.h /usr/include/features.h \
   >     /usr/include/sys/cdefs.h /usr/include/gnu/stubs.h \
   >     /usr/lib/gcc-lib/i486-suse-linux/2.95.3/include/stddef.h \
   >     /usr/include/bits/types.h /usr/include/bits/pthreadtypes.h \
   >     /usr/include/bits/sched.h /usr/include/libio.h \
   >     /usr/include/_G_config.h /usr/include/wchar.h \
   >     /usr/include/bits/wchar.h /usr/include/gconv.h \
   >     /usr/lib/gcc-lib/i486-suse-linux/2.95.3/include/stdarg.h \
   >     /usr/include/bits/stdio_lim.h
   > ```

2. ```shell
   gcc -M main.c
   ```

   > 输出
   >
   > ```makefile
   > main.o : main.c defs.h
   > ```

### 自动生成

```makefile
%.d: %.c
    @set -e; rm -f $@; \
    $(CC) -M $(CPPFLAGS) $< > $@.$$$$; \
    sed 's,\($*\)\.o[ :]*,\1.o $@ : ,g' < $@.$$$$ > $@; \
    rm -f $@.$$$$
```

# 命令

- make按顺序一条一条的执行命令。
- 每条命令的开头必须以 `Tab` 键开头，除非命令是紧跟在依赖规则后面的分号后的。
- 命令行之间中的空格或是空行会被忽略，但是如果该空格或空行是以Tab键开头的，那么make会认为其是一个空命令。

## 显示命令

- **默认输出命令**：make会把其要执行的命令行在命令执行前输出到屏幕上。
- **不输出当前命令**：命令前加`@`：。
- **只输出命令，不执行**：参数 `-n` 或 `--just-print` 。
- **不输入任何命令**：参数`-s` 或 `--silent` 或 `--quiet` 。

## 命令执行

- **互不影响**：命令之间互不影响。

  ```makefile
  exec:
      cd /home/hchen
      pwd
  ```

  > pwd显示 `makefile` 所在目录。

- **影响方式**：如果要让上一条命令的结果应用在下一条命令时，应该使用分号分隔这两条命令。

  ```makefile
  exec:
      cd /home/hchen; pwd
  ```

  > pwd显示 `/home/hchen`。

## 命令出错

1. **命令出错**：任何命令执行错误，都会终止执行。

2. **忽略命令出错**：命令前加 `-`。

   ```makefile
   clean:
       -rm -f *.o
   ```

3. **忽略所有命令出错**：参数 `-i` 或 `--ignore-errors` 。

4. **命令出错，只停止该目标**：参数 `-k` 或 `--keep-going` 。

5. **`.IGNORE` 目标**：这个规则中的所有命令将会忽略错误。

## 嵌套执行make

> 我们可以在给每个模块写一个 `makefile` ，然后写一个总控 `makefile`。有如下两种写法。

```makefile
subsystem:
    cd subdir && $(MAKE)
```

```makefile
subsystem:
    $(MAKE) -C subdir
```

### 变量传递

> 总控 `Makefile` 的变量可以传递到下级 `Makefile` 中【显示声明】，但不会覆盖下层的 `Makefile` 中所定义的变量。
>
> 除非 `-e` 参数。

- **传递所有变量**

  ```makefile
  export
  ```

- **传递指定变量**

  ```makefile
  export <variable ...>;
  ```

  > 示例
  >
  > ```makefile
  > # 示例1
  > export variable = value
  > 
  > # 示例2
  > variable = value
  > export variable
  > 
  > # 示例3
  > export variable := value
  > 
  > # 示例4
  > variable := value
  > export variable
  > ```

- **禁止传递指定变量**

  ```makefile
  unexport <variable ...>;
  ```

- **`SHELL` 和 `MAKEFLAGS`**

  >  `MAKEFLAGS` 是 `make` 的参数信息。
  >
  > 这两个参数无论是否 **export**，总会传递到下层。

## 定义命令包

如果 `Makefile` 中出现一些相同命令序列，可以为这些相同的命令序列定义一个变量；语法如下：

```makefile
define run-yacc
yacc $(firstword $^)
mv y.tab.c $@
endef
```

> - 以 `define` 开始，以 `endef` 结束。
>
> - `run-yacc` ：命令包的名字。主意不要和其他变量重名。

### 使用

```makefile
foo.c : foo.y
    $(run-yacc)
```

# 变量

## 变量基础

1. **命名**
   - 可以包含字符、数字，下划线。
   - 可以是数字开头。
   - 不包含 `:` 、 `#` 、 `=` 或是空字符【空格、回车】。
   - 大小写敏感。
   - 命名通常全大写。
2. **声明**：声明时需要给予初值。
3. **使用**：需要给在变量名前加上 `$` 符号，最好用小括号 `()` 或大括号 `{}` 把变量括起来。
4. **`$` 符号**：用 `$$` 来表示。
5. **使用场景**：规则中的「目标」、「依赖」、「命令」以及「新的变量」。

## 变量中的变量

> 在定义变量的值时，我们可以使用其它变量来构造变量的值。

- `=`：变量不一定非要是已定义好的值，其也可以使用后面定义的值。

  ```makefile
  foo = $(bar)
  bar = $(ugh)
  ugh = Huh?
  
  all:
      echo $(foo)
  ```

  - *好处*：可以把变量的真实值推到后面来定义。

  - *坏处*：递归定义。

    ```makefile
    CFLAGS = $(CFLAGS) -O
    
    A = $(B)
    B = $(A)
    ```

- `:=`：前面的变量不能使用后面的变量，只能使用前面已定义好了的变量。

- `?=`

  ```makefile
  FOO ?= bar
  ```

  > 如果FOO没有被定义过，那么变量FOO的值就是「bar」；
  >
  > 如果FOO先前被定义过，什么也不做。

## 高级用法

### 变量值替换

- **替换变量中的共有部分**

  ```makefile
  $(var:a=b)
  # or
  ${var:a=b}
  ```

  > 把变量「var」中所有「以a字串结尾」的「a」替换成「b」字串。

  - e.g.

    ```makefile
    foo := a.o b.o c.o
    bar := $(foo:.o=.c)
    ```

    > 把 `$(foo)` 中所有以 `.o` 字串“结尾”全部替换成 `.c` 
    >
    >  `$(bar)` 的值是「a.c b.c c.c」。

- **静态模式**

  > 依赖于被替换字串中的有相同的模式。

  ```makefile
  foo := a.o b.o c.o
  bar := $(foo:%.o=%.c)
  ```

### 变量的值作为变量

```makefile
x = $(y)
y = z
z = Hello
a := $($(x))
```

> `$(a)` = Hello

```makefile
first_second = Hello
a = first
b = second
all = $($a_$b)
```

> `$(all)` = Hello

## 追加变量值

`+=` ：给变量追加值。

1. **变量之前没有定义**： `+=` -> `=` 。
2. **变量之前有定义**：继承前次赋值符。

## override 指示符

对于命令行参数设置的变量，`Makefile`会忽略对这个变量的赋值。**可使用 `override` 强制生效**。

```makefile
override <variable>; = <value>;
override <variable>; := <value>;
override <variable>; += <more text>;
```

- 也可以重写多行变量。

  ```makefile
  override define foo
  bar
  endef
  ```

## 多行变量

- **开始关键字**：**define**。
- **变量名称**：跟在 **define** 后面。
- **变量的值**：重起一行定义变量的值，可以包含函数、命令、文字，或是其它变量。
- **结尾关键字**：**endef**

## 环境变量

1. `makeflie` 的变量值 > 环境变量 。
2. 环境变量能传到下层 `Makeflie`，而文件中的变量需要显示置顶 `export`。

## 目标变量

> **Target-specific Variable**
>
> 可以为某个目标设置局部变量。
>
> **作用范围**：这条规则以及连带规则中。

```makefile
<target ...> : <variable-assignment>;

<target ...> : overide <variable-assignment>
```

- **variable-assignment**：可以是任何赋值表达式。
- **override**：针对于**系统环境变量**，或是make**命令行指定的变量**。

## 模式变量

> **Pattern-specific Variable**
>
> 给定一种“模式”，可以把变量定义在符合这种模式的所有目标上。

```makefile
<pattern ...>; : <variable-assignment>;

<pattern ...>; : override <variable-assignment>;
```

# 条件判断

## 语法

- **if**

  ```makefile
  <conditional-directive>
  <text-if-true>
  endif
  ```

- **if - else**

  ```makefile
  <conditional-directive>
  <text-if-true>
  else
  <text-if-false>
  endif
  ```

## 组成

- **conditional-directive**：条件关键字，有4个。

  - **ifeq**：比较两个参数是否相同，参数可以使用函数。

    ```makefile
    ifeq (<arg1>, <arg2>)
    ifeq '<arg1>' '<arg2>'
    ifeq "<arg1>" "<arg2>"
    ifeq "<arg1>" '<arg2>'
    ifeq '<arg1>' "<arg2>"
    ```

  - **ifneq**：不同为真。

  - **ifdef**：如果变量的值非空，那到表达式为真。

    ```makefile
    ifdef <variable-name>
    ```

    > 示例
    >
    > ```makefile
    > bar =
    > foo = $(bar)
    > ifdef foo
    >     frobozz = yes
    > else
    >     frobozz = no
    > endif
    > ```

  - **ifndef**：和 `ifdef` 相反。

- **conditional-directive、else、endif**

  > 可以有空格，但不能是 `tab` 开始。

# 函数

> 函数调用后，函数的返回值可以当做变量来使用。
>
> 为了风格的统一，函数和变量的括号最好一样。

```makefile
$(<function> <arguments>)
# or
${<function> <arguments>}
```

- **function**：函数名，和参数间以空格分隔。
- **arguments**：函数的参数，参数间以逗号 `,` 分隔。

## 字符串处理函数

### subst

```makefile
$(subst <from>,<to>,<text>)
```

- **名称**：字符串替换函数
- **功能**：把字串 `<text>` 中的 `<from>` 字符串替换成 `<to>` 。
- **返回**：函数返回被替换过后的字符串。

### patsubst

```makefile
$(patsubst <pattern>,<replacement>,<text>)
```

- **名称**：模式字符串替换函数。
- **功能**：查找 `<text>` 中的**单词**【单词以「空格」、「Tab」或「回车」「换行」分隔】是否符合模式 `<pattern>` ，如果匹配的话，则以 `<replacement>` 替换。这里， `<pattern>` 可以包括通配符 `%` ，表示任意长度的字串。如果 `<replacement>` 中也包含 `%` ，那么， `<replacement>` 中的这个 `%` 将是 `<pattern>` 中的那个 `%` 所代表的字串。（可以用 `\` 来转义，以 `\%` 来表示真实含义的 `%` 字符）
- **返回**：函数返回被替换过后的字符串。

### strip

```makefile
$(strip <string>)
```

- **名称**：去空格函数。
- **功能**：去掉 `<string>` 字串中开头和结尾的空字符。
- **返回**：返回被去掉空格的字符串值。

### findstring

```makefile
$(findstring <find>,<in>)
```

- **名称**：查找字符串函数。
- **功能**：在字串 `<in>` 中查找 `<find>` 字串。
- **返回**：如果找到，那么返回 `<find>` ，否则返回空字符串。

### filter

```makefile
$(filter <pattern...>,<text>)
```

- **名称**：过滤函数。
- **功能**：以 `<pattern>` 模式过滤 `<text>` 字符串中的**单词**，保留符合模式 `<pattern>` 的单词。可以有多个模式。
- **返回**：返回符合模式 `<pattern>` 的字串。

### filter-out

```makefile
$(filter-out <pattern...>,<text>)
```

- **名称**：反过滤函数。
- **功能**：以 `<pattern>` 模式过滤 `<text>` 字符串中的单词，去除符合模式 `<pattern>` 的单词。可以有多个模式。
- **返回**：返回不符合模式 `<pattern>` 的字串。

### sort

```shell
$(sort <list>)
```

- **名称**：排序函数。
- **功能**：给字符串 `<list>` 中的**单词**排序【升序】。
- **返回**：返回排序后的字符串。
- **示例**： `$(sort foo bar lose)` 返回 `bar foo lose` 。
- **备注**： `sort` 函数会去掉 `<list>` 中相同的单词。

### word

```makefile
$(word <n>,<text>)
```

- **名称**：取单词函数。
- **功能**：取字符串 `<text>` 中第 `<n>` 个**单词**。【从1开始】
- **返回**：返回字符串 `<text>` 中第 `<n>` 个单词。如果 `<n>` 比 `<text>` 中的单词数要大，那么返回空字符串。

### wordlist

```makefile
$(wordlist <ss>,<e>,<text>)
```

- **名称**：取单词串函数。
- **功能**：从字符串 `<text>` 中取从 `<ss>` 开始到 `<e>` 的单词串。 `<ss>` 和 `<e>` 是一个数字。
- **返回**：返回字符串 `<text>` 中从 `<ss>` 到 `<e>` 的单词字串。如果 `<ss>` 比 `<text>` 中的单词数要大，那么返回空字符串。如果 `<e>` 大于 `<text>` 的单词数，那么返回从 `<ss>` 开始，到 `<text>` 结束的单词串。

### words

```makefile
$(words <text>)
```

- **名称**：单词个数统计函数。
- **功能**：统计 `<text>` 中字符串中的单词个数。
- **返回**：返回 `<text>` 中的单词数。
- **示例**： `$(words, foo bar baz)` 返回值是 `3` 。
- **备注**：如果我们要取 `<text>` 中最后的一个单词，我们可以这样： `$(word $(words <text>),<text>)` 。

### firstword

```makefile
$(firstword <text>)
```

- **名称**：首单词函数——firstword。
- **功能**：取字符串 `<text>` 中的第一个单词。
- **返回**：返回字符串 `<text>` 的第一个单词。

## 文件名操作函数

### dir

```makefile
$(dir <names...>)
```

- **名称**：取目录函数。
- **功能**：从文件名序列 `<names>` 中取出目录部分。指最后一个反斜杠（ `/` ）之前的部分。如果没有反斜杠，那么返回 `./` 。
- **返回**：返回文件名序列 `<names>` 的目录部分。

### notdir

```makefile
$(notdir <names...>)
```

- **名称**：取文件函数。
- **功能**：从文件名序列 `<names>` 中取出非目录部分。非目录部分是指最后一个反斜杠（ `/` ）之后的部分。
- **返回**：返回文件名序列 `<names>` 的非目录部分。

### suffix

```makefile
$(suffix <names...>)
```

- **名称**：取后缀函数。
- **功能**：从文件名序列 `<names>` 中取出各个文件名的后缀。
- **返回**：返回文件名序列 `<names>` 的后缀序列，如果文件没有后缀，则返回空字串。

### basename

```makefile
$(basename <names...>)
```

- **名称**：取前缀函数。
- **功能**：从文件名序列 `<names>` 中取出各个文件名的前缀部分。
- **返回**：返回文件名序列 `<names>` 的前缀序列，如果文件没有前缀，则返回空字串。

### addsuffix

```makefile
$(addsuffix <suffix>,<names...>)
```

- **名称**：加后缀函数。
- **功能**：把后缀 `<suffix>` 加到 `<names>` 中的每个单词后面。
- **返回**：返回加过后缀的文件名序列。

### addprefix

```makefile
$(addprefix <prefix>,<names...>)
```

- **名称**：加前缀函数。
- **功能**：把前缀 `<prefix>` 加到 `<names>` 中的每个单词前面。
- **返回**：返回加过前缀的文件名序列。

### join

```makefile
$(join <list1>,<list2>)
```

- **名称**：连接函数。
- **功能**：把 `<list2>` 中的单词对应地加到 `<list1>` 的单词后面。如果 `<list1>` 的单词个数要比 `<list2>` 的多，那么， `<list1>` 中的多出来的单词将保持原样。如果 `<list2>` 的单词个数要比 `<list1>` 多，那么， `<list2>` 多出来的单词将被复制到 `<list1>` 中。
- **返回**：返回连接过后的字符串。
- **示例**： `$(join aaa bbb , 111 222 333)` 返回值是 `aaa111 bbb222 333` 。

## foreach 函数

```makefile
$(foreach <var>,<list>,<text>)
```

- **名称**：循环处理函数。
- **功能**
  - 把 `<list>` 中的单词逐一取出放到参数 `<var>` 所指定的变量中，然后再执行 `<text>` 所包含的表达式。
  - 每一次 `<text>` 会返回一个字符串，循环过程中， `<text>` 的所返回的每个字符串会以空格分隔，最后当整个循环结束时， `<text>` 所返回的每个字符串所组成的整个字符串（以空格分隔）将会是foreach函数的返回值。
- **返回**：`<text>` 表达式处理过的单词序列。

## if 函数

```makefile
$(if <condition>,<then-part>)
# or
$(if <condition>,<then-part>,<else-part>)
```

- **名称**：判断函数。
- **功能**：`<condition>` 参数是if的表达式，如果其返回的为非空字符串，那么这个表达式就相当于返回真，于是， `<then-part>` 会被计算，否则 `<else-part>` 会被计算。
- **返回值**：如果 `<condition>` 为真（非空字符串），那个 `<then-part>` 会是整个函数的返回值，如果 `<condition>` 为假（空字符串），那么 `<else-part>` 会是整个函数的返回值，此时如果 `<else-part>` 没有被定义，那么，整个函数返回空字串。

## call函数

> 唯一一个可以用来创建新的参数化的函数。

```makefile
$(call <expression>,<parm1>,<parm2>,...,<parmn>)
```

- **名称**：向表达式传递参数。

- **功能**：

  > 可以写一个非常复杂的表达式，这个表达式中，你可以定义许多参数，然后你可以call函数来向这个表达式传递参数。 
  >
  > `<expression>` 参数中的变量，如 `$(1)` 、 `$(2)` 等，会被参数 `<parm1>` 、 `<parm2>` 、 `<parm3>` 依次取代。

- **返回值**：`<expression>` 的返回值。

- **示例**： 

  ```makefile
  reverse =  $(2) $(1)
  foo = $(call reverse,a,b)
  ```

  > 返回值：`b a`。

- **注意**：第2个及其之后的参数中的空格会被保留。

## origin

```makefile
$(origin <variable>)
```

- **名称**：变量来源函数。

- **功能**：`variable`是变量名称，不是引用, 获取变量来源。

  > 可以写一个非常复杂的表达式，这个表达式中，你可以定义许多参数，然后你可以call函数来向这个表达式传递参数。 
  >
  > `<expression>` 参数中的变量，如 `$(1)` 、 `$(2)` 等，会被参数 `<parm1>` 、 `<parm2>` 、 `<parm3>` 依次取代。

- **返回值**

  - `undefined`： `variable` 从来没有定义过。
  - `default`：默认的定义，比如`CC`变量。
  - `environment`：环境变量，并且没指定 `-e` 参数。
  - `file`：变量被定义在Makefile中。
  - `command line`：命令行定义的参数。
  - `override`：被override指示符重新定义。
  - `automatic`：命令运行中的自动化变量。

## shell函数

```makefile
$(shell command line)
```

- **名称**：sh函数。
- **功能**：执行操作系统函数。
- **返回值**：`command line` 的返回值。
- **注意**：和反引号「`」是相同的功能。

## 控制make的函数

### error

```makefile
$(error <text ...>)
```

- **名称**：报错函数。
- **功能**：产生一个致命的错误， `<text ...>` 是错误信息。
- **注意**：error函数不会在一被使用就会产生错误信息，所以如果你把其定义在某个变量中，并在后续的脚本中使用这个变量，那么也是可以的。

### warning

```makefile
$(warning <text ...>)
```

- **名称**：警告函数。
- **功能**：产生一段警告， `<text ...>` 是警告信息。

# 隐含规则

1. **编译C程序的隐含规则**

   > `*.o` 的目标的依赖目标会自动推导为 `*.c` ，并且其生成命令是 `$(CC) –c $(CPPFLAGS) $(CFLAGS)`。

2. **编译C++程序的隐含规则**

   > `*.o` 的目标的依赖目标会自动推导为 `*.cc` 或是 `*.C` ，并且其生成命令是 `$(CXX) –c $(CPPFLAGS) $(CFLAGS) `。

3. **链接Object文件的隐含规则**

   > `n` 目标依赖于 `n.o` ，通过运行C的编译器来运行链接程序生成【一般是 ld 】，
   >
   > 其生成命令是： `$(CC) $(LDFLAGS) n.o $(LOADLIBES) $(LDLIBS) `。

## 取消预设值隐藏规则

 `-r` 或 `--no-builtin-rules` 选项。

## 隐含规则使用的变量

 `-R` 或 `--no–builtin-variables` 参数来取消自定义变量对隐含规则的影响。

### 命令的变量

- `AR` : 函数库打包程序。默认命令是 `ar`
- `AS` : 汇编语言编译程序。默认命令是 `as`
- `CC` : C语言编译程序。默认命令是 `cc`
- `CXX` : C++语言编译程序。默认命令是 `g++`
- `CO` : 从 RCS文件中扩展文件程序。默认命令是 `co`
- `CPP` : C程序的预处理器（输出是标准输出设备）。默认命令是 `$(CC) –E`
- `FC` : Fortran 和 Ratfor 的编译器和预处理程序。默认命令是 `f77`
- `GET` : 从SCCS文件中扩展文件的程序。默认命令是 `get`
- `LEX` : Lex方法分析器程序（针对于C或Ratfor）。默认命令是 `lex`
- `PC` : Pascal语言编译程序。默认命令是 `pc`
- `YACC` : Yacc文法分析器（针对于C程序）。默认命令是 `yacc`
- `YACCR` : Yacc文法分析器（针对于Ratfor程序）。默认命令是 `yacc –r`
- `MAKEINFO` : 转换Texinfo源文件（.texi）到Info文件程序。默认命令是 `makeinfo`
- `TEX` : 从TeX源文件创建TeX DVI文件的程序。默认命令是 `tex`
- `TEXI2DVI` : 从Texinfo源文件创建军TeX DVI 文件的程序。默认命令是 `texi2dvi`
- `WEAVE` : 转换Web到TeX的程序。默认命令是 `weave`
- `CWEAVE` : 转换C Web 到 TeX的程序。默认命令是 `cweave`
- `TANGLE` : 转换Web到Pascal语言的程序。默认命令是 `tangle`
- `CTANGLE` : 转换C Web 到 C。默认命令是 `ctangle`
- `RM` : 删除文件命令。默认命令是 `rm –f`

### 命令参数的变量

- `ARFLAGS` : 函数库打包程序AR命令的参数。默认值是 `rv`
- `ASFLAGS` : 汇编语言编译器参数。【当明显地调用 `.s` 或 `.S` 文件时】
- `CFLAGS` : C语言编译器参数。
- `CXXFLAGS` : C++语言编译器参数。
- `COFLAGS` : RCS命令参数。
- `CPPFLAGS` : C预处理器参数。【C 和 Fortran 编译器也会用到】。
- `FFLAGS` : Fortran语言编译器参数。
- `GFLAGS` : SCCS “get”程序参数。
- `LDFLAGS` : 链接器参数。【如： `ld` 】
- `LFLAGS` : Lex文法分析器参数。
- `PFLAGS` : Pascal语言编译器参数。
- `RFLAGS` : Ratfor 程序的Fortran 编译器参数。
- `YFLAGS` : Yacc文法分析器参数。

## 隐含规则链

一个目标可能被一系列的隐含规则所作用。把这一系列的隐含规则叫做「隐含规则链」。

- **中间目标**：make自动推导出的目标。

  > 和一般目标的区别。
  >
  > 1. 除非中间的目标不存在，才会引发中间规则。
  > 2. 只要目标成功产生，那么，产生最终目标过程中，所产生的中间目标文件会被以 `rm -f` 删除。

- **指定中间目标**

  ```makefile
  .INTERMEDIATE : mid
  ```

- **阻止自动删除中间目标**

  ```makefile
  .SECONDARY : sec
  ```

## 模式规则

```makefile
%.o : %.c
    $(CC) -c $(CFLAGS) $(CPPFLAGS) $< -o $@
```

> - `%` 表示长度任意的非空字符串。e.g. `s.%.c` 则表示以 `s.` 开头， `.c` 结尾的文件名。
> - `$<` 表示所有依赖目标的每个值。
> - `$@` 表示所有的目标的每个值。

### 自动化变量

- **功能**：把模式中所定义的一系列的文件自动地挨个取出，直至所有的符合模式的文件都取完了。

- **使用场景**：只应出现在规则的命令中。

#### 变量列表

- `$@` : 表示规则中的目标文件集。在模式规则中，如果有多个目标，那么， `$@` 就是匹配于目标中模式定义的集合。

- `$%` 

  - 当目标是**函数库文件**，表示规则中的目标成员名。例如，如果一个目标是 `foo.a(bar.o)` ，那么， `$%` 就是 `bar.o` ， `$@` 就是 `foo.a` 。
  - 当目标不是**函数库文件**【Unix下是 `.a` ，Windows下是 `.lib` 】，那么，其值为空。

- `$<` : 依赖目标中的第一个目标名字。

  > 如果依赖目标是以模式【即 `%` 】定义的，那么 `$<` 将是符合模式的一系列的文件集。注意，其是**一个一个取出来**的。

- `$?` : 所有比目标新的依赖目标的集合。以空格分隔。

- `$^` : 所有的依赖目标的集合。以空格分隔。如果在依赖目标中有多个重复的，那么这个变量会**去除重复的依赖目标**，只保留一份。

- `$+` : 这个变量很像 `$^` ，也是所有依赖目标的集合。只是它**不去除重复的依赖目标**。

- `$*` : 这个变量表示目标模式中 `%` 及其之前的部分。

  - 如果目标是 `dir/a.foo.b` ，并且目标的模式是 `a.%.b` ，那么， `$*` 的值就是 `dir/a.foo` 。
  - 这个变量对于构造有关联的文件名是比较有用。
  - 如果目标中没有模式的定义，那么 `$*` 也就不能被推导出，但是，如果目标文件的后缀是make所识别的，那么 `$*` 就是除了后缀的那一部分。
  - 例如：如果目标是 `foo.c` ，因为 `.c` 是make所能识别的后缀名，所以， `$*` 的值就是 `foo` 。这个特性是GNU make的，很有可能不兼容于其它版本的make，所以，你应该**尽量避免使用 `$*`** ，除非是在隐含规则或是静态模式中。如果目标中的后缀是make所不能识别的，那么 `$*` 就是空值。

#### 变量扩展

- `$(@D)`：表示 `$@` 的目录部分【不以斜杠作为结尾】。

  > 如果 `$@` 值是 `dir/foo.o` ，那么 `$(@D)` 就是 `dir` ，而如果 `$@` 中没有包含斜杠的话，其值就是 `.` （当前目录）。

- `$(@F)`：表示 `$@` 的文件部分。相当于函数 `$(notdir $@)` 。
- `$(*D)`, `$(*F)`
- `$(<D)`, `$(<F)`：分别表示依赖文件的目录部分和文件部分。
- `$(^D)`, `$(^F)`：分别表示所有依赖文件的目录部分和文件部分。【去重】
- `$(+D)`, `$(+F)`：分别表示所有依赖文件的目录部分和文件部分。【不去重】
- `$(?D)`, `$(?F)`：分别表示被更新的依赖文件的目录部分和文件部分。

### 模式的匹配

- **茎**： `%` 所匹配的内容。

- **包含有斜杠**：目录部分会首先被移开，然后进行匹配，成功后，再把目录加回去。

  > e.g.：模式 `e%t` ，文件 `src/eat` 匹配于该模式，于是 `src/a` 就是其「**茎**」。

### 重载隐含规则

- **重载**

  ```makefile
  %.o : %.c
      $(CC) -c $(CPPFLAGS) $(CFLAGS) -D$(date)
  ```

- **取消**

  ```makefile
  %.o : %.s
  ```

  > 不写命令就行。

## 后缀规则

> - 后缀规则中所定义的后缀应该是make所认识的。
> - 后缀规则不允许任何的依赖文件。
> - 后缀规则没有命令，毫无意义。他不会**移去内建的隐含规则**。

- **双后缀**： `.c.o` 相当于 `%o : %c` 

  ```makefile
  .c.o:
      $(CC) -c $(CFLAGS) $(CPPFLAGS) -o $@ $<
  ```

- **单后缀**： `.c` 相当于 `% : %.c` 。

- **.SUFFIXES**：定义或删除make知道的后缀。

  ```makefile
  .SUFFIXES:            # 删除默认的后缀
  .SUFFIXES: .c .o .h   # 定义自己的后缀
  ```

- `-r` 、 `-no-builtin-rules`：使默认的后缀列表为空。

- `SUFFIXE`：定义默认的后缀列表。

  > 用 `.SUFFIXES` 来改变后缀列表，不要改变变量 `SUFFIXE` 的值。

## 隐含规则搜索算法

> 函数库文件模式，这个算法会被运行两次，第一次是找目标，如果没有找到的话，那么进入第二次，第二次会把 `member` 当作T来搜索。

1. 把T的目录部分分离出来。叫D，而剩余部分叫N。（如：如果T是 `src/foo.o` ，那么，D就是 `src/` ，N就是 `foo.o` ）
2. 创建所有匹配于T或是N的模式规则列表。
3. 如果在模式规则列表中有匹配所有文件的模式，如 `%` ，那么从列表中移除其它的模式。
4. 移除列表中没有命令的规则。
5. 对于第一个在列表中的模式规则：
   1. 推导其“茎”S，S应该是T或是N匹配于模式中 `%` 非空的部分。
   2. 计算依赖文件。把依赖文件中的 `%` 都替换成“茎”S。如果目标模式中没有包含斜框字符，而把D加在第一个依赖文件的开头。
   3. 测试是否所有的依赖文件都存在或是理当存在。（如果有一个文件被定义成另外一个规则的目标文件，或者是一个显式规则的依赖文件，那么这个文件就叫“理当存在”）
   4. 如果所有的依赖文件存在或是理当存在，或是就没有依赖文件。那么这条规则将被采用，退出该算法。
6. 如果经过第5步，没有模式规则被找到，那么就做更进一步的搜索。对于存在于列表中的第一个模式规则：
   1. 如果规则是终止规则，那就忽略它，继续下一条模式规则。
   2. 计算依赖文件。（同第5步）
   3. 测试所有的依赖文件是否存在或是理当存在。
   4. 对于不存在的依赖文件，递归调用这个算法查找他是否可以被隐含规则找到。
   5. 如果所有的依赖文件存在或是理当存在，或是就根本没有依赖文件。那么这条规则被采用，退出该算法。
   6. 如果没有隐含规则可以使用，查看 `.DEFAULT` 规则，如果有，采用，把 `.DEFAULT` 的命令给T使用。

# 函数库文件

> 对Object文件（程序编译的中间文件）的打包文件。

## 成员

一个函数库文件由多个文件组成。你可以用如下格式指定函数库文件及其组成。

```makefile
archive(member)
```

> 一般来说，这种用法基本上就是为了 `ar` 命令来服务的。
>
> e.g.
>
> ```makefile
> foolib(hack.o kludge.o): hack.o kludge.o
>     ar cr foolib hack.o kludge.o
> ```

## 隐含规则

1. 如果这个目标是 `a(m)` 形式的，会把目标变成 `(m)` 。

## 后缀规则

```makefile
.c.a:
    $(CC) $(CFLAGS) $(CPPFLAGS) -c $< -o $*.o
    $(AR) r $@ $*.o
    $(RM) $*.o
```

> 等价于
>
> ```makefile
> (%.o) : %.c
>     $(CC) $(CFLAGS) $(CPPFLAGS) -c $< -o $*.o
>     $(AR) r $@ $*.o
>     $(RM) $*.o
> ```
>
> 

## 注意

函数库打包文件生成时，小心使用make的并行机制（ `-j` 参数）。

如果多个 `ar` 命令在同一时间运行在同一个函数库打包文件上，就很有可以损坏这个函数库文件。

尽量不要用 `-j`。