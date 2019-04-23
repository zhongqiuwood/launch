## GovStoreKey:"gov"
| key                                                   | value                                               | number(key)                                                    | value details | value size  | clean up                                   | 备注     |
| ---------------------------------                     | ---------------------------                         | --------------------                                           | -------       | --------    | --------                                   | -------- |
| newProposalID                                         | next proposalID to be generated                     | 1                                                              | NA            | < 1k        | NA                                         |          |
| proposals:\${ProposalID}                              | struct x/gov.TextProposal                           | 满足最少抵押额度和不满足最少抵押额度且仍处于抵押周期的提案数量 | NA            | > 1k<br>(max6k) | 抵押超时仍不满足最小抵押额度               |          |
| proposals:\${ProposalID}                              | struct x/gov.DexListProposal                        | 满足最少抵押额度和不满足最少抵押额度且仍处于抵押周期的提案数量 | NA            | < 1k        | 抵押超时仍不满足最小抵押额度               |          |
| proposals:\${ProposalID}                              | struct x/gov.ParameterProposal                      | 满足最少抵押额度和不满足最少抵押额度且仍处于抵押周期的提案数量 | NA            | < 1k        | 抵押超时仍不满足最小抵押额度               |          |
| votes:\${proposalID}:\${VoterAddr}                    | struct x/gov.Vote                                     | 所有投票的票数                                                 | NA            | < 1k        | 投票被统计后                               |          |
| deposits:\${proposalID}:\${DepositorAddr}             | struct x/gov.Deposit                               | 所有抵押的次数（包含发起提案时的抵押）                         | NA            | < 1k        | 抵押超时仍不满足最小抵押额度或投票被统计后 |          |
| activeProposalQueue\${Time}\${ProposalIDBE}           | proposalID of Proposal in Voting Period             | 处在投票阶段的提案                                             | NA            | < 1k        | 投票被统计后                               |          |
| inactiveProposalQueue\${Time}\${ProposalIDBE}         | proposalID of Proposal in Deposit Period            | 处在抵押阶段并且抵押额度不满足最少抵押额度的提案               | NA            | < 1k        | 进入投票阶段或抵押超时仍不满足最小抵押额度 |          |
| waitingProposalQueue\${BlockHeightBE}\${ProposalIDBE} | proposalID of Proposal to be excuted in BlockHeight | 投票通过并待处理的提案                                         | NA            | < 1k        | 处理后                                     | 见[^注]         |

## [^注]: 
* time为格式：`2006-01-02T15:04:05.000000000`的时间经过format后的字节数组；
* proposalIDBE为proposalID的大端编码；
* BlockHeightBE为BlockHeight的大端编码。e.g.:以waitingProposalQueue\${BlockHeightBE}\${ProposalIDBE}为例说明，
* waitingProposalQueue1001:1该KV对中key中100为BlockHeight（`为方面说明这里直接用100表示实际为100的大端编码`），
* 1为proposalID（同BlockHeight需要大端编码）；value中1即proposalID。
* 数据库中存在waitingProposalQueue1001:1，
* waitingProposalQueue1012:2，
* waitingProposalQueue1023:3三个键值对，
* 通过waitingProposalQueue101作为key获取value可以获取1和2两个value。