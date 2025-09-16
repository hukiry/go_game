一、基本文件和目录操作
	命令	用法示例	说明
	echo		echo Hello World		输出文本到控制台
	@echo off	@echo off				关闭命令回显（常放在开头）
	cd / chdir	cd C:\MyFolder			切换当前目录
	md / mkdir	mkdir build				创建目录
	rd / rmdir	rmdir /s /q build		删除目录（/s 删除子目录，/q 安静模式）
	del			del *.o /q				删除文件（/q 安静模式）
	copy		copy src.txt dest.txt	复制文件
	xcopy		xcopy src dest /s /i /y	复制文件夹及内容
	move		move src dest			移动文件或目录
二、变量与环境
	命令	示例	说明
	set			set VAR=value					定义变量
	set /p		set /p VAR=请输入:				从用户输入获取变量
	setlocal	setlocal enabledelayedexpansion	启用延迟变量扩展
	%VAR%		echo %VAR%						访问变量
	!VAR!		echo !VAR!						延迟扩展模式下访问变量
三、条件判断
	命令	示例	说明
	if				if "%VAR%"=="1" echo yes			条件判断字符串/数字
	if exist		if exist file.txt echo Found		判断文件/目录是否存在
	if not exist	if not exist folder mkdir folder	条件不存在时执行
	if /i			if /i "%VAR%"=="Y"					忽略大小写比较

	注意：在 Windows 批处理中，if 和圆括号 () 的用法必须在同一行开始，否则容易报错。推荐用 goto 跳转方式处理多行逻辑。

四、循环与跳转
	命令	示例	说明
	for		for %%f in (*.c) do echo %%f	遍历文件或列表
	goto	goto label						跳转到标签
	:label	:build_windows					定义跳转标签
	call	call other.bat					调用另一个批处理文件
五、执行外部命令
	命令	示例	说明
	start		start notepad.exe							启动程序或新窗口
	cmd /c		cmd /c dir									执行命令后退出
	powershell	powershell -Command "Compress-Archive ..."	调用 PowerShell 命令
六、错误处理
	命令	示例	说明
	exit		exit /b 1					退出批处理（可带返回码）
	errorlevel	if errorlevel 1 echo fail	判断上条命令的返回码
	`											`
七、其他常用命令
	命令	示例	说明
	pause		pause				暂停，等待用户按键
	rem			rem 注释内容			批处理注释
	::			:: 注释内容			另一种注释方法
	title		title Lua Build		设置 CMD 窗口标题
	color		color 0A			改变字体颜色
	time / date	time /t				获取当前时间/日期
	call :label	call :subroutine	调用子程序标签
	shift		shift				移动参数变量（%1→%2）
	set /a		set /a SUM=1+2		算术运算