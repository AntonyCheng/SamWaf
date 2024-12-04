docker run --rm -v "$PWD":/media/sf_SamWaf -w /media/sf_SamWaf -e CGO_ENABLED=1 -e GOPROXY=https://goproxy.cn,direct golang:1.21.4 go build -v -ldflags="-X SamWaf/global.GWAF_RELEASE=true -X SamWaf/global.GWAF_RELEASE_VERSION_NAME=20241204 -X SamWaf/global.GWAF_RELEASE_VERSION=v1.3.8 -s -w -extldflags "-static"" -o /media/sf_SamWaf/release/SamWafLinux64 main.go  && upx -9 /media/sf_SamWaf/release/SamWafLinux64