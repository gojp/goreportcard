# Go Report Card

旧版存在的问题：

* 旧版不支持golangci-lint
* gometalinter 已经弃用了
* 不支持私库

预期维护的功能：

* 增加更多linter（golangci-lint 支持 40+ linter）
* 支持私库
* 生成可视化报告
* 分支可选
* 直接通过链接获取测试报告 （如 xxx.com/report?repo=git.medlinker.com/service/common）
* 静态代码检查评分
* 缓存repo的评估报告
