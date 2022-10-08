package models

type Row struct {
	Amount0         string `json:"amount0"`
	Amount1         string `json:"amount1"`
	ContractAddress string `json:"contract_address"`
	EvtBlockNumber  int    `json:"evt_block_number"`
	EvtBlockTime    string `json:"evt_block_time"`
	EvtIndex        int    `json:"evt_index"`
	EvtTxHash       string `json:"evt_tx_hash"`
	Liquidity       string `json:"liquidity"`
	Recipient       string `json:"recipient"`
	Sender          string `json:"sender"`
	SqrtPriceX96    string `json:"sqrtPriceX96"`
	Tick            int    `json:"tick"`
}
