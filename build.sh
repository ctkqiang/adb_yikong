#!/bin/bash

echo "开始构建 ADB易控..."

rm -rf build/
mkdir -p build

VERSION=$(grep "Version" internal/constants/constants.go | cut -d'"' -f2)
APP_NAME="ADB易控"

echo "应用名称: $APP_NAME"
echo "版本: $VERSION"

echo "正在构建应用..."
go build -o "build/$APP_NAME" main.go

if [ $? -eq 0 ]; then
    echo "构建成功！"
    echo "应用已生成: build/$APP_NAME"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "正在创建macOS应用包..."
        
        APP_DIR="build/$APP_NAME.app"
        CONTENTS_DIR="$APP_DIR/Contents"
        
        mkdir -p "$CONTENTS_DIR/MacOS"
        mkdir -p "$CONTENTS_DIR/Resources"
        
        mv "build/$APP_NAME" "$CONTENTS_DIR/MacOS/"
        
        cat > "$CONTENTS_DIR/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>
    <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>com.ctkqiang.adbcontroller</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleExecutable</key>
    <string>$APP_NAME</string>
    <key>LSUIElement</key>
    <false/>
</dict>
</plist>
EOF
        
        echo "macOS应用包创建完成: $APP_DIR"
    fi
    
else
    echo "构建失败！"
    exit 1
fi

echo "构建完成！"