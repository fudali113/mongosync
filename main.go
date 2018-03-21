package main

import "flag"

func main() {
	var src, dst string
	flag.StringVar(&src, "src", "localhost:27017", "数据源数据库地址")
	flag.StringVar(&dst, "dst", "localhost:27018", "目标数据库地址")

}
