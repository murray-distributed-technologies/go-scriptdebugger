package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
	"github.com/libsv/go-bt/v2/bscript/interpreter/debug"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Missing arguments. Usage: go-scriptdebugger <rawtx> <lockingScript>")
	}

	/*lockingScript, err := bscript.NewFromASM("e8 OP_EQUAL")
	if err != nil {
		log.Printf("Error parsing locking script: %v\n", err)
		return
	}

	unlockingScript, err := bscript.NewFromASM("e7 OP_2 OP_ADD")
	if err != nil {
		log.Printf("Error parsing unlocking script: %v\n", err)
		return
	}*/

	debugger := debug.NewDebugger()
	debugger.AttachBeforeStep(func(state *interpreter.State) {
		var result string
		if len(state.DataStack) == 0 {
			result += "00000000  <empty>\n"
		}
		for _, frame := range state.DataStack {
			if len(frame) == 0 {
				result += "00000000  <empty>\n"
			}
			result += hex.Dump(frame)
		}
		log.Printf("Stack Before:\n\n%v\n", result)
		log.Printf("OP: %v\n", state.Opcode().Name())
	})
	debugger.AttachAfterStep(func(state *interpreter.State) {
		var result string
		for _, frame := range state.DataStack {
			if len(frame) == 0 {
				result += "00000000  <empty>\n"
			}
			result += hex.Dump(frame)
		}
		log.Printf("Stack After:\n\n%v\n", result)
		log.Printf("\n----------------------------------------------------------------\n")
		cont := ""
		fmt.Scanln(&cont)
	})
	debugger.AttachAfterStackPush(func(state *interpreter.State, data []byte) {
		frames := make([]string, len(state.DataStack))
		for i, frame := range state.DataStack {
			frames[i] = hex.EncodeToString(frame)
		}
		log.Printf("Pushing data [%x]\n", data)

	})

	txStr := os.Args[1]
	tx, err := bt.NewTxFromString(txStr)
	if err != nil {
		log.Fatalf("failed to read tx: %v", err)
	}

	scriptPubKey, err := bscript.NewFromHexString(os.Args[2])
	if err != nil {
		log.Fatalf("failed to read script: %v", err)
	}

	out := &bt.Output{LockingScript: scriptPubKey, Satoshis: uint64(5000)}

	if err := interpreter.NewEngine().Execute(
		interpreter.WithTx(tx, 0, out),
		//interpreter.WithScripts(lockingScript, unlockingScript),
		interpreter.WithAfterGenesis(),
		interpreter.WithForkID(),
		interpreter.WithDebugger(debugger),
	); err != nil {
		log.Println(err)
	}

	// Output:
	// 68656c6c6f
	// 68656c6c6f|777f726c64
	// 777f726c64|68656c6c6f
	// 777f726c6468656c6c6f
	// 8a0e597fd66749ca1a2f098f4ef706422c63a96dceef4abfd74517b10cd12f63
	// 8a0e597fd66749ca1a2f098f4ef706422c63a96dceef4abfd74517b10cd12f63|8376118fc0230e6054e782fb31ae52ebcfd551342d8d026c209997e0127b6f74
	//
	// false stack entry at end of script execution
}
