package form

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/xx-labs/sleeve/wallet"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"strconv"
)

// ValidateNodeID returns an error if the base 64 encoded string cannot be
// validated as a node ID. This functions adheres to the ValidateFunc type.
func ValidateNodeID(str string) error {
	if len(str) == 0 {
		return errors.New("Required.")
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	nodeID, err := id.Unmarshal(data)
	if err != nil {
		return err
	}

	if nodeID.GetType() != id.Node {
		return errors.Errorf("ID is of type %s; must be of type %s",
			nodeID.GetType(), id.Node)
	}

	return nil
}

// ValidateXXNetworkAddress returns an error if the xx network address is
// invalid. This functions adheres to the ValidateFunc type.
func ValidateXXNetworkAddress(str string) error {
	if len(str) == 0 {
		return errors.New("Required.")
	}

	ok, err := wallet.ValidateXXNetworkAddress(str)
	if !ok || err != nil {
		return errors.Errorf("Invalid wallet address: %s", err.Error())
	}

	return nil
}

// ValidateMultiplier returns an error if the xx network address is
// invalid. This functions adheres to the ValidateFunc type.
func ValidateMultiplier(max float32) ValidateFunc {
	return func(str string) error {
		if len(str) == 0 {
			return errors.New("Required.")
		}

		f64, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}

		f := float32(f64)

		if f < 0 || f > max {
			return errors.Errorf("value must be between %.3f and %.3f", 0.0, max)
		}

		return nil
	}
}

// ValidateFilePath returns an error if the file path is invalid. This functions
// adheres to the ValidateFunc type.
func ValidateFilePath(str string) error {
	if len(str) == 0 {
		return errors.New("Required.")
	}

	if _, err := utils.ReadFile(str); err != nil {
		return err
	}

	return nil
}
