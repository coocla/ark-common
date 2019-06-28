package tool_test

import (
	"ark-common/utils/tool"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDCP(t *testing.T) {
	Convey("测试加密与解密", t, func() {
		row := "abc"
		Convey("加密字符串", func() {
			So(tool.ECP(row), ShouldNotEqual, row)
		})

		Convey("解密字符串", func() {
			So(tool.DCP(tool.ECP(row)), ShouldEqual, row)
		})
	})
}

func TestNewRSAKeyPair(t *testing.T) {
	Convey("生成一对SSH密钥对", t, func() {
		privateKey, publicKey := tool.NewRSAKeyPair()
		So(privateKey, ShouldNotBeNil)
		So(publicKey, ShouldNotBeNil)
	})
}
