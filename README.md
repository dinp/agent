Agent
==========

这个模块部署在资源池的所有节点机器上，主要用于收集数据做汇报

收集的数据分成两部分

- 机器剩余内存，调度模块拿到这个信息之后才能做调度
- container情况，这是读取的Docker Daemon的接口，list所有container

对于container这块，目前只是拿到了本机的container列表，得知这些container的PublicPort，汇报给server，server把路由信息写入redis中的路由表

#### Q1：如何得知某个container是哪个app的呢？

刚开始的做法是把app的名称写入image的url中，container本身是可以知道使用了哪个image创建的，这样就知道了app与container的对
应关系。但是这个规范比较强，太具侵入性

现在的做法是把app的名称写入ENV，在创建container的时候写入，之后再通过inspect拿到ENV["APP\_NAME"]

## 注意：

- agent会把本机的ip汇报给server，server用此构建路由表，那本机的ip是如何获取的呢？多个ip的情况怎么办？现在的做法是只拿内网ip，然后过滤掉回环ip，只要网卡名称使用eth打头的，过滤掉Docker创建的虚拟网卡，如果用户配置了localIp，就直接使用用户配置的localIp
- agent是个前台进程，真正线上部署的时候可以使用god或者supervisor之类的管理

## 配置项说明

- **debug**: true/false 只影响打印的log
- **localIp**: 本机Ip地址，server会用这个ip地址和docker daemon通信
- **servers**: server的地址
- **interval**: 心跳周期，单位是秒
- **timeout**: 连接server的超时时间，单位是毫秒
- **docker**: Docker Daemon的接口地址，推荐Docker Daemon监听一个127.0.0.1的tcp接口，unix socket也可以，不过要注意文件权限了

## install

```
mkdir -p $GOPATH/src/github.com/dinp
cd $GOPATH/src/github.com/dinp; git clone https://github.com/dinp/agent.git
cd agent
go get ./...

# check cfg.json, depend docker daemon and server
./control start
```

