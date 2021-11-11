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
                "ethAddr": "0x9bb963c0752431df1eb1810bbe0c5405daef8015",
                "djtxAddr": "X-dijets1v6ztr2wwj633ccqxjch3hzqtprlx2h55sdna9r",
                "initialAmount": 500000000000000,
                "unlockSchedule": [
                    {
                        "amount": 4375000000000000,
                        "locktime": 1634565635
                    }
                ]
            },
            {
                "ethAddr": "0x9bb963c0752431df1eb1810bbe0c5405daef8015",
                "djtxAddr": "X-dijets1v6ztr2wwj633ccqxjch3hzqtprlx2h55sdna9r",
                "initialAmount": 500000000000000,
                "unlockSchedule": [
                    {
                        "amount": 4375000000000000,
                        "locktime": 1635256835
                    }
                ]
            },
            {
                "ethAddr": "0x9bb963c0752431df1eb1810bbe0c5405daef8015",
                "djtxAddr": "X-dijets1v6ztr2wwj633ccqxjch3hzqtprlx2h55sdna9r",
                "initialAmount": 500000000000000,
                "unlockSchedule": [
                    {
                        "amount": 4375000000000000,
                        "locktime": 1635692435
                    }
                ]
            },
            {
                "ethAddr": "0x9bb963c0752431df1eb1810bbe0c5405daef8015",
				"djtxAddr": "X-dijets1v6ztr2wwj633ccqxjch3hzqtprlx2h55sdna9r",
                "initialAmount": 500000000000000,
                "unlockSchedule": [
                    {
                        "amount": 4375000000000000,
                        "locktime": 1636556435
                    }
                ]
            }
        ],
        "startTime": 1599696000,
        "initialStakeDuration": 31536000,
        "initialStakeDurationOffset": 5400,
        "initialStakedFunds": [
            "X-dijets1v6ztr2wwj633ccqxjch3hzqtprlx2h55sdna9r"
        ],
        "initialStakers": [
            {
                "nodeID": "NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg",
                "rewardAddress": "X-dijets18rr79nevvpgda5llysysckpmnve7xyadt9jmsv",
                "delegationFee": 1000000
            },
            {
                "nodeID": "NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ",
                "rewardAddress": "X-dijets18rr79nevvpgda5llysysckpmnve7xyadt9jmsv",
                "delegationFee": 500000
            },
            {
                "nodeID": "NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN",
                "rewardAddress": "X-dijets18rr79nevvpgda5llysysckpmnve7xyadt9jmsv",
                "delegationFee": 250000
            },
            {
                "nodeID": "NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu",
                "rewardAddress": "X-dijets18rr79nevvpgda5llysysckpmnve7xyadt9jmsv",
                "delegationFee": 125000
            },
            {
                "nodeID": "NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5",
                "rewardAddress": "X-dijets18rr79nevvpgda5llysysckpmnve7xyadt9jmsv",
                "delegationFee": 62500
            }
        ],
        "cChainGenesis": "{\"config\":{\"chainId\":43114,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"apricotPhase1BlockTimestamp\":0,\"apricotPhase2BlockTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC\":{\"balance\":\"0x295BE96E64066972000000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}",
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
