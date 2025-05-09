# PipControler(Pip掌控者)

*让安装第三方库不再繁琐，让Python开发更高效*

## 介绍:

PipControler旨在为了简化Python第三方库安装，它与传统Pip的区别：

| | PioControl | Pip |
| --- | --- | --- |
| 运行时自动安装 | Yes | No |
| 安装时间 | 准备运行时 | 运行前 |
| ... | ... | ... |

作者:*Yufeng Gosling*

语言:*Golang*

开源许可证:*GPLv3*

版本:*v0.2.0*

## 安装


安装工具链
```
# 基于Debain的发行版
sudo apt install git golang perl
# 基于RHEL的发行版
# Fedora 22, RHEL 8以下
sudo yum install git golang perl
# Fedora 22, RHEL 8及以上
sudo dnf install git golang perl
```

克隆仓库
```
cd /home/
git clone https://github.com/YufengGosling/PipControler.git
cd PipControler
```

编译
```
# 自动脚本
sudo bash install.sh

# 如果要手动编译
# 编译
go build -o bin/ -ldflags="-s -w" ./cmd/...
# 如果要debug
go build -o bin/ ./cmd/...
```

复制文件
```
cp -r scripts bin/scripts
```

添加环境变量
```
# bash
echo export PATH=$PATH:/home/PipControler/bin/ > ./bashrc
# 如果是zsh
echo export PATH=$PATH:/home/PipControler/bin/ > ./zshrc
```

## 使用
安装完成后，目前有下面命令可以使用:
| 命令 | 用处 |
| --- | --- |
| ipp | 扫描目录下的Python源码文件并自动安装依赖的第三方库 |
| pipcontroler | 显示信息(-v版本，-h帮助等) |
| 等待更多 | ... |

## 下一步打算
- 增量扫描
- 一键运行安装

## 贡献
如果您感兴趣，可以为该项目贡献代码，不过请遵循GPLv3许可证

# 许可证
该项目使用GPLv3许可证，详细请看[LICENSE](http://github.com/YufengGosling/PipControler/LICENSE)
