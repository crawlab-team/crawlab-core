package utils

func IsMaster() bool {
	return EnvIsTrue("node.master", false)
}
