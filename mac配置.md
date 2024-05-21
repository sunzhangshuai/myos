# 安装iterm2

1. 安装：访问网址：[iterm2](https://iterm2.com/)，点击Downloads。
2. 配置字体：`settings -> profiles -> text -> font`。可以修改字体和大小。字体建议：`MesloLGS NF`，大小18。
3. 配置样式
   1. 下载 [样式库](# https://github.com/eillsu/iTerm2-Chinese-Tutorial 下载样式库)
   2. `settings -> profiles -> Colors -> Color Presets -> improt`，找到下载的 `schemes` 文件夹，选择 `Solarized Dark Higher Contrast.itermcolors`。再切换到 `Solarized Dark Higher Contrast` 。

# 安装homebrew

1. 安装

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

2. 相关命令

```shell
brew –help 								# 查看帮助
brew install wget 				# 下载软件
brew uninstall git 				# 卸载软件
brew search git 					# 搜索软件
brew list 								# 查看软件列表
brew update 							# 升级homebrew本身
brew outdated 						# 查看软件可更新的版本
brew upgrade 							# 更新所有包
brew upgrade git 					# 更新指定包
brew cleanup             	# 清理所有包的旧版本
brew cleanup git    			# 清理指定包的旧版本
brew cleanup -n          	# 查看可清理的旧版本包
brew pin git      				# 锁定包，不再更新
brew unpin git   					# 取消锁定
brew info git    					# 显示某个包的信息
brew info             		# 显示安装了包数量，文件数量，和总占用空间
brew deps git 						# 查看已安装的包的依赖，树形显示
```

3. 软件路径

```shell
/opt/homebrew/Cellar/[软件名]/[版本号]/   		# 软件位置
/opt/homebrew/bin														# 软件可执行文件的软链位置
/opt/homebrew/etc														# 软件的配置文件位置
/opt/homebrew/var														# 软件的数据位置	
```

4. 更多软件的下载库。如果有些软件brew上不存在，可以新增下载源。

```shell
brew tap 																		# 更新已存在的软件库，并列出
brew tap elastic/tap												# 新增软件库
brew untap elastic/tap											# 删除软件库
brew install elastic/tap/elasticsearch-full	# 安装软件库中的软件，需要加库名前缀
```

```shell
/opt/homebrew/Library/Taps																# 软件库位置
/opt/homebrew/Library/Taps/elastic/homebrew-tap/Formula		# 三方库的下载列表
```

# 安装oh-my-zsh

1. 安装

```shell
sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
```

2. 变更shell

```shell
brew install zsh				# 安装zsh
chsh -l  								# 查看可以用的shell
chsh -s	[shell路径]			 # 指定默认shell
# 重启iterm2
```

3. 安装插件

```shell
cd ~/.oh-my-zsh/custom/plugins     																	# 进入插件目录
git clone https://github.com/zsh-users/zsh-autosuggestions   				# 下载自动补全
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git	# 下载语法高亮
plugins=(git zsh-autosuggestions zsh-syntax-highlighting)						# 编辑~/.zshrc文件，修改plugins
source ~/.zshrc
```

4. 安装theme

```shell
cd ~/.oh-my-zsh/custom/themes																				# 进入样式目录
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git		# 下载p10k
ZSH_THEME="powerlevel10k/powerlevel10k"															# 编辑~/.zshrc文件，修改ZSH_THEME
source ~/.zshrc																											# 根据提示配置p10k
# 不满意可以重新配置
p10k configure
```

# 安装vim

1. 安装

```shell
brew install vim
brew install macvim
```

2. 修改配置文件 `~/.vimrc`

```shell
"显示行号
set nu

"启动时隐去援助提示
set shortmess=atI

"语法高亮
syntax on

"不需要备份
set nobackup

set nocompatible

"没有保存或文件只读时弹出确认
set confirm

"鼠标可用
set mouse=a

"tab缩进
set tabstop=4
set shiftwidth=4
set expandtab
set smarttab

"文件自动检测外部更改
set autoread

"c文件自动缩进
set cindent

"自动对齐
set autoindent

"智能缩进
set smartindent

"高亮查找匹配
set hlsearch

"显示匹配
set showmatch

"显示标尺，就是在右下角显示光标位置
set ruler

"去除vi的一致性
set nocompatible

"设置键盘映射，通过空格设置折叠
nnoremap <space> @=((foldclosed(line('.')<0)?'zc':'zo'))<CR>
""""""""""""""""""""""""""""""""""""""""""""""
"不要闪烁
set novisualbell

"启动显示状态行
set laststatus=2

"浅色显示当前行
autocmd InsertLeave * se nocul

"用浅色高亮当前行
autocmd InsertEnter * se cul

"显示输入的命令
set showcmd

"被分割窗口之间显示空白
set fillchars=vert:/
set fillchars=stl:/
set fillchars=stlnc:/

" vundle 环境设置
filetype off
set rtp+=~/.vim/bundle/Vundle.vim
"vundle管理的插件列表必须位于 vundle#begin() 和 vundle#end() 之间
call vundle#begin()
Plugin 'VundleVim/Vundle.vim'
Plugin 'altercation/vim-colors-solarized'
Plugin 'tomasr/molokai'
Plugin 'vim-scripts/phd'
Plugin 'Lokaltog/vim-powerline'
Plugin 'octol/vim-cpp-enhanced-highlight'
Plugin 'Raimondi/delimitMate'
Plugin 'VundleVim/YouCompleteMe'
" 插件列表结束
call vundle#end()
filetype plugin indent on

" 配色方案
set background=dark
colorscheme torte
"colorscheme molokai
"colorscheme phd

" 禁止显示菜单和工具条
set guioptions-=m
set guioptions-=T

" 总是显示状态栏
set laststatus=2

" 禁止折行
set nowrap

" 设置状态栏主题风格
let g:Powerline_colorscheme='solarized256'

syntax keyword cppSTLtype initializer_list

" 基于缩进或语法进行代码折叠
"set foldmethod=indent
set foldmethod=syntax
" 启动 vim 时关闭折叠代码
set nofoldenable

"允许用退格键删除字符
set backspace=indent,eol,start

"编码设置
set encoding=utf-8

"共享剪切板
set clipboard=unnamed

" Don't write backup file if vim is being called by "crontab -e"
au BufWrite /private/tmp/crontab.* set nowritebackup nobackup
" Don't write backup file if vim is being called by "chpass"
au BufWrite /private/etc/pw.* set nowritebackup nobackup
```

3. 安装插件管理器

```shell
cd ~/.vim/bundle																			# 进入vim插件目录
git clone https://github.com/VundleVim/Vundle.vim.git	# 下载插件目录管理器
```

4. 安装插件。

```shell
# 要安装的插件放在 call vundle#end() 和 call vundle#end() 之间，参考.vimrc
vim								# 进入vim编辑器
:PluginInstall		# 安装插件
:PluginList				# 插件列表
:PluginUpdate     # 更新所有插件
:PluginClean			# 修改.vimrc文件后，执行可以清理不需要的插件
:h vundle					# 查看帮助文档
```



4. `YouCompleteMe` 插件需要单独安装。这是个vim代码提示的插件，比较强大。

```shell
cd ~/.vim/bundle																				# 进入vim插件目录
git clone https://github.com/Valloric/YouCompleteMe.git	# 下载插件
cd ~/.vim/bundle/YouCompleteMe													# 进入插件位置
git submodule update --init --recursive                 # 下载子模块
./install.py --all																			# 安装
```

# 安装JetBrains

1. 安装 [Jetbrains Toolbox App](https://www.jetbrains.com/toolbox-app/)
2. 在 Jetbrains Toolbox App 安装需要的idea。
3. idea插件。

```shell
# Chinese (Simplified) Language Pack / 中文语言包		汉化
# GitToolBox																				git加强
# Material Theme UI																	样式库，在工具 -> Material Theme 切换样式，可以选 Ocenain
# Translation																				翻译，选中文本右键可翻译
# Atom Material Icons																icon加强
```

4. 字体
   1. `设置 -> 编辑器 -> 配色方案 -> 配色方案字体`，选择 `Jetbraibs Mono` 字体，开启连写，大小18，行高1.4
   2. `设置 -> 编辑器 -> 配色方案 -> 控制台字体`，修改控制台字体。
   3. `设置 -> 外观与行为 -> 外观`，修改项目栏等的字体和大小。