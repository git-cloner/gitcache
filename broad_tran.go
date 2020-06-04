package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

type ConfigStruct struct {
	GenesisPublickey  string `json:"genesis_publickey"`
	GenesisPrivatekey string `json:"genesis_privatekey"`
	PeerPublickey     string `json:"peer_publickey"`
}

func GetPublic() ConfigStruct {
	// open config.json
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result ConfigStruct
	if err := json.Unmarshal([]byte(byteValue), &result); err != nil {
		fmt.Println(err)
	}
	return result
}

/**
* toPublicKey:get token publicKey
* amount: token's amount
* repostory: clone adress
 */
func BroadcastTransaction(amount string, repostory string) ResultObj {
	var configJson = GetPublic()

	// create msg
	var broadcastTranMsg BroadTranMsg
	var Code = "CODE"
	var PublicKey = configJson.GenesisPublickey
	broadcastTranMsg.Token = Code
	broadcastTranMsg.From = PublicKey
	broadcastTranMsg.To = configJson.PeerPublickey
	broadcastTranMsg.Amount = amount
	broadcastTranMsg.Repostory = repostory
	jsonstring, _ := json.Marshal(broadcastTranMsg)
	// encode msg
	var base64msg = base64.StdEncoding.EncodeToString([]byte(jsonstring))
	// sign
	var PrivateKey = configJson.GenesisPrivatekey
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
