# Freev2ray

1. 自动从 https://github.com/freefq/free www.youhou8.com 等免费发布站点抓取vmess、trojan服务器节点，并建立代理。
1. 内置配置文件支持透明代理，非常方便集成到Openwrt。
1. 如果有其他免费站点，需要增加支持，可以提交issue。

## build
```sh
git clone https://github.com/xiqinging/freev2ray
cd freev2ray
go mod tidy
cd main
go build -o freev2ray
```


## 使用
```
 $ ./freev2ray -h

Usage:
    ./freev2ray [Options] <Subcommand>

Options:
    -d,--default-config    the default config file

Subcommand:
    b64trojan              trojan outbound, use base64 fetcher
    b64vmess               vmess outbound, use base64 fetcher
    zkqtrojan              trojan outbound, use zkq fetcher
```


