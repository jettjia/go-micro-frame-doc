 

# 案例

 ![image-20220816203733917](01-goconvey.assets/image-20220816203733917.png)

```
$ go test -v 查看更多信息

$ goconvey  浏览器查看
```

 ![image-20220816204020687](01-goconvey.assets/image-20220816204020687.png)



文档：https://github.com/smartystreets/goconvey/wiki/Assertions#general-equality

比如查看 equal 的判断





# 测试报告

```shell
go test ./... -v -coverprofile=report/cover.out 2>&1 | go-junit-report > report/ut_report.xml

gocov convert report/cover.out | gocov-xml > report/coverage.xml

gocov convert report/cover.out | gocov-html > report/coverage.html
```

