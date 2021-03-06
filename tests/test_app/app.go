package main

import (
	"bytes"
	"fmt"
	"os"

	asuracli "github.com/teragrid/asura/client"
	"github.com/teragrid/asura/types"
	"github.com/teragrid/teralibs/log"
)

func startClient(asuraType string) asuracli.Client {
	// Start client
	client, err := asuracli.NewClient("tcp://127.0.0.1:46658", asuraType, true)
	if err != nil {
		panic(err.Error())
	}
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	client.SetLogger(logger.With("module", "asuracli"))
	if err := client.Start(); err != nil {
		panicf("connecting to asura_app: %v", err.Error())
	}

	return client
}

func setOption(client asuracli.Client, key, value string) {
	_, err := client.SetOptionSync(types.RequestSetOption{key, value})
	if err != nil {
		panicf("setting %v=%v: \nerr: %v", key, value, err)
	}
}

func commit(client asuracli.Client, hashExp []byte) {
	res, err := client.CommitSync()
	if err != nil {
		panicf("client error: %v", err)
	}
	if !bytes.Equal(res.Data, hashExp) {
		panicf("Commit hash was unexpected. Got %X expected %X", res.Data, hashExp)
	}
}

func deliverTx(client asuracli.Client, txBytes []byte, codeExp uint32, dataExp []byte) {
	res, err := client.DeliverTxSync(txBytes)
	if err != nil {
		panicf("client error: %v", err)
	}
	if res.Code != codeExp {
		panicf("DeliverTx response code was unexpected. Got %v expected %v. Log: %v", res.Code, codeExp, res.Log)
	}
	if !bytes.Equal(res.Data, dataExp) {
		panicf("DeliverTx response data was unexpected. Got %X expected %X", res.Data, dataExp)
	}
}

/*func checkTx(client asuracli.Client, txBytes []byte, codeExp uint32, dataExp []byte) {
	res, err := client.CheckTxSync(txBytes)
	if err != nil {
		panicf("client error: %v", err)
	}
	if res.IsErr() {
		panicf("checking tx %X: %v\nlog: %v", txBytes, res.Log)
	}
	if res.Code != codeExp {
		panicf("CheckTx response code was unexpected. Got %v expected %v. Log: %v",
			res.Code, codeExp, res.Log)
	}
	if !bytes.Equal(res.Data, dataExp) {
		panicf("CheckTx response data was unexpected. Got %X expected %X",
			res.Data, dataExp)
	}
}*/

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
