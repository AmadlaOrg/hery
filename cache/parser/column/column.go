package column

import "github.com/AmadlaOrg/hery/cache/database"

var (
	// Id returns simple database column structure for `Id`
	Id = database.Column{
		ColumnName: "Id",
		DataType:   "string",
		Default: `(
lower(hex(randomblob(4))
|| '-'
|| hex(randomblob(2))
|| '-4'
|| substr(hex(randomblob(2)), 2)
|| '-'
|| substr('89ab', abs(random() % 4) + 1, 1)
|| substr(hex(randomblob(2)), 2)
|| '-'
|| hex(randomblob(6)))
)`,
	}
	InsertDateTime = database.Column{
		ColumnName: "InsertDateTime",
		DataType:   "DATETIME",
	}
	UpdateDateTime = database.Column{
		ColumnName: "UpdateDateTime",
		DataType:   "DATETIME",
	}
)
