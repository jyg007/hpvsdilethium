package main

import (
	"encoding/hex"
	"context"
	"fmt"
	"os"
	"crypto/sha256"

	"dilethium/ep11"
	pb "dilethium/grpc"
)


func main()( ) {

    cryptoClient := getGrep11Server()
   	defer disconnectGrep11Server() 


    pk := make([]byte, hex.DecodedLen(len(os.Getenv("PK"))))
    hex.Decode(pk, []byte(os.Getenv("PK")))

    sign := make([]byte, hex.DecodedLen(len(os.Getenv("SIGN"))))
    hex.Decode(sign, []byte(os.Getenv("SIGN")))

	signData := sha256.Sum256([]byte(os.Args[1]))
 
    //*****************************************************************
 	verifyInitRequest := &pb.VerifyInitRequest{
		Mech:   &pb.Mechanism{Mechanism: ep11.CKM_IBM_DILITHIUM},
		PubKey: pk,
	}
	verifyInitResponse, err := cryptoClient.VerifyInit(context.Background(), verifyInitRequest)
	if err != nil {
		panic(fmt.Errorf("VerifyInit error: %s", err))
	}

	verifyRequest := &pb.VerifyRequest{
		State:     verifyInitResponse.State,
		Data:      []byte(signData[:]),
		Signature: sign,
	}

	_, err = cryptoClient.Verify(context.Background(), verifyRequest)
	if ok, ep11Status := Convert(err); !ok {
		if ep11Status.Code == ep11.CKR_SIGNATURE_INVALID {
			panic(fmt.Errorf("Invalid signature"))
		} else {
			panic(fmt.Errorf("Verify error: [%d]: %s", ep11Status.Code, ep11Status.Detail))
		}
	}
	fmt.Println("Verified")
	//*****************************************************************
	
}