CH3_TASK := " 3.1 3.2 3.3 3.4 3.5 3.6 3.7 3.8 3.9 3.10 3.11 3.12 3.13 3.14 3.15 3.16 3.17 3.18 3.19 "
CH4_TASK :=" 4.1 4.2 4.3 4.4 4.5 4.6 4.7 4.8 4.9 4.10 4.11 4.12"
CH5_TASK := " 5.1 5.2 5.3 5.4 5.5 5.6 5.7 5.8 5.9 5.10 5.11 5.12 5.13 5.14 5.15 5.16 5.17 5.18 "
CH6_TASK := " 6.1 6.2 6.3 6.4 6.5 6.6 6.7 6.8 6.9 6.10 6.11 6.12 6.13 6.14 6.15 6.16 "
CH7_TASK := " 7.9 7.10 "

OTHER := " checkWord bundle unBundle spaceLineSep sampleForm clicheGenerator"

test: ${OTHER} ${CH3_TASK} ${CH4_TASK} ${CH5_TASK} ${CH6_TASK} ${CH7_TASK}

checkWord:
	@echo "checkWord 检查密码文件"
	awk -f ch3/checkpassword.awk /etc/passwd

bundle:
	@echo "bundle 打包到文件"
	awk '{ gsub(/.*\//, "", FILENAME); print FILENAME, $$0 > "bundle";} ' ch2/*

unBundle:
	@echo "unBundle 解压到文件夹"
	awk '$$1 != prev { close(prev); prev = $$1; } { print substr($$0, index($$0, " ") + 1)  > "unbundle/"$$1}' bundle

spaceLineSep:
	@echo "spaceLineSep 由空行分隔的记录"
	awk -f ch3/addressbook.awk ch3/file/addressbooktest.txt

sampleForm:
	@echo "sampleForm 简单报表"
	awk -f ch4/normal/prep.awk file/countries.txt | awk -f ch4/normal/form.awk

clicheGenerator:
	@echo "clicheGenerator 废话生成器"
	awk -f ch5/clichegenerator.awk ch5/file/clichetest.txt 5

sort:
	@echo "sort 排序算法测试"
	awk -f ch7/sort.awk ch7/file/sorttest.txt

3.1:
	@echo "习题 3.1 修改程序 sum3, 使它忽略空行。"
	awk -f ch3/sum.awk  ch3/file/sumtest.txt

3.2:
	@echo "习题 3.2 为数值添加更加一般的正则表达式. 它会如何影响运行时间?"
	@echo "todo"

3.3:
	@echo "习题 3.3 如果把第 2 个 for 语句的 numcol 测试拿掉, 会产生什么影响?"
	@echo "涉及到数值的强制类型转换，转为前缀数字或0，不会出现非字符的字段了。"

3.4:
	@echo "习题 3.4 写一个程序, 这个程序读取一个 条目–数额 对列表, 对列表中的每一个条目, 累加它的数额; 在结束时, 打印条目以及它的总数额, 条目按照字母顺序排列。"
	awk -f ch3/amount.awk  ch3/file/amounttest.txt

3.5: 3.6
	@echo "习题 3.5 根据比例决定星号的个数, 使得当数据过多时, 一行的长度不会超过屏幕的宽度。"

3.6:
	@echo "习题 3.6 修改 histogram, 把输入分拣到指定数量的桶中, 根据目前为止看到的数据调整每个桶的范围。"
	awk 'BEGIN { for (i = 1; i <= 2000; i++) print int(101*rand())} '| awk -f ch3/histogram.awk

3.7:
	@echo "习题 3.7 修复 sumcomma (带逗号的数字求和程序): 检查数字中的逗号是否处于正确的位置上。"
	awk -f ch3/addcomma.awk ch3/file/addcommatest.txt
	awk -f ch3/sumcomma.awk ch3/file/sumcommatest.txt

3.8:
	@echo "习题 3.8 将日期转换成某种形式, 这种形式允许你对日期进行算术运算, 例如计算某两个日期之间的天数。"
	awk -f ch3/date.awk ch3/file/datetest.txt

3.9: 3.12
	@echo "习题 3.9 利用前面提到过的程序 addcomma, 为金额加上逗号。"

3.10: 3.12
	@echo "习题 3.10 对于负的, 或特别大的金额, 程序 prchecks 处理得并不是很好. 修改程序: 拒绝金额为负的打印请求, 同时能将数额特别巨大的金额分成两行打印出来。"

3.11: 3.12
	@echo "习题 3.11 函数 numtowords 有时会在一行中连着打印两个空格, 还会打印出像 “one dollars”这样有错误的句子, 你会如何消除这些瑕疵?"

3.12:
	@echo "习题 3.12 修改 prchecks: 在适当的地方, 为金额的大写形式加上连字符, 比如 “twenty-one dollars”"
	awk -f ch3/prchecks.awk ch3/file/prcheckstest.txt

3.13:
	@echo "习题 3.13 如何修改这个程序, 使得它可以处理具有多种分隔符的文本?"
	awk -f ch3/p12check.awk ch3/file/p12checktest.txt

3.14:
	@echo "习题 3.14 增强 checkgen 的功能, 使得我们可以原封不动地向程序传递一段代码, 例如创建一 个 BEGIN 来设置字段分隔符。"
	awk -f ch3/checkgen.awk ch3/file/checkgentest.txt

3.15: 3.16
	@echo "习题 3.15 重写 compat, 不用 asplit, 而是用正则表达式来识别关键词, 内建函数等. 比较两个版本的复杂度与速度。"

3.16:
	@echo "习题 3.16 因为 awk 的变量不需要事先声明, 所以如果用户不小心把变量名写错了, awk 并不会检测到该错误。 写一个程序, 这个程序搜索文件中只出现一次的名字. 为了让这个程序更具实用价值, 你可能需要对函数的定义及其用到的变量花点心思。"
	awk -f ch3/compat.awk ch3/file/compat.awk

3.17:
	@echo "习题 3.17 比较不同版本 bundle 与 unbundle 的时间效率和空间效率, 这些不同的版本用到 了不同的头部信息与尾部信息, 对程序的性能与复杂性之间的折衷进行评价。"
	@echo "todo"

3.18:
	@echo "习题 3.18 修改第 1 个 awk 程序: 检查输入数据中是否包含魔术字符串 !!#。"
	awk -f ch3/addressbooksort1.awk ch3/file/addressbooktest.txt | sort | awk -f ch3/addressbooksort2.awk

3.19:
	@echo "习题 3.19 写一个命令 lookup x y, 该命令从已知的文件中打印所有符合条件的多行记录, 条件 是记录含有名字为 x 且值为 y 的项。"
	ch3/lookup date 1/5/87

4.1: 4.4
	@echo "习题 4.1 修改 form3 与 form4: 从别处获取日期, 而不是将日期硬编码到代码中。"

4.2: 4.4
	@echo "习题 4.2 由于四舍五入, 由 form3 与 form4 打印的项并不总是等于对应列的小计, 你会如何修正这个问题?"

4.3: 4.4
	@echo "习题 4.3 表格格式化程序假定所有数字的小数部分的位数都是相同的, 修改它, 使得即使这个假定不成立, 程序也可以正确地工作。"

4.4:
	@echo "习题 4.4 增强 table 的功能, 增强后的 table 允许输入数据中出现一个格式说明行序列, 这个序列说明了如何格式化每一列随后的数据. (tbl 就是这样控制输出格式的)。"
	awk -f ch4/hard/prep.awk pass=1 file/countries.txt pass=2 file/countries.txt | awk -f ch4/hard/form.awk

4.5:
	@echo "习题 4.5 修改 info: 参数通过 ARGV 传递进来, 而不是由 shell 进行替换。"
	ch4/info "USA|England"

4.6:
	@echo "习题 4.6 本节实现的 join 不会进行错误检查, 也不会检查文件是否是有序的. 修复这些问题, 在修复之后, 程序会变得多大?"
	sort file/countries.txt > file/temp.countries.txt
	sort file/capitals.txt > file/temp.capitals.txt
	awk -f ch4/join.awk file/temp.countries.txt file/temp.capitals.txt

4.7: 4.8
	@echo "习题 4.7 实现 join 的另一个版本, 它将一个文件整个读入内存, 然后再执行连接操作. 与原来的版本相比, 哪个更简单?"

4.8:
	@echo "习题 4.8 修改 join: 它可以按照输入文件的任意一个字段或字段组来进行连接, 并可以按照任意的顺序, 有选择地输出某些字段。"
	awk -f ch4/joincolumn.awk leftKey="COUNTRY" rightKey="COUNTRY" outputField="CONTINENT COUNTRY CAPITAL AREA POPULATION" file/capitals.txt file/countries.txt

4.9: 4.12
	@echo "习题 4.9 如果你的系统不支持 awk 的 system 函数, 请修改 qawk, 修改后的 qawk 把命令序列 写入一个或多个文件中, 这些文件可以被单独地执行。"

4.10: 4.12
	@echo "习题 4.10 构造导出表的时候, qawk 为每一个命令调用一次 system, 修改 qawk, 修改后的 qawk 把所有的构造表格的命令都收集到一个字符串中, 这样就可以只需要调用一次 system。"

4.11: 4.12
	@echo "习题 4.11 修改 qawk, 使得它可以检查即将作为输入数据的导出文件是否已经存在. 如果文件 已经存在, 并且基文件在导出文件被导出之后未曾修改过, 那么我们就可以直接使用该导出文件, 而不用重新计算. 可以参考第 七 章的程序 make。"

4.12:
	@echo "习题 4.12 提供一种输入并编辑多行查询的方式. 只要对 qawk 作些许修改, 就可以把多行查询收集起来. 一种可选的方案是调用你喜欢的编辑器, 另一种方式是用 awk 实现一个简单的编辑器。"
	awk -f ch4/qawk.awk ch4/file/qawktest.txt

5.1:
	@echo "习题 5.1 测试 rand 的输出是不是真的随机数。"
	@echo "答：不是"

5.2:
	@echo "习题 5.2 写一个程序, 该程序生成 1 到 n 之间的 k 个互不相同的随机整数, 要求程序的时间复杂度与 k 成正比。"
	awk -f ch5/randInts.awk 10 5

5.3:
	@echo "习题 5.3 写一个随机生成四手桥牌的程序。"
	awk -f ch5/poker.awk

5.4:
	@echo "习题 5.4 写一套语法规则, 该规则能够生成关于某学科的, 听起来貌似合理的文本 — 学科可以 是商业, 政治, 计算机或其他你感兴趣的领域。"
	awk -f ch5/sentgen.awk ch5/file/expression.txt depth=100 ch5/file/expressiontest.txt

5.5: 5.6
	@echo "习题 5.5 在某些语法规则下, 语句生成程序很有可能落入到这样一种境地: 推导过程越来越长, 却没有停下来的迹象, 添加一条机制, 使得程序可以限制推导过程的长度。"

5.6:
	@echo "习题 5.6 给语法规则加上权重, 使得对同一个非终结符来说, 它的各个展开规则被选中的概率是不同的。"
	awk -f ch5/sentgen.awk ch5/file/grammar.txt depth=10 ch5/file/grammartest.txt

5.7:
	@echo "习题 5.7 实现一个非递归的语句生成程序。"
	@echo "todo"

5.8:
	@echo "习题 5.8 除了加法外, 再新增几种数学运算. 另外, 如果用户的回答是错误的, 显示一条提示信息。"
	awk -f ch5/arith.awk - 100 auto

5.9:
	@echo "习题 5.9 修改 quiz, 使得同一道问题最多出现一次。"
	awk -f ch5/quiz.awk ch5/file/quiztest.txt name number auto

5.10:
	@echo "习题 5.10 修改单词计数程序: 不区分单词大小写, 于是 The 与 the 被当作是同一个单词。"
	awk -f ch5/sentgen.awk ch5/file/grammar.txt depth=10 ch5/file/grammartest.txt | awk -f ch5/wordfreq.awk

5.11:
	@echo "习题 5.11 写一个程序, 该程序计算某个文档中句子的个数及每个句子的长度。"
	awk -f ch5/linefreq.awk ch5/file/word.txt

5.12:
	@echo "习题 5.12 写一个 control-break 程序, 用来计算单词的个数. 与 wordfreq 相比, 它的性能表现 如何?。"
	awk -f ch5/sentgen.awk ch5/file/grammar.txt depth=10 ch5/file/grammartest.txt | awk '{gsub(/[.,:;!?(){}]/, ""); for (i = 1; i <= NF; i++) {print $i | "sort -rn"}}' | awk -f ch5/wordfreq2.awk

5.13:
	@echo "习题 5.13 修改 fmt: 对齐输出文本的右边空白。"
	awk -f ch5/fmt.awk ch5/file/fmttest.txt

5.14:
	@echo "习题 5.14 增强 fmt 的功能, 使得它可以通过识别文档中可能的标题, 列表等信息, 推断出文档的正确格式。 这次不是直接对文档进行格式化, 而是生成排版程序 (例如 troff, TEX 等) 的格式化命令。"
	@echo "todo"

5.15:
	@echo "习题 5.15 如果遗漏了符号名末尾的下划线, 会发生什么事？"
	@echo "如果有其他正常段落和缺失的符号一样，会被错误替换"

5.16:
	@echo "习题 5.16 修改 xref: 可以侦测到某个符号名的多次定义。"
	awk -f ch5/xref.awk ch5/file/document.txt > ch5/file/xref.temp.awk
	awk -f ch5/file/xref.temp.awk ch5/file/document.txt

5.17:
	@echo "习题 5.17 修改 xref, 使得它可以生成你所喜爱的文本编辑器或流式编辑器 (比如 sed) 的编辑 命令, 而非 awk 命令. 这会对编辑器的性能产生什么影响?"
	@echo "todo"

5.18:
	@echo "习题 5.18 你有没有办法让 xref 只需要对输入数据遍历一次? 【遍历一次】 对定义的放置位置而言, 隐含着什么限制条件?"
	@echo "符号名必须在使用前定义"
	awk -f ch5/xrefone.awk ch5/file/document.txt

5.19: 5.21
	@echo "习题 5.19 为 kwic 添加一个 【停止列表】 (stop list), 【停止列表】 指的是不能被当作关键词的单词集合, 比如 a, the。"

5.20: 5.21
	@echo "习题 5.20 修改 kwic, 使得它可以显示尽可能多的行内数据, 方法是当遇到一行的末尾时回绕,而非截断。"

5.21:
	@echo "习题 5.21 编写一个创建重要语汇索引 (concordance) 的程序: 对每一个重要的单词, 显示所有的出现了该单词的句子或短语。"
	ch5/kwic ch5/file/fmt.temp.txt

6.1: 6.3
	@echo "习题 6.1 修改 asm, 打印程序与内存的内容, 就像上面显示的那样。"

6.2: 6.3
	@echo "习题 6.2 增强解释器的功能, 打印指令的执行轨迹。"

6.3:
	@echo "习题 6.3 适当扩大汇编语言的规模, 比如添加错误处理代码与其他条件判断指令. 为了方便用户 使用, 你会怎么处理立即数, 比如 add = 1 (如果不支持立即数, 就必须要求用户自己创建一个 名为 one 的内存单元)?"
	awk -f ch6/asm.awk ch6/file/asm.txt auto

6.4:
	@echo "习题 6.4 写一个反汇编程序, 把内存中的内容转换成对应的汇编语言。"
	awk -f ch6/backasm.awk ch6/file/backasm.txt

6.5:
	@echo "习题 6.5 查看一台真实的机器 (比如 Apple-II 和 Commodore 的 6502 芯片, 或 IBM PC 及其兼容机的 8086 芯片族), 尝试为它的汇编语言子集写一个汇编程序。"
	@echo "todo"

6.6:
	@echo "习题 6.6 修改 graph: 对图进行转置, 即 x 轴在垂直方向上, 而 y 轴在水平方向上, 而且可以指 定是否以对数的方式显示 x 坐标或 y 坐标。"
	awk -f ch6/transpose.awk ch6/file/graph.txt | awk -f ch6/graph.awk

6.7: 6.8
	@echo "习题 6.7 为 graph 添加一条命令, 这条命令可以让 graph 从文件中读取数据。"

6.8:
	@echo "习题 6.8 自动计算刻度的位置。"
	awk -f ch6/transpose.awk ch6/file/graph.txt

6.9:
	@echo "习题 6.9 如果你的系统提供了图形库, 修改 graph, 为图形库生成命令。"
	echo "todo"

6.10:
	@echo "习题 6.10 实现一个新版的 sortgen, 它提供了访问 sort 全部功能的方法, 要求程序能够检测 出不一致的请求, 比如同时按照字典序与数值对数据进行排序。"
	echo "todo"

6.11:
	@echo "习题 6.11 在对输入语言规范性不作过多要求的前提下, 如何提高 sortgen 的精确度。"
	awk -f ch6/sortgen.awk ch6/file/sortgentest.txt

6.12:
	@echo "习题 6.12 写一个程序, 把一个排序命令翻译成对应的英文句子, 然后再把句子作为 sortgen 的输入。"
	echo "todo"

6.13:
	@echo "习题 6.13 为 calc2 的标准值 (例如 π 与 e) 添加内建变量. 添加一个内建变量, 用来表示最后一个输入行的运算结果. 添加两个栈操作运算符, 分别用来完成栈顶元素的复制, 与交换栈顶的两个元素。"
	awk -f ch6/calc2.awk  ch6/file/calc.txt

6.14: 6.16
	@echo "习题 6.14 构造一组测试集, 对 calc3 进行详尽地测试。"

6.15: 6.16
	@echo "习题 6.15 为中缀计算器添加指数运算, 内建函数和变量. 与逆波兰式计算器的实现作比较。"

6.16:
	@echo "习题 6.16 加强 calc3 的错误处理功能。"
	awk -f ch5/sentgen.awk ch5/file/expression.txt depth=100 ch5/file/expressiontest.txt |  awk -f ch6/calc3.awk

7.9:
	@echo "习题 7.9 修改 tsort, 使得它可以处理图中孤立的节点。"
	awk -f ch7/tsort.awk ch7/file/tsorttest.txt

7.10:
	@echo "习题 7.10 修改 rtsort, 按照通常的顺序打印节点, 也就是先打印前驱节点. 你是否可以在不修改 rtsort 的前提下达到同样的效果?"
	awk -f ch7/tdfssort.awk ch7/file/tsorttest.txt

.PHONY: test ${OTHER} ${CH3_TASK} ${CH4_TASK} ${CH5_TASK} ${CH6_TASK} ${CH7_TASK}
