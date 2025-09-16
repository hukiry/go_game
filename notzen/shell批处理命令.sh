一、基础命令
	命令		作用			示例
	echo	输出到终端		echo "Hello World"
	pwd		显示当前路径		pwd
	cd		切换目录			cd ~/Downloads
	ls		列出目录内容		ls -l（长列表）ls -a（显示隐藏文件）
	clear	清屏				clear
	date	查看日期和时间	date
	whoami	当前用户			whoami
	uname	系统信息			uname -a
二、文件和目录操作
	命令		作用				示例
	mkdir	创建目录			mkdir new_folder
	rmdir	删除空目录		rmdir old_folder
	rm -r	删除目录及内容	rm -r folder_name
	touch	创建空文件		touch file.txt
	cp		复制文件/目录	 	cp file1.txt file2.txt / cp -r dir1 dir2
	mv		移动/重命名		mv old.txt new.txt
	rm		删除文件			rm file.txt
	cat		查看文件内容		cat file.txt
	less/more	分页查看文件	less file.txt
	head	查看文件开头		head -n 10 file.txt
	tail	查看文件结尾		tail -n 10 file.txt
三、文本处理
	命令	 	作用					示例
	grep	查找文本				grep "keyword" file.txt
	sed		流编辑器，文本替换	sed 's/old/new/g' file.txt
	awk		文本分析和处理		awk '{print $1}' file.txt
	cut		截取文本				cut -d ',' -f 1 file.txt
	sort	排序					sort file.txt
	uniq	去重					`sort file.txt
	wc		统计字数/行数/字符	wc -l file.txt
四、流程控制
	命令 / 语法						作用				示例
	if ... then ... fi				条件判断			if [ -f file.txt ]; then echo "存在"; fi
	if ...; then ...; else ...; fi	带 else 条件		if [ -d dir ]; then echo "目录"; else echo "不存在"; fi
	[ condition ]					条件表达式		-f file（文件存在）-d dir（目录存在）
	for var in ...; do ...; done	循环				for f in *.txt; do echo $f; done
	while ...; do ...; done			循环 			while [ condition ]; do ...; done
	until ...; do ...; done			循环直到条件成立	until [ -f file ]; do sleep 1; done
	case ... in ... esac			多分支选择		case $var in 1) echo "one";; 2) echo "two";; esac
	break							跳出循环			break
	continue						继续下一次循环	continue
五、变量与运算
	语法					说明				示例
	VAR=value			定义变量			name="Tom"
	$VAR				访问变量			echo $name
	${VAR}				推荐用法			echo ${name}
	read VAR			读取用户输入		read name
	export VAR=value	环境变量			export PATH=$PATH:/usr/local/bin
	let 或 $(( ))		数学运算			a=5; b=3; c=$((a+b))
	expr				另一种数学运算	expr 5 + 3
六、函数
	myfunc() {
	    echo "Hello $1"
	}
	myfunc Tom

	$1、$2 等表示传入参数

	return 返回状态码（0表示成功，非0表示失败）


七、文件测试与判断条件
	条件				含义				示例
	-f file			是否为文件		[ -f file.txt ]
	-d dir			是否为目录		[ -d folder ]
	-e path			文件或目录是否存在[ -e file_or_dir ]
	-r file			是否可读			[ -r file.txt ]
	-w file			是否可写			[ -w file.txt ]
	-x file			是否可执行		[ -x file.sh ]
	string1=string2	字符串相等		[ "$a" = "$b" ]
	-z string		字符串长度为0		[ -z "$a" ]
	-n string		字符串长度不为0	[ -n "$a" ]
八、文件权限
	命令		作用		示例
	chmod	改权限	chmod +x script.sh（加可执行）
	chown	改属主	chown user:group file.txt
	chgrp	改属组	chgrp staff file.txt
九、系统与进程
	命令		作用			示例
	ps		查看进程		ps aux
	top		实时进程		top
	kill	杀死进程		kill -9 PID
	df		查看磁盘空间	df -h
	du		查看目录占用	du -sh folder
	free	查看内存		free -h
十、压缩与打包
	命令		作用			示例
	tar		打包/解包	tar -cvf archive.tar folder / tar -xvf archive.tar
	tar.gz	压缩			tar -czvf archive.tar.gz folder / tar -xzvf archive.tar.gz
	zip		压缩			zip -r archive.zip folder
	unzip	解压			unzip archive.zip
十一、重定向与管道
	语法	含义					示例
	>	输出重定向（覆盖）	echo hello > file.txt
	>>	输出重定向（追加）	echo world >> file.txt
	<	输入重定向			wc -l < file.txt
	`	`					管道
	2>	错误重定向			command 2> error.log
	&>	输出和错误同时重定向	command &> output.log
十二、常用工具
	工具		功能				示例
	java	Java 执行		java -jar apktool.jar
	adb		Android 设备操作	adb install app.apk
	curl	下载/请求		curl -O url
	wget	下载文件			wget url
	find	查找文件			find . -name "*.txt"
	xargs	命令参数传递		`find . -name "*.txt"