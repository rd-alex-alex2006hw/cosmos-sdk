package gov

import crypto "github.com/tendermint/go-crypto"

type Vote struct {
	ProposalID      int64         `json:"proposal_id"`       //  proposalID of the proposal
	Option          string        `json:"option"`            //  option from OptionSet chosen by the voter
	ValidatorPubKey crypto.PubKey `json:"validator_pub_key"` //  PubKey of the validator voter wants to tie its vote to
}

type Votes struct {
	YesVotes        int64
	NoVotes         int64
	NoWithVetoVotes int64
	AbstainVotes    int64
}

//-----------------------------------------------------------

// Proposal
type Proposal struct {
	Title        string    `json:"title"`         //  Title of the proposal
	Description  string    `json:"description"`   //  Description of the proposal
	ProposalType string    `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	TotalDeposit sdk.Coins `json:"total_deposit"` //  Current deposit on this proposal. Initial value is set at InitialDeposit

	Deposits []Deposit `json:"deposits"` //  Current deposit on this proposal. Initial value is set at InitialDeposit

	SubmitBlock          int64 `json:"submit_block"`            //  Height of the block where TxGovSubmitProposal was included
	VotingStartBlock     int64 `json:"voting_start_block"`      //  Height of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	InitTotalVotingPower int64 `json:"init_total_voting_power"` //  Total voting power when proposal enters voting period (default 0)
	InitProcedure        int16 `json:"init_procedure"`          //  Procedure number of the active procedure when proposal enters voting period (default -1)

	Votes Votes
}

// Procedure
type Procedure struct {
	VotingPeriod      int64             `json:"voting_period"`      //  Length of the voting period. Initial value: 2 weeks
	MinDeposit        int64             `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	OptionSet         []string          `json:"option_set"`         //  Options available to voters. {Yes, No, NoWithVeto, Abstain}
	ProposalTypes     []string          `json:"proposal_type"`      //  Types available to submitters. {PlainTextProposal, SoftwareUpgradeProposal}
	Threshold         rational.Rational `json:"threshold"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto              rational.Rational `json:"veto"`               //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	MaxDepositPeriod  int64             `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	GovernancePenalty int64             `json:"governance_penalty"` //  Penalty if validator does not vote
	IsActive          bool              `json:"is_active"`          //  If true, procedure is active. Only one procedure can have isActive true.
}

// Deposit
type Deposit struct {
	Depositer crypto.address `json:"depositer"` //  Address of the depositer
	Amount    sdk.Coins      `json:"amount"`    //  Deposit amount
}

type ValidatorGovInfo struct {
	ProposalId      int64 //  Id of the Proposal this validator
	ValidatorAddr   int64 //  Address of the validator
	InitVotingPower int64 //  Voting power of validator when proposal enters voting period
	Minus           int64 //  Minus of validator, used to compute validator's voting power
}
