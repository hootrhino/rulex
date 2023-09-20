# 软路由
## 简介
RhinoPi 的软路由配置，该配置主要基于Ubuntu Ip table实现，理论上来说只要有多个网卡就适用，但是当前该功能仅适配于Rhino系列的产品，如果需要移植请注意网卡参数。
## 环境要求
如果你需要移植这个功能到你自己的设备商，需要安装`dnsmasq`:
```sh
sudo apt install dnsmasq
```
同时吧Linux自带的DNS服务器关了
```sh
sudo systemctl disable systemd-resolved
sudo systemctl stop systemd-resolved
```
启动：
```sh
sudo systemctl start dnsmasq
sudo systemctl restart dnsmasq
```
最后测试：
```sh
rhino@RH-PI1:~$ dig A www.baidu.com

; <<>> DiG 9.11.3-1ubuntu1.16-Ubuntu <<>> A www.baidu.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 8809
;; flags: qr rd ra; QUERY: 1, ANSWER: 3, AUTHORITY: 13, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;www.baidu.com.                 IN      A

;; ANSWER SECTION:
www.baidu.com.          99      IN      CNAME   www.a.shifen.com.
www.a.shifen.com.       79      IN      A       182.61.200.7
www.a.shifen.com.       79      IN      A       182.61.200.6

;; AUTHORITY SECTION:
com.                    66948   IN      NS      j.gtld-servers.net.
com.                    66948   IN      NS      f.gtld-servers.net.
com.                    66948   IN      NS      a.gtld-servers.net.
com.                    66948   IN      NS      c.gtld-servers.net.
com.                    66948   IN      NS      l.gtld-servers.net.
com.                    66948   IN      NS      m.gtld-servers.net.
com.                    66948   IN      NS      e.gtld-servers.net.
com.                    66948   IN      NS      b.gtld-servers.net.
com.                    66948   IN      NS      h.gtld-servers.net.
com.                    66948   IN      NS      k.gtld-servers.net.
com.                    66948   IN      NS      i.gtld-servers.net.
com.                    66948   IN      NS      g.gtld-servers.net.
com.                    66948   IN      NS      d.gtld-servers.net.

;; Query time: 5 msec
;; SERVER: 192.168.199.1#53(192.168.199.1)
;; WHEN: Wed Sep 20 15:43:43 CST 2023
;; MSG SIZE  rcvd: 325

```
更多教程参考这里：https://computingforgeeks.com/install-and-configure-dnsmasq-on-ubuntu

## 原理
软路由实现和网卡有关，RhinoPi有以下网卡：
- eth0: 以太网口0
- eth1: 以太网口1
- lo: 本地回环
- usb0: 4G网卡
- wlan0: WIFI网卡

而软路由的实现原理就是利用这几个网卡互相进行桥接转发。