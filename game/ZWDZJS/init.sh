#!/bin/bash

# 项目名称
APP_NAME="Star Dust"
EXECUTABLE_NAME="star_dust"
ICONSET_NAME="icon.iconset"
ICNS_NAME="icon.icns"

# 创建目录结构
mkdir -p "$APP_NAME.app/Contents/MacOS"
mkdir -p "$APP_NAME.app/Contents/Resources"

# 创建 Info.plist 文件
cat > "$APP_NAME.app/Contents/Info.plist" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$EXECUTABLE_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>com.yourcompany.$EXECUTABLE_NAME</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.9</string>
    <key>CFBundleIconFile</key>
    <string>$ICNS_NAME</string>
</dict>
</plist>
EOL

# 编译 Go 程序
go build -o $EXECUTABLE_NAME main.go

# 生成 .icns 文件
if [ -d "$ICONSET_NAME" ]; then
    iconutil -c icns $ICONSET_NAME
    mv "icon.icns" "$APP_NAME.app/Contents/Resources/$ICNS_NAME"
fi

# 移动可执行文件
mv $EXECUTABLE_NAME "$APP_NAME.app/Contents/MacOS/"

echo "$APP_NAME.app is ready."
