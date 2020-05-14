package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

type MsgTx struct {
	PrivateKey string
	PublicKey  string
	Msg        string
}

type ResultObj struct {
	Result bool
	Info   string
	Error  string
}

// post {method jsonrpc params id} to 26657/broadcast_tx_commit
func BroadCastMsg(json MsgTx) ResultObj {
	// encode msg
	var base64msg = base64.StdEncoding.EncodeToString([]byte(json.Msg))
	// sign
	_privatekey, _ := hex.DecodeString(json.PrivateKey)
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
			"\"publickey\":\"" + json.PublicKey + "\"," +
			"\"sign\":\"" + sign + "\"," +
			"\"msg\":\"" + base64msg + "\"" +
			"}"
		fmt.Println(baseInitData)
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
