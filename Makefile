BINARY_NAME := ozen
# ソースファイル
SRC := ./main.go
# 出力ディレクトリ
DIST := dist

# プロダクションビルド用フラグ
# -s: シンボルテーブルの削除 (デバッガで使えなくなるがサイズが減る)
# -w: DWARFデバッグ情報の削除
LDFLAGS := -ldflags "-s -w"

# ターゲット定義（ファイル名との衝突を避ける）
.PHONY: all clean linux windows

all: clean linux windows

# Linux x64 (AMD64) ビルド
linux:
	@echo "Building for Linux (x64)..."
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY_NAME)-linux-amd64 $(SRC)
	@echo "Done: $(DIST)/$(BINARY_NAME)-linux-amd64"

# Windows x64 (AMD64) ビルド
windows:
	@echo "Building for Windows (x64)..."
	@mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY_NAME)-windows-amd64.exe $(SRC)
	@echo "Done: $(DIST)/$(BINARY_NAME)-windows-amd64.exe"

# 生成物の削除
clean:
	@echo "Cleaning..."
	rm -rf $(DIST)
