# 基于ARP欺骗实现的网络剪刀手

`sudo go run main.go`

建议以管理员身份运行，因为发数据包需要管理员权限

第一步:扫描同局域网下的主机
![输入图片说明](https://images.gitee.com/uploads/images/2021/0316/211103_a283bd95_8582605.png "屏幕截图.png")

第二步:选择网关（一般.1结尾的是网关）
![输入图片说明](https://images.gitee.com/uploads/images/2021/0316/211213_7866d1bb_8582605.png "屏幕截图.png")

第三步：点击CutOff切断目标主机的网络

![输入图片说明](https://images.gitee.com/uploads/images/2021/0316/211248_ef9a2b8f_8582605.png "屏幕截图.png")
（注意不要把自己的网给掐了 :sweat_smile: ）


有两种ARP欺骗思路：

- 方式一：向目标主机发送ARP欺骗报文，告诉它错误的网关MAC地址
- 方式二：向网关发送ARP欺骗报文，告诉它错误的目标主机MAC地址

第二种方式更加隐蔽,更加奏效
![输入图片说明](https://images.gitee.com/uploads/images/2021/0609/185822_c1afd466_8582605.png "屏幕截图.png")

