package http

import "fmt"

func getClientLogicName(dir string, packageName string) string {
	return dir + "/" + packageName + "_logic.go"
}

// client: 函数
func getClientFuncName(dir string, packageName string) string {
	return dir + "/" + packageName + ".go"
}

// client: 结构体定义
func getClientTypeName(dir string, packageName string) string {
	return dir + "/" + packageName + "_type.go"
}

// server: main包前缀
func getServerPrefixMain(dir string, goModLastName string) string {
	return fmt.Sprintf("%s/api/cmd/%s", dir, goModLastName)
}

// server: main包全路径
func getServerMainName(dir string, gomod string) string {
	return dir + "/" + gomod + ".go"
}

// server:type目录前缀
func getServerTypePrefix(dir string, packageName string) string {
	return fmt.Sprintf("%s/api/internal/types/%s/", dir, packageName)
}

// server:type全路径
func getServerTypeName(dir string, packageName string) string {
	return fmt.Sprintf("%s/%s_type.go", dir, packageName)
}

// server:svc目录前缀
func getServerSvcPrefix(dir string) string {
	return fmt.Sprintf("%s/api/internal/svc/", dir)
}

// server:svc全路径
func getServerSvcName(dir string) string {
	return fmt.Sprintf("%s/servicecontext.go", dir)
}

// server:svc目录前缀
func getServerConfigPrefix(dir string) string {
	return fmt.Sprintf("%s/api/internal/config/", dir)
}

// server:svc全路径
func getServerConfigName(dir string) string {
	return fmt.Sprintf("%s/config.go", dir)
}

// server:logic
func getLogicPrefix(dir string, packageName string) string {
	return fmt.Sprintf("%s/api/internal/logic/%s/", dir, packageName)
}

func getLogicName(dir string, handler string) string {
	return fmt.Sprintf("%s/%s_logic.go", dir, handler)
}

// server:handler
func getHandlerPrefix(dir string, packageName string) string {
	return fmt.Sprintf("%s/api/internal/handler/%s/", dir, packageName)
}

// server:handler 全路径
func getHandlerName(dir string, handler string) string {
	return fmt.Sprintf("%s/%s_handler.go", dir, handler)
}

// server:handler
func getRoutesPrefix(dir string) string {
	return fmt.Sprintf("%s/api/internal/handler/", dir)
}

// server: routes 全路径
func getRoutesName(dir string) string {
	return fmt.Sprintf("%s/routes.go", dir)
}

// server: etc 目录前缀
func getEtcPrefix(dir string) string {
	return fmt.Sprintf("%s/api/etc/", dir)
}

// server: etc 目录前缀
func getEtcName(dir string, goModLastName string) string {
	return fmt.Sprintf("%s/%s.yaml", dir, goModLastName)
}
