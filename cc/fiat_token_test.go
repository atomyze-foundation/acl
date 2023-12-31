package cc

import (
	"errors"
	"fmt"

	"github.com/atomyze-foundation/foundation/core/acl"
	"github.com/atomyze-foundation/foundation/core/types"
	"github.com/atomyze-foundation/foundation/core/types/big"
	"github.com/atomyze-foundation/foundation/token"
)

// FiatToken - base struct
type FiatToken struct {
	token.BaseToken
}

// NewFiatToken creates fiat token
func NewFiatToken(bt token.BaseToken) *FiatToken {
	return &FiatToken{bt}
}

// TxEmit - emits fiat token
func (mt *FiatToken) TxEmit(sender *types.Sender, address *types.Address, amount *big.Int) error {
	if !sender.Equal(mt.Issuer()) {
		return errors.New("unauthorized")
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount should be more than zero")
	}

	if err := mt.TokenBalanceAdd(address, amount, "txEmit"); err != nil {
		return err
	}
	return mt.EmissionAdd(amount)
}

func (mt *FiatToken) QueryGetRight(channel, chaincode, role, operation, address string) (bool, error) {
	stub := mt.GetStub()
	if stub == nil {
		return false, fmt.Errorf("getting stub failed, stub is nil")
	}

	params := []string{channel, chaincode, role, operation, address}
	haveRight, err := acl.GetAccountRight(stub, params)
	if err != nil {
		return false, err
	}
	if haveRight.HaveRight {
		return true, nil
	}

	return false, nil
}
