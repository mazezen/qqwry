#!/usr/bin/env bash

# =============================================
# 自动下载最新纯真IP库 qqwry.dat
# 来源：https://github.com/metowolf/qqwry.dat
# =============================================

set -e

# 下载地址（GitHub latest release 的直链）
URL="https://github.com/metowolf/qqwry.dat/releases/latest/download/qqwry.dat"

# 保存的文件名
OUTPUT_FILE="qqwry.dat"

# 临时文件（避免下载中断覆盖旧文件）
TEMP_FILE="qqwry.dat.tmp"

echo "正在下载最新 qqwry.dat ..."

# 使用 curl 下载（推荐），如果没有 curl 可用 wget
if command -v curl >/dev/null 2>&1; then
    curl -L -o "$TEMP_FILE" --fail "$URL" || {
        echo "下载失败，请检查网络或 GitHub 是否可访问"
        exit 1
    }
else
    if command -v wget >/dev/null 2>&1; then
        wget -O "$TEMP_FILE" "$URL" || {
            echo "下载失败，请检查网络或 GitHub 是否可访问"
            exit 1
        }
    else
        echo "错误：需要安装 curl 或 wget"
        exit 1
    fi
fi

# 检查文件是否正常（大小应该在 25MB 左右或更大）
FILE_SIZE=$(stat -c %s "$TEMP_FILE" 2>/dev/null || stat -f %z "$TEMP_FILE" 2>/dev/null || ls -l "$TEMP_FILE" | awk '{print $5}')
MIN_SIZE=$((20 * 1024 * 1024))  # 至少 20MB

if [ "$FILE_SIZE" -lt "$MIN_SIZE" ]; then
    echo "下载的文件太小（${FILE_SIZE} bytes），可能是损坏或下载失败"
    rm -f "$TEMP_FILE"
    exit 1
fi

# 替换原文件
mv "$TEMP_FILE" "$OUTPUT_FILE"
echo "更新成功！"
echo "文件: $OUTPUT_FILE"
echo "大小: $((FILE_SIZE / 1024 / 1024)) MB"
echo "完成时间: $(date '+%Y-%m-%d %H:%M:%S')"

# 可选：显示文件修改时间
ls -l "$OUTPUT_FILE"