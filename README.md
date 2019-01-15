## 一、个人需求来源：
我在多个机器上的代码有同步需求，往往在一个地方写的代码不好同步到第二个地方，当然现在有很多同步方式例如 github、oneNote 等等，但是感觉起来如果自己能捣腾一个轻量、好使的工具，能够实时的通过请求将文件同步过来还是蛮爽的，正好最近看了下 P2P 、内网穿透的东西，感觉很牛X，所以这里个人立项了 ^_^ !


## 二、项目计划：
### 阶段一：
实现在服务器中转的情况下同步文件数据；

### 阶段二：
内网穿透，在底层硬件支持的情况下实现 P2P 文件传输；


## 三、服务器中转文件计划：
### client 端：
> * 轮询 -- 实现本地目录的递归扫描，并以一定的数据结构保存文件目录及结构（树）；
> * 与 server 端建立连接，定时上报本地的文件数据，保证服务端的数据一致性；
> * 完成文件对比算法，校验与对比文件的一致性，可以调研下能否通过 git 解决；

### server 端：
> * 与 client 保持长连接，建立 Map 索引会话实时同步数据；
> * 文件中转，A、B 多个客户端之间的文件交互需要通过这里上传与下载（可以尝试文件读取 && 文件上传两种方式）；



