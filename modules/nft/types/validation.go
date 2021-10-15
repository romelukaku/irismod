package types

import (
	"fmt"
	"regexp"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/token/types"
)

const (
	DoNotModify = "[do-not-modify]"
	MinClassLen = 3
	MaxClassLen = 64

	MaxTokenURILen = 256

	ReservedPeg  = "peg"
	ReservedIBC  = "ibc"
	ReservedHTLT = "htlt"
	ReservedTIBC = "tibc"
)

var (
	// IsAlphaNumeric only accepts [a-z0-9]
	IsAlphaNumeric = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	// IsBeginWithAlpha only begin with [a-z]
	IsBeginWithAlpha = regexp.MustCompile(`^[a-z].*`).MatchString

	keywords          = strings.Join([]string{ReservedPeg, ReservedIBC, ReservedHTLT, ReservedTIBC}, "|")
	regexpKeywordsFmt = fmt.Sprintf("^(%s).*", keywords)
	regexpKeyword     = regexp.MustCompile(regexpKeywordsFmt).MatchString
)

// ValidateClassID verifies whether the  parameters are legal
func ValidateClassID(classID string) error {
	if len(classID) < MinClassLen || len(classID) > MaxClassLen {
		return sdkerrors.Wrapf(ErrInvalidClass, "the length of class(%s) only accepts value [%d, %d]", classID, MinClassLen, MaxClassLen)
	}
	if !IsBeginWithAlpha(classID) || !IsAlphaNumeric(classID) {
		return sdkerrors.Wrapf(ErrInvalidClass, "the class(%s) only accepts alphanumeric characters, and begin with an english letter", classID)
	}
	return ValidateKeywords(classID)
}

// ValidateTokenID verify that the tokenID is legal
func ValidateTokenID(tokenID string) error {
	if len(tokenID) < MinClassLen || len(tokenID) > MaxClassLen {
		return sdkerrors.Wrapf(ErrInvalidTokenID, "the length of nft id(%s) only accepts value [%d, %d]", tokenID, MinClassLen, MaxClassLen)
	}
	if !IsBeginWithAlpha(tokenID) || !IsAlphaNumeric(tokenID) {
		return sdkerrors.Wrapf(ErrInvalidTokenID, "nft id(%s) only accepts alphanumeric characters, and begin with an english letter", tokenID)
	}
	return nil
}

// ValidateTokenURI verify that the tokenURI is legal
func ValidateTokenURI(tokenURI string) error {
	if len(tokenURI) > MaxTokenURILen {
		return sdkerrors.Wrapf(ErrInvalidTokenURI, "the length of nft uri(%s) only accepts value [0, %d]", tokenURI, MaxTokenURILen)
	}
	return nil
}

// Modified returns whether the field is modified
func Modified(target string) bool {
	return target != types.DoNotModify
}

// ValidateKeywords checks if the given classId begins with `ClassKeywords`
func ValidateKeywords(classId string) error {
	if regexpKeyword(classId) {
		return sdkerrors.Wrapf(ErrInvalidClass, "invalid classId: %s, can not begin with keyword: (%s)", classId, keywords)
	}
	return nil
}
