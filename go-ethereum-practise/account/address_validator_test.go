package account

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestAddressValidate(t *testing.T) {
	addr1 := "0x323b5d4c32345ced77393b3530b1eed0f346429d"
	result, err := AddressValidate(common.Client(), addr1)
	if err != nil {
		if err == common.NotValidAddress {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Log(result)
	}

	addr2 := "0xZYXb5d4c32345ced77393b3530b1eed0f346429d"
	result, err = AddressValidate(common.Client(), addr2)
	if err != nil {

		if err == common.NotValidAddress {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Log(result)
	}

	addr3 := "0xe41d2489571d322189246dafa5ebde1f4699f498"
	result, err = AddressValidate(common.Client(), addr3)
	if err != nil {

		if err == common.NotValidAddress {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Log(result)
	}

	addr4 := "0x8e215d06ea7ec1fdb4fc5fd21768f4b34ee92ef4"
	result, err = AddressValidate(common.Client(), addr4)
	if err != nil {

		if err == common.NotValidAddress {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Log(result)
	}
}
