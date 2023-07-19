# diaoyu project
## datago
> 自己数据管理程序
> 
## 账号/密码
> admin/123456
> 
## 要使用上传图片功能请配置oss
> client, err := oss.New("https://oss-us-west-1.aliyuncs.com", "LTAI5tLn6jhY3Nt9TXGpuekK", "")
> bucket, err := client.Bucket("notea-us")
> attachment := model.NewAttachmentForInner("https://notea-us.oss-us-west-1.aliyuncs.com/" + path)
> 上述几行代码