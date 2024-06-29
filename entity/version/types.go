package version

const (
	Match        = `.+@v\d+\.\d+\.\d+`
	Format       = `^v(\d+)(\.\d+)?(\.\d+)?$`
	FormatForDir = `(.+)@v(\d+\.\d+\.\d+)`
)
