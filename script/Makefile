# 变量定义
ADMIN_APP := mxclub-admin
MINI_APP := mxclub-mini
SRC_DIR := ../apps/
BUILD_DIR := ./bin

.PHONY: all build-admin build-mini test-admin test-mini clean run-admin run-mini vet-admin vet-mini

# 默认目标
all: build-admin build-mini

# 构建 admin 应用
build-admin:
	@echo "Building admin application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(ADMIN_APP) $(SRC_DIR)/$(ADMIN_APP)

# 构建 mini 应用
build-mini:
	@echo "Building mini application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(MINI_APP) $(SRC_DIR)/$(MINI_APP)

# 运行 admin 应用
run-admin: build-admin
	@echo "Running admin application..."
	@$(BUILD_DIR)/$(ADMIN_APP)

# 运行 mini 应用
run-mini: build-mini
	@echo "Running mini application..."
	@$(BUILD_DIR)/$(MINI_APP)

# 测试
test:
	@echo "Running admin tests..."
	go list ../... | grep -v config | xargs go test -run .

# 静态分析
vet:
	@echo "Running go vet on admin..."
	go list ../... | xargs go vet -composites=false

# 清理构建文件
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)