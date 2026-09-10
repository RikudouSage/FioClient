package db

type loadOptions struct {
	where            []string
	limit            uint
	orderByField     string
	orderByDirection string
	bindValues       []any
}

type LoadOption func(options *loadOptions)

func WithWhere(where string, bind ...any) LoadOption {
	return func(options *loadOptions) {
		options.where = append(options.where, where)
		options.bindValues = append(options.bindValues, bind...)
	}
}

func WithLimit[TInt ~int | ~uint](limit TInt) LoadOption {
	return func(options *loadOptions) {
		options.limit = uint(limit)
	}
}

func WithOrderBy(column, direction string) LoadOption {
	return func(options *loadOptions) {
		options.orderByField = column
		options.orderByDirection = direction
	}
}

func WithAccountNumber(accountNumber string) LoadOption {
	return WithWhere("account_number = ?", accountNumber)
}
