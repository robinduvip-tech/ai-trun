#!/bin/bash
# CCX 全量打包脚本
# 用法: bash build-all.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DESKTOP_DIR="$SCRIPT_DIR/desktop"
FRONTEND_DIR="$SCRIPT_DIR/frontend"
BACKEND_DIR="$SCRIPT_DIR/backend-go"
DIST_DIR="$SCRIPT_DIR/dist"

# 工具路径
GO="/c/Program Files/Go/bin/go.exe"
BUN="/c/Users/dubin/AppData/Local/Microsoft/WinGet/Packages/Oven-sh.Bun_Microsoft.Winget.Source_8wekyb3d8bbwe/bun-windows-x64/bun.exe"
WAILS3="$HOME/go/bin/wails3.exe"
MAKENSIS="/c/Program Files (x86)/NSIS/makensis.exe"

# Go 环境
export GOPATH="$SCRIPT_DIR/.gocache"
export GOMODCACHE="$SCRIPT_DIR/.gomodcache"
export GOTMPDIR="$SCRIPT_DIR/.gotmp"
export GOPROXY="https://goproxy.cn,direct"

# PATH
export PATH="/c/Program Files/Go/bin:/c/Users/dubin/AppData/Local/Microsoft/WinGet/Packages/Oven-sh.Bun_Microsoft.Winget.Source_8wekyb3d8bbwe/bun-windows-x64:$HOME/go/bin:/c/Program Files (x86)/NSIS:$PATH"

echo "=========================================="
echo "  CCX 全量打包"
echo "=========================================="

# Step 1: 安装主前端依赖并构建
echo ""
echo "[1/4] 构建主前端..."
cd "$FRONTEND_DIR"
"$BUN" install
"$BUN" run build

# Step 2: 构建后端 (ccx-go.exe)
echo ""
echo "[2/4] 构建后端 ccx-go.exe..."
mkdir -p "$DIST_DIR"
rm -rf "$BACKEND_DIR/frontend/dist"
mkdir -p "$BACKEND_DIR/frontend/dist"
cp -r "$FRONTEND_DIR/dist/"* "$BACKEND_DIR/frontend/dist/"

cd "$BACKEND_DIR"
VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "dev")
BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(cd "$SCRIPT_DIR" && git rev-parse --short HEAD 2>/dev/null || echo "unknown")

CGO_ENABLED=1 "$GO" build \
  -ldflags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -s -w" \
  -o "$DIST_DIR/ccx-go.exe" .

echo "  -> $DIST_DIR/ccx-go.exe"

# Step 3: 复制后端到桌面目录
cp "$DIST_DIR/ccx-go.exe" "$DESKTOP_DIR/bin/ccx-go.exe"
cp "$DIST_DIR/ccx-go.exe" "$DESKTOP_DIR/build/windows/nsis/ccx-go.exe"

# Step 4: 构建桌面应用 + 安装包
echo ""
echo "[3/4] 构建桌面应用 ccx-desktop.exe..."
cd "$DESKTOP_DIR"
"$WAILS3" task build ARCH=amd64 PRODUCTION=true

echo ""
echo "[4/4] 打包 NSIS 安装程序..."
cd "$DESKTOP_DIR/build/windows/nsis"
# makensis 需要 Windows 格式路径
"$MAKENSIS" -DARG_WAILS_AMD64_BINARY="F:\\workspace\\ai-trun\\desktop\\bin\\ccx-desktop.exe" project.nsi

echo ""
echo "=========================================="
echo "  ✅ 打包完成!"
echo "=========================================="
echo ""
echo "产物列表:"
echo "  桌面应用:   $DESKTOP_DIR/bin/ccx-desktop.exe"
echo "  安装程序:   $DESKTOP_DIR/bin/ccx-desktop-amd64-installer.exe"
echo "  后端服务:   $DIST_DIR/ccx-go.exe"
echo ""
