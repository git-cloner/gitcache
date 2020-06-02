package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

type BroadTranMsg struct {
	Token     string
	From      string
	To        string
	Amount    string
	Repostory string
}

type BroadTranInfo struct {
	Public string
	Sign   string
	Msg    string
}

/**
* toPublicKey:get token publicKey
* amount: token's amount
* repostory: clone adress
 */
func BroadcastTransaction(toPublicKey string, amount string, repostory string) ResultObj {
	// create msg
	var broadcastTranMsg BroadTranMsg
	var Code = "CODE"
	var PublicKey = "9CE9108CD5243B401CD1A7EDE1921E3F7FCF3ADDCDD70F0AA5AF31350D1B55B1"
	broadcastTranMsg.Token = Code
	broadcastTranMsg.From = PublicKey
	broadcastTranMsg.To = toPublicKey
	broadcastTranMsg.Amount = amount
	broadcastTranMsg.Repostory = repostory
	jsonstring, _ := json.Marshal(broadcastTranMsg)
	// encode msg
	var base64msg = base64.StdEncoding.EncodeToString([]byte(jsonstring))
	// sign
	var PrivateKey = "3d3a226a1c3f72af75270cc4f9475ebfbd7c7150b313aa3a27ee73a67a8df15f9ce9108cd5243b401cd1a7ede1921e3f7fcf3addcdd70f0aa5af31350d1b55b1"
	_privatekey, _ := hex.DecodeString(PrivateKey)
	var privateKey ed25519.PrivKeyEd25519
	copy(privateKey[:], _privatekey)
	signStr, err := privateKey.Sign([]byte(base64msg))
	// define response
	var res ResultObj

	if err == nil {
		// sign successfully
		sign := hex.EncodeToString(signStr)
		url := "http://localhost:26657"
		// defined struct
		var baseInitData = "{" +
			"\"publickey\":\"" + PublicKey + "\"," +
			"\"sign\":\"" + sign + "\"," +
			"\"msg\":\"" + base64msg + "\"" +
			"}"
		var baseInput = []byte(baseInitData)
		var encodingString = base64.StdEncoding.EncodeToString(baseInput)
		var post = "{\"method\":\"broadcast_tx_commit\",\"jsonrpc\":\"2.0\",\"params\":{\"tx\":\"" + encodingString + "\"},\"id\":\"\"}"
		var jsonStr = []byte(post)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			res.Result = false
			res.Info = ""
			res.Error = err.Error()
			return res
		}
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := &http.Client{}
		// send request
		resp, err := client.Do(req)

		if err != nil {
			res.Result = false
			res.Info = ""
			res.Error = err.Error()
			return res
		}
		defer resp.Body.Close()
		// reponse result
		body, err := ioutil.ReadAll((resp.Body))
		if resp.StatusCode == 200 {
			res.Result = true
			res.Info = string(body)
			res.Error = ""
			return res
		} else {
			// sign failed
			res.Result = false
			res.Info = ""
			res.Error = err.Error()
			return res
		}
	} else {
		res.Result = false
		res.Info = ""
		res.Error = err.Error()
		return res
	}
}
