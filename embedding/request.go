package embedding

func NewRequest(inputs []string, option Option) Request {
	return Request{
		Inputs: inputs,
		Option: option,
	}
}

type Request struct {
	Inputs []string
	Option Option
}
