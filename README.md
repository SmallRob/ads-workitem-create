# ads-workitem-create
add workitem in ads，azure，which can create workitem and close theme

# how to use
# pckage files
> CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ads-linux
> CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ads.exe
> CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ads-osx

# cli
work --file ./logs/WorkTemplate202204.txt
work --gti xxx
work --close ./logs/workid.log

# make a work template 生成工时模板
ads work t

# create item 创建工作项
ads work a

# close item 关闭工作项
ads work c
