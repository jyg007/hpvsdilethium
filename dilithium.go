package main


import (
	//"bytes"
	"encoding/hex"
	"context"
	"encoding/asn1"
	"fmt"
	"hpvsdilithium/ep11"
	pb "hpvsdilithium/grpc"

)

func main()( ) {

	cryptoClient := getGrep11Server()
	defer disconnectGrep11Server() 

	dilithiumStrengthParam, err := asn1.Marshal(asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 2, 267, 1, 6, 5}) // Round 2 strength)
	if err != nil {
		panic(fmt.Errorf("Unable to encode parameter OID: %s", err))
	}

	publicKeyTemplate := ep11.EP11Attributes{
		ep11.CKA_IBM_PQC_PARAMS: dilithiumStrengthParam,
		ep11.CKA_VERIFY:         true,
		ep11.CKA_EXTRACTABLE:    false,
	}
	privateKeyTemplate := ep11.EP11Attributes{
		ep11.CKA_SIGN:        true,
		ep11.CKA_EXTRACTABLE: false,
	}
	generateDilKeyPairRequest := &pb.GenerateKeyPairRequest{
		Mech:            &pb.Mechanism{Mechanism: ep11.CKM_IBM_DILITHIUM},
		PubKeyTemplate:  AttributeMap(publicKeyTemplate),
		PrivKeyTemplate: AttributeMap(privateKeyTemplate),
	}

	// Dilithium Key Pair generation
	generateDilKeyPairResponse, err := cryptoClient.GenerateKeyPair(context.Background(), generateDilKeyPairRequest)
	if ok, ep11Status := Convert(err); !ok {
		if ep11Status.Code == ep11.CKR_MECHANISM_INVALID {
			fmt.Println("Dilithium mechanism is not supported on the remote HSM")
			return
		} else {
			panic(fmt.Errorf("Generate Dilithium key pair error: %s", err))
		}
	}

    pubkeyhex := make([]byte, hex.EncodedLen(len(generateDilKeyPairResponse.PubKeyBytes)))
    hex.Encode(pubkeyhex, generateDilKeyPairResponse.PubKeyBytes)
    fmt.Println("Public Dilithium key: " +string(pubkeyhex)+"\n")

    privkeyhex := make([]byte, hex.EncodedLen(len(generateDilKeyPairResponse.PrivKeyBytes)))
    hex.Encode(privkeyhex, generateDilKeyPairResponse.PrivKeyBytes)
    fmt.Println("Private Dilithium key: " +string(privkeyhex))

}
