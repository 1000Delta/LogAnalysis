package models

// Container 表示 Docker 容器信息
type Container struct {
	Name string
	Image string
	// Ports 表示端口映射，"19200:9200" -> 19200 => 9200
	Ports map[int]int
	// Volumes 表示数据卷绑定，
	Volumes map[*Volume]string
}

type volumeDriver int

// 数据卷类型
const (
	VolumeDriverLocal volumeDriver = iota
)

// Volume 表示数据卷信息
type Volume struct {
	Name   string
	Driver string
}

// NewVolume 新建数据卷信息对象
func NewVolume(name string, driver volumeDriver) *Volume {
	vol := &Volume{
		Name: name,
	}
	switch driver {
	case VolumeDriverLocal:
		vol.Driver = "local"
	default:
		return nil
	}
	return vol
}
