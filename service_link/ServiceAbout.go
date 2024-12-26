package service_link

// Linker 定义了一个接口，要求实现创建链接和获取类型名称的方法
type Linker interface {
	CreateLink() (any, error)
	GetName() string
}
