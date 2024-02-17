# 本仓库包含两部分内容
## practice文件夹下是前期的练习内容
## go_redis文件夹下是使用Go语言实现Redis功能的代码。

# 实现Redis协议解析器
## REdis Serialization Protocol(RESP)
* 正常回复
* 错误回复
* 整数
* 多行字符串
* 数组
### 正常回复
* 以 + 开头 以 \r\n 结尾的字符串形式
```bash
+ok/r/r
```
### 错误回复
* 以 - 开头，以 /r/n 结尾的字符串形式
```bash
-Error message\r\n
```
### 整数
* 以：开头，以/r/n结尾的字符串形式
```bash
:123456\r\n
```
### 字符串
* 以$开头，后跟实际发送字节数，以/r/n结尾
  "ShiShi"
```bash
$6\r\nShiShi\r\n
```
### 数组
* 以 “*” 开头，后跟成员个数
* SET key value
```bash
*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n 
```