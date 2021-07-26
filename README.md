# [ARP Spoofing]
<img src="https://img-blog.csdnimg.cn/20210726194704703.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dlaXhpbl80MzQxNTY0NA==,size_16,color_FFFFFF,t_70" width="400px">

# [Usage]

- Commands:
  - clear              clear the screen
  - cut                通过ARP欺骗切断局域网内某台主机的网络
  - exit               exit the program
  - help               display help
  - hosts              主机管理功能
  - loot               查看嗅探到的敏感信息
  - middle-attack      中间人攻击
  - scan               扫描内网中存活的主机
  - set                配置参数
  - show               展示信息
  - sniff              嗅探用户名和密码
  - webspy             嗅探http报文

# [Example]
首先启动redis数据库，然后：
- Linux/macOS : sudo go run main.go 
- windows: go run main.go
程序会自动连接redis数据库
- step0. show options 检查各项配置是否正确,如果配置不正确，可以使用 set key value 设置选项key的值为value
- step1. scan 扫描局域网中的主机
- step2. hosts 查看所有扫描到的主机
- step3  cut 向某台主机发送arp欺骗报文,启动后 发送数据包的协程将在后台默默运行
- step4  cut stop 停止发送

其他功能：
- webspy 可以嗅探所有流经本机网卡的http包,启动webspy前，建议向使用middleattack将目标主机的流量导过来
<img src="https://img-blog.csdnimg.cn/2021072619170764.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dlaXhpbl80MzQxNTY0NA==,size_16,color_FFFFFF,t_70"/>
- sniff 嗅探有敏感信息的http数据包，并存入redis，可以通过loot命令查看收集到的数据包

# [TODO] 

1. 查看所有 cutted 主机
2. 添加关闭中间人攻击的功能
3. 每个主机只能被中间人攻击一次（设置一个被攻击的主机集合）
4. 优化webspy的功能
5. 检查中间人攻击模块 篡改数据包的mac地址是否确实篡改了
如果确实篡改了，为什么 wireshark 抓不到
6. 解决本主机在启动中间人攻击模块后上网慢的问题
7. 设置一个过滤器，只监听欺骗双方的数据包,查下过滤器的语法