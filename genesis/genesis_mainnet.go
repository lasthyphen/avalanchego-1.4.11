// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	"github.com/ava-labs/avalanchego/utils/units"
)

var (
	mainnetGenesisConfigJSON = `{
		"networkID": 1,
		"allocations": [
			{
				"ethAddr": "0x86c556954057da5a98ac9e5c5610fbac49e038d5",
				"djtxAddr": "X-djtx1gzzt7eaxl3qrdj26dw9lnkvpnxrzrl3ermgjne",
				"initialAmount": 240000000000,
				"unlockSchedule": [
					{
						"amount": 360000000000,
						"locktime": 1607472000
					},
					{
						"amount": 360000000000,
						"locktime": 1615248000
					},
					{
						"amount": 360000000000,
						"locktime": 1623024000
					},
					{
						"amount": 360000000000,
						"locktime": 1630800000
					},
					{
						"amount": 360000000000,
						"locktime": 1638576000
					},
					{
						"amount": 360000000000,
						"locktime": 1646352000
					}
				]
			},
			{
				"ethAddr": "0x3ba84aa8e5613bcfcbfe6c0736c7527b08ab669c",
				"djtxAddr": "X-djtx1yktwy8clp5jdl23rns0w0d7d2nggsp5mwcgtnj",
				"initialAmount": 10000000000000000,
				"unlockSchedule": [
					{
						"amount": 10000000000000000,
						"locktime": 1633824000
					}
				]
			}
		],
		"startTime": 1599696000,
		"initialStakeDuration": 31536000,
		"initialStakeDurationOffset": 5400,
		"initialStakedFunds": [
			"X-djtx1gzzt7eaxl3qrdj26dw9lnkvpnxrzrl3ermgjne"
		],
		"initialStakers": [
			{
				"nodeID": "NodeID-2MRhHKpa9SLZgRSkqAdcSrMGfKjvLuAdv",
				"rewardAddress": "X-djtx1gzzt7eaxl3qrdj26dw9lnkvpnxrzrl3ermgjne",
				"delegationFee": 1000000
			}
		],
		"cChainGenesis": "{\"config\":{\"chainId\":43112,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"apricotPhase1BlockTimestamp\":0,\"apricotPhase2BlockTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC\":{\"balance\":\"0x295BE96E64066972000000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}",
		"message": "{{ fun_quote }}"
	}`

	// MainnetParams are the params used for mainnet
	MainnetParams = Params{
		TxFee:                units.MilliDjtx,
		CreationTxFee:        10 * units.MilliDjtx,
		UptimeRequirement:    .6, // 60%
		MinValidatorStake:    2 * units.KiloDjtx,
		MaxValidatorStake:    3 * units.MegaDjtx,
		MinDelegatorStake:    25 * units.Djtx,
		MinDelegationFee:     20000, // 2%
		MinStakeDuration:     2 * 7 * 24 * time.Hour,
		MaxStakeDuration:     365 * 24 * time.Hour,
		StakeMintingPeriod:   365 * 24 * time.Hour,
		EpochFirstTransition: time.Unix(1607626800, 0),
		EpochDuration:        6 * time.Hour,
	}
)
