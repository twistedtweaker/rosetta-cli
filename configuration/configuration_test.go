package configuration

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/coinbase/rosetta-cli/internal/utils"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

var (
	whackyConfig = &Configuration{
		Construction: &ConstructionConfiguration{
			Network: &types.NetworkIdentifier{
				Blockchain: "sweet",
				Network:    "sweeter",
			},
			OnlineURL:  "http://hasudhasjkdk",
			OfflineURL: "https://ashdjaksdkjshdk",
			Currency: &types.Currency{
				Symbol:   "FIRE",
				Decimals: 100,
			},
			MinimumBalance:   "1002",
			MaximumFee:       "1",
			CurveType:        types.Edwards25519,
			AccountingModel:  UtxoModel,
			TransferScenario: DefaultTransferScenario,
		},
		Data: &DataConfiguration{
			OnlineURL:                         "https://asjdlkasjdklajsdlkj",
			BlockConcurrency:                  12,
			TransactionConcurrency:            2,
			ActiveReconciliationConcurrency:   100,
			InactiveReconciliationConcurrency: 2938,
			InactiveReconciliationFrequency:   3,
		},
	}
	invalidNetwork = &Configuration{
		Construction: &ConstructionConfiguration{
			Network: &types.NetworkIdentifier{
				Blockchain: "?",
			},
		},
	}
	invalidCurrency = &Configuration{
		Construction: &ConstructionConfiguration{
			Currency: &types.Currency{
				Decimals: 12,
			},
		},
	}
	invalidCurve = &Configuration{
		Construction: &ConstructionConfiguration{
			CurveType: "hello",
		},
	}
	invalidAccountingModel = &Configuration{
		Construction: &ConstructionConfiguration{
			AccountingModel: "hello",
		},
	}
	invalidMinimumBalance = &Configuration{
		Construction: &ConstructionConfiguration{
			MinimumBalance: "-1000",
		},
	}
	invalidMaximumFee = &Configuration{
		Construction: &ConstructionConfiguration{
			MaximumFee: "hello",
		},
	}
)

func TestLoadConfiguration(t *testing.T) {
	var tests = map[string]struct {
		provided *Configuration
		expected *Configuration

		err bool
	}{
		"nothing provided": {
			provided: &Configuration{},
			expected: DefaultConfiguration(),
		},
		"no overwrite": {
			provided: whackyConfig,
			expected: whackyConfig,
		},
		"invalid network": {
			provided: invalidNetwork,
			err:      true,
		},
		"invalid currency": {
			provided: invalidCurrency,
			err:      true,
		},
		"invalid curve type": {
			provided: invalidCurve,
			err:      true,
		},
		"invalid accounting model": {
			provided: invalidAccountingModel,
			err:      true,
		},
		"invalid minimum balance": {
			provided: invalidMinimumBalance,
			err:      true,
		},
		"invalid maximum fee": {
			provided: invalidMaximumFee,
			err:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Write configuration file to tempdir
			tmpfile, err := ioutil.TempFile("", "test.json")
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			err = utils.SerializeAndWrite(tmpfile.Name(), test.provided)
			assert.NoError(t, err)

			// Check if expected fields populated
			config, err := LoadConfiguration(tmpfile.Name())
			if test.err {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, config)
			}
			assert.NoError(t, tmpfile.Close())
		})
	}
}
