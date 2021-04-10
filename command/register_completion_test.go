package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCompletion(t *testing.T) {
	skipShort(t)
	TestRegister(t)
	for {
		got := RegisterResponse{}
		err := PaymentTerminal.Completion(&got)
		if err != nil {
			assert.NoError(t, err)
			break
		}
		if got.Transaction != nil && got.Transaction.Result != Result_Pending {
			if got.Transaction.Result == Result_Success {
				// TODO assert result values
				_ = 0
			}
			break
		}
		// assert.EqualValues(t, want, got)
	}
}

// func TestRegisterProcess(t *testing.T) {
// 	want := RegisterResponse{}
// 	testBytes, err := util.Load("testdata/.hex")
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	c := Command{}
// 	c.Unmarshal(&testBytes)
// 	got := RegisterResponse{}
// 	got.Process(&c)
// 	assert.EqualValues(t, want, got)
// }

// func TestRegisterProcess2(t *testing.T) {
// 	testBytes, err := util.Load("testdata/.hex")
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	c := Command{}
// 	err = c.Unmarshal(&testBytes)
// 	assert.Error(t, err)
// }

// func TestRegisterProcess3(t *testing.T) {
// 	testBytes, err := util.Load("testdata/.hex")
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	want := Command{
// 		CtrlField: instr.CtrlField{
// 			Class: 0x06,
// 			Instr: 0xD3,
// 			Length: blen.Length{
// 				Kind:  blen.BINARY,
// 				Size:  3,
// 				Value: 1069,
// 			},
// 			RawDataLength: 0,
// 		},
// 		Data: apdu.DataUnit{
// 			Data:    []byte{},
// 			BMPOBJs: []bmp.OBJ{},
// 			TLVContainer: tlv.Container{
// 				Objects: []tlv.DataObject{
// 					{[]byte{0xD3}, []byte{0x00}},
// 				},
// 			},
// 		},
// 	}
// 	got := Command{}
// 	if err = got.Unmarshal(&testBytes); assert.NoError(t, err) {
// 		if assert.EqualValues(t, want.CtrlField, got.CtrlField) &&
// 			assert.Equal(t, 2, len(got.Data.TLVContainer.Objects)) {
// 			got2 := RegisterResponse{}
// 			got2.Process(&got)
// 			assert.EqualValues(t, 3, got2.Transaction.Data.PrintOut.Type)
// 			assert.EqualValues(t, 1057, len(got2.Transaction.Data.PrintOut.Text))
// 		}
// 	}

// }
