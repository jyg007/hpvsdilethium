package main


import (
	//"bytes"
	"encoding/hex"
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	"dilethium/ep11"
	pb "dilethium/grpc"
)

func main()( ) {

	cryptoClient := getGrep11Server()
	defer disconnectGrep11Server() 

    sk := make([]byte, hex.DecodedLen(len(os.Getenv("SK"))))
    hex.Decode(sk, []byte(os.Getenv("SK")))


    //*****************************************************************
	signInitRequest := &pb.SignInitRequest{
		Mech:    &pb.Mechanism{Mechanism: ep11.CKM_IBM_DILITHIUM},
		PrivKey:  sk,
	}
	signInitResponse, err := cryptoClient.SignInit(context.Background(), signInitRequest)
	if err != nil {
		panic(fmt.Errorf("SignInit error: %s", err))
	}

	signData := sha256.Sum256([]byte(os.Args[1]))
	signRequest := &pb.SignRequest{
		State: signInitResponse.State,
		Data:  signData[:],
	}

	SignResponse, err := cryptoClient.Sign(context.Background(), signRequest)
	if err != nil {
		panic(fmt.Errorf("Sign error: %s", err))
	}
 	//*****************************************************************
	
 	sign := make([]byte, hex.EncodedLen(len(SignResponse.Signature)))
    hex.Encode(sign, SignResponse.Signature)
    fmt.Println("signature: " +string(sign))

}