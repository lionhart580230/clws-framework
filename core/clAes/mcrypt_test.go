package clAes

import "testing"
import "fmt"

const encryptKey = "b2ee8d97f922bafd659f16061c6e41cc"

func TestEncode(t *testing.T) {

	pValues := Encode([]byte(`p=123&hao=123123&nn=dddd`), []byte(encryptKey))

	fmt.Printf(">> pValues: (%v)\n", string(pValues))
}


func TestDecodeParam(t *testing.T) {


	DecodeParam("eyJpdiI6IlRqRjNkemR0WlVOWlpVOU1lRE5JWVE9PSIsInZhbHVlIjoiY1phZ0tSbEE2S3hhb1J1YXM3a1U3VXAzZ0Q0ckExcmRSWHlHa1F5M1Y0MD0ifQ==", encryptKey)

}