package up

type ListOption struct {
	name, value string
}

func NewListOption(name, value string) ListOption {
	return ListOption{
		name:  name,
		value: value,
	}
}
