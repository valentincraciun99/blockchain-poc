package blockchain

type TxOutput struct {
	Value  int
	PubKey string
}

type TxOutputs struct {
	Outputs []TxOutput
}

type TxInput struct {
	ID        []byte
	OutIndex  int
	Signature string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
