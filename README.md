## SSH Tunnel

#### 编译
> 会在bin目录生成可执行文件
```sh
make && cd bin/
```

#### 配置文件（tunnel.yaml）
```yaml
tunnels: # root
  elasticsearch: # 名称
    auth: /path/to/.ssh/id_rsa # 密钥的绝对地址
    local: 0.0.0.0:9200 # 本地监听地址
    server: root@182.92.193.203:22 # 测试服务器地址
    remote: es-cn-v641aso9d001eihdi.elasticsearch.aliyuncs.com:9200 # elasticsearch的地址
```

#### 启动
```
Tunnel.exe -h
Tunnel.exe -c ./tunnel.yaml
```