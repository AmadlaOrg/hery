package version

const (
	Match        = `.+@(v\d+\.\d+\.\d+(-(alpha|beta|rc)\.\d+)?|v\d+\.\d+\.\d+-\d{14}-[a-f0-9]{12})$`
	Format       = `^v(\d+)(\.\d+)?(\.\d+)?(-(alpha|beta|rc)\.\d+)?$`
	FormatForDir = `(.+)@((v\d+\.\d+\.\d+(-(alpha|beta|rc)\.\d+)?|v\d+\.\d+\.\d+-\d{14}-[a-f0-9]{12}))$`
)
