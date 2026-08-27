# league-sched — Go 语言循环赛赛程生成与积分榜 HTTP 后端服务，支持赛程编排、比赛结果录入和实时排名模拟

给定参赛队名单，按圈法排出单循环或主客场双循环赛程；再读入比分计算积分、净胜球与排名。队数为奇数须轮空、比分非法须拒绝入表；同一批结果按积分→净胜→进球→队名排序，模拟排名与手工录入必须一致。

## 构建 / 运行 / 测试

```text
go build ./...                                        # 编译
go run . -addr :8080                                  # 启动 HTTP 服务（/api/fixtures, /api/table, /api/simulate）
go test ./...                                         # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
