package deployed

// 设置viper默认值回调
var defaultValueFuncs []func()

// ViperAddDefaultFunc 增加设置viper默认值回调
func ViperAddDefaultFunc(f func()) {
	defaultValueFuncs = append(defaultValueFuncs, f)
}

// ViperInitDefault 初始化viper默认值回调
func ViperInitDefault() {
	for _, f := range defaultValueFuncs {
		f()
	}
}
