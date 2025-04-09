#!/bin/bash

# 关闭以太网接口 eth0
ifconfig eth0 down

# 发送 AT 命令到 /dev/ttyUSB2
echo "AT+QNETDEVCTL=3,1,1" > /dev/ttyUSB2

# 启用 usb0 接口
ip link set usb0 up

# 配置 usb0 接口的 IP 地址和子网掩码
ifconfig usb0 192.168.43.100 netmask 255.255.255.0

# 添加默认路由
ip route add default via 192.168.43.1 dev usb0

echo "Network configuration completed."
