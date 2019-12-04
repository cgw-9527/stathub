## 特性

- Golang 编写
- 只需一个文件就可以完成一切
- 支持SSL安全，支持域名访问，安全并简单


添加一个客户端

####  Linux (如果有wget)

wget --no-check-certificate -O - http://status.ulord.one:15944/node | sh

####  Linux (如果有curl)

curl --insecure http://status.ulord.one:15944/node | sh

