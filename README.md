# mongosync

### 快速开始
```
./mongosync -src=192.168.30.249:28717 -dst=localhost:27017 -name=test
```

### 帮助
```
./mongosync -h
  -dst string
        目标数据库地址,支持任何 mongo 官方支持的连接字符串 (default "localhost:27017")
  -excludes ，
        只有不在此集合中的 NS 才会被同步， 多个可以使用 ， 分割； 可以使用 `dbName.*`（目前只支持这一种格式）只匹配某条数据库下面的 ns
  -h    帮助信息
  -includes ，
        只有在此集合中的 NS 才会被同步, 多个可以使用 ， 分割; 可以使用 `dbName.*`（目前只支持这一种格式）只匹配某条数据库下面的 ns
  -interval int
        同步间隔时间; unit: second (default 60)
  -limit int
        每次从oplog.rs读取多少条数据进行转化 (default 1000)
  -name string
        转换上下文的名字, 推荐为每个转换设置一个特殊的名字; 默认值为 dst 参数
  -op-str ,
        加载哪些 op type 的数据进行转换， 默认以 , 分割 (default "i,u,d")
  -src string
        数据源数据库地址,支持任何 mongo 官方支持的连接字符串 (default "localhost:27017")
  -update-ts-len int
        转换多少条数据同步一次 mongo.sync.log 里面的 ts 参数， 该 ts 参数用于下次获取数据的起点 (default 10)
  -v    版本信息
```

### 特性
* 根据 `oplog.rs` 增量同步主库内容到从库(非主从关系，只是一个比喻)
* 同步失败的记录和原因会存在数据库中，保证我们能够追溯原因



### 注意
* `interval` 参数劲量超过 5 s, 因为 `oplog.rs` 数据较多，查询可能会比较耗时，实际情况应该结合 `limit` 参数一起制定