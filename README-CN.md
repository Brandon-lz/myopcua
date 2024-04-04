# myopcua
OPC-UA转HTTP小服务。随时随地部署，功能强大。

## 特点
- 轻量级，不依赖环境，无需安装，直接运行。
- 支持linux、windows、macos。
- 支持OPC-UA读写
- 稳定可靠，断线重连，不会崩溃。
- 高性能读取，支持最小10ms间隔读取，丝滑流畅。
- 支持多种数据格式，包括JSON、CSV、XML、Excel等。
- 支持多种认证方式，包括用户名密码、证书、令牌等。


## 快速开始

### 下载
找到对应的平台可执行文件下载： https://github.com/Brandon-lz/myopcua/releases/tag/v0.1.0

### 配置
按照config-example.toml文件配置好相关参数，然后保存为config.toml。


### 运行
下载后请注意添加执行权限，然后运行：

```
chmod +x myopcua-linux
./myopcua-linux
```
打开http://localhost:8080/docs/index.html#/ 即可访问API文档，按照上面的步骤配置好参数后，点击“Try it out”按钮，即可进行测试。
仅需2步即可轻松完成配置

## 联系方式
- 邮箱：advanced_to@163.com
- QQ：1791325181
- 商务合作加微信：advance_to  （软件本身开源免费，有问题请提issue，不是商业合作或加群请勿扰）
- 微信群（加微信备注：opcua加群）
