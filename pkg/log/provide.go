package log

func ProvideLogger() *LogCustom {
	return NewLogger()
}
