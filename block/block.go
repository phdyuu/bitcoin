package block

type Block struct {
	TimeStamp		int64
	Data 			string
	PrevHash		[]byte
	Hash 			[]byte
	Nonce			int64
}
