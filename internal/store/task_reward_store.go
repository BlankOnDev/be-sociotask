package store

type JenisCategory string

const (
	CryptoUsdt1 JenisCategory = "crypto_usdt_1"
	CryptoUsdt2 JenisCategory = "crypto_usdt_2"
	CryptoUsdt3 JenisCategory = "crypto_usdt_3"
)

type RewardTask struct {
	ID         int
	RewardType JenisCategory
	RewardName string
}
