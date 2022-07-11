package form

import (
	"github.com/pkg/errors"
	"github.com/xx-labs/sleeve/wallet"
	"gitlab.com/xx_network/primitives/id/idf"
	"gitlab.com/xx_network/primitives/utils"
	"net/mail"
	"strconv"
)

// ValidateXXNetworkAddress returns an error if the xx network address is
// invalid. This function adheres to the ValidateFunc type.
func ValidateXXNetworkAddress(str string) (interface{}, error) {
	if len(str) == 0 {
		return nil, errors.New("Required.")
	}

	ok, err := wallet.ValidateXXNetworkAddress(str)
	if !ok || err != nil {
		return nil, errors.Errorf("Invalid wallet address: %s", err.Error())
	}

	return str, nil
}

// ValidateEmail returns an error if the email is invalid. This function adheres
// to the ValidateFunc type.
func ValidateEmail(str string) (interface{}, error) {
	if len(str) == 0 {
		return "", nil
	}

	_, err := mail.ParseAddress(str)

	return str, err
}

// ValidateMultiplier returns an error if the xx network address is
// invalid.
func ValidateMultiplier(max uint64) ValidateFunc {
	return func(str string) (interface{}, error) {
		if len(str) == 0 {
			return nil, errors.New("Required.")
		}

		u64, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}

		if u64 < 0 || u64 > max {
			return nil, errors.Errorf(
				"value must be between %d and %d", 0, max)
		}

		return u64, nil
	}
}

// ValidateFilePath returns an error if the file path is invalid. This function
// adheres to the ValidateFunc type.
func ValidateFilePath(str string) (interface{}, error) {
	if len(str) == 0 {
		return nil, errors.New("Required.")
	}

	file, err := utils.ReadFile(str)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ValidateIdfPath returns an error if the IDF file path is invalid. This
// function adheres to the ValidateFunc type.
func ValidateIdfPath(str string) (interface{}, error) {
	if len(str) == 0 {
		return nil, errors.New("Required.")
	}

	_, nid, err := idf.UnloadIDF(str)
	if err != nil {
		return nil, err
	}

	return nid.HexEncode(), nil
}

// ValidateCheckbox returns an error if the checkbox is not checked. This
// function adheres to the ValidateFunc type.
func ValidateCheckbox(str string) (interface{}, error) {
	if len(str) == 0 {
		return nil, errors.New("Required.")
	}

	return true, nil
}
