package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/abci/types"
)

func refund(ctx sdk.Context, gm governanceMapper, int16 proposalId, Proposal proposal) {
	newDeposits := new([]Deposit)
	for _, d := range proposal.Deposits {
		newDeposits = append(newDeposits, Deposit{d.Depositer, 0})
		gm.AddCoins(ctx, d.Depositer, d.Amount)
	}
	proposal.Deposits = newDeposits
	gm.SetProposal(proposalId, proposal)
}

func checkProposal(ctx sdk.Context, gm governanceMapper) {
	proposalId := gm.PeekProposalQueue(ctx)
	if proposalId == nil {
		return
	}
	proposal := gm.GetProposal(ctx, proposalId)

	// Urgent proposal accepted
	if proposal.Votes.YesVotes/proposal.InitTotalVotingPower >= 2/3 {
		gm.PopProposalQueue(ctx)
		refund(ctx, gm, proposalId, proposal)
		return checkProposal()
	}

	currentBlock := ctx.BlockHeight()
	proceduer := gm.GetProdeuer(proposal.InitProceduer)
	// Proposal reached the end of the voting period
	if currentBlock == proposal.VotingStartBlock+proceduer.VotingPeriod {
		gm.PopProposalQueue(ctx)
		activeProcedure := gm.GetParam(ctx, ActiveProcedure)
		currentBondedValidators := gm.GetParam(ctx, CurrentBondedValidators)

		// Slash validators if not voted
		for _, v := range currentBondedValidators {
			validatorGovInfo := gm.GetValidatorGovInfo(ctx, proposalId, v.Address)
			if validatorGovInfo.InitVotingPower != nil {
				validatorOption := gm.GetOption(ctx, proposalId, v.Address)
				if validatorOption == nil {
					// slash
				}
			}
		}

		//Proposal was accepted
		noAbstainTotal := proposal.Votes.YesVotes + proposal.Votes.NoVotes + proposal.Votes.NoWithVetoVotes
		if proposal.Votes.YesVotes/noAbstainTotal > 0.5 && proposal.Votes.NoWithVetoVotes/noAbstainTotal < 1/3 {
			refund(ctx, gm, proposalId, proposal)
			return checkProposal()
		}
	}
}

func NewBeginBlocker(gm governanceMapper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		checkProposal(ctx, gm)
		return abci.ResponseBeginBlock{}
	}
}
