<div align="center">
<p align="center">
  <a href="./README.md">English</a> |
  <a href="./README-zh.md">简体中文</a> 
</p>
</div>

# myopcua
Flexible OPC-UA to HTTP API server. Portable and powerful.

## Features
- Lightweight, environment-independent, no installation required, run directly.
- Supported on Linux, Windows, macOS.
- Support OPC-UA read and write.
- Stable and reliable, supports reconnection when disconnected, will not crash.
- High-performance read, supports minimum 10ms interval read, smooth and fluid.
- Support multiple data formats, including JSON, CSV, XML, Excel, etc.
- Support multiple authentication methods, including username and password, certificates, tokens, etc.
## Quick Start
### Download
Find the corresponding platform executable file for download: https://github.com/Brandon-lz/myopcua/releases/tag/v0.1.0

### Configuration
Configure the relevant parameters according to the config-example.toml file and save it as config.toml.

### Running
After downloading, please add execution permissions and run:
```shell
chmod +x myopcua-linux
./myopcua-linux
```
Open http://localhost:8080/docs/index.html#/ to access the API documentation. After configuring the parameters as described above, click the “Try it out” button to start testing. It only takes 2 steps to easily complete the configuration.

## Contact Information
- Email: advanced_to@163.com
- QQ: 1791325181
- Business cooperation add WeChat: advance_to (The software itself is open source and free, please raise issues if there are problems, no disturbance for non-commercial cooperation or joining groups)
- WeChat group (Add WeChat with the note: opcua join group)