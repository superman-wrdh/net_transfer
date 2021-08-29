## tcp文件传输

#### 功能
    tcp传输文件 支持大文件,支持文件和文件夹,单文件最大支持4GB

#### 编译
    Macos sh  build-darwin.sh
    Linux sh build-linux.sh 
    编译后可执行文件输出到target目录

#### 启动
    服务端 ./net_transfer 
    客户端 ./net_transfer -op=client -ip=your ip

#### 效果图
![使用截图](doc/image20210829.png)