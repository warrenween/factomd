// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// This is a minimum wallet to be used to test the coin
// There isn't much in the way of interest in security 
// here, but rather provides a mechanism to create keys
// and sign transactions, etc.

package wallet

import (
    "fmt"
    "github.com/agl/ed25519"
    "math/rand"
    sc "github.com/FactomProject/simplecoin"
    "github.com/FactomProject/simplecoin/database"
)

// The wallet interface uses bytes.  This is because we want to 
// handle fixed length values in our maps and the database.  If
// we try to use strings, then the lengths vary based on encoding
// and that complicates the implementation without really making
// the interface more usable by developers.
type ISCWallet interface {
    database.ISCDatabase
    GenerateAddress(name []byte, m int, n int) (sc.IAddress, error)
    
    GetAddressBalance(addr sc.IAddress) uint64
    GetAddressDetailsAddr(addr []byte) IWalletEntry
    GetAddressList() (names [][]byte, addresses []sc.IAddress)
    GetAddressListByName(name []byte) (names[][]byte)
    
    /** Transaction calls **/
    CreateTransaction() sc.ITransaction
    UpdateInput(sc.ITransaction, int, sc.IAddress, uint64) error
    AddInput(sc.ITransaction, sc.IAddress, uint64) error
    AddOutput(sc.ITransaction, sc.IAddress, uint64) error
    AddECOutput(sc.ITransaction, sc.IAddress, uint64) error
    Validate(sc.ITransaction) (bool,error)
    ValidateSignatures(sc.ITransaction) bool
    SignInputs(sc.ITransaction) (bool, error)   // True if all inputs are signed
    
    GetECRate() uint64
    
    SubmitTransaction(sc.ITransaction) error
}

var factoshisPerEC uint64 = 100000

var oneSCW SCWallet

type SCWallet struct {
    ISCWallet
    db database.MapDB
    r *rand.Rand
}

var _ ISCWallet = (*SCWallet)(nil)

func (b SCWallet) String() string {
    txt,err := b.MarshalText()
    if err != nil {return "<error>" }
    return string(txt)
}

func (SCWallet) GetDBHash() sc.IHash {
    return sc.Sha([]byte("SCWallet"))
}

func (w *SCWallet) SignInputs(trans sc.ITransaction) (bool, error) {
    
    data,err := trans.MarshalBinarySig()    // Get the part of the transaction we sign
    if err != nil { return false, err }    
    
    var numSigs int = 0
    
    inputs  := trans.GetInputs()
    rcds    := trans.GetRCDs()
    for i,rcd := range rcds {
        rcd1, ok := rcd.(*sc.RCD_1)
        if ok {
            pub := rcd1.GetPublicKey()
            we := w.db.GetRaw([]byte(sc.W_ADDRESS_PUB_KEY),pub).(*WalletEntry)
            if we != nil {
                var pri [sc.SIGNATURE_LENGTH]byte
                copy(pri[:],we.private[0])
                bsig := ed25519.Sign(&pri,data)
                sig := new(sc.Signature)
                sig.SetSignature(0,bsig[:])
                sigblk := new(sc.SignatureBlock)
                sigblk.AddSignature(sig)
                trans.SetSignatureBlock(i,sigblk)
                numSigs += 1
            }
        }
    }
    return numSigs == len(inputs), nil
}

func (w *SCWallet) GetECRate() uint64 {
    return factoshisPerEC
}

func (w *SCWallet) GetAddressDetailsAddr(name []byte) IWalletEntry {
    return w.db.GetRaw([]byte("wallet.address.addr"),name).(IWalletEntry)
}

func (w *SCWallet) GenerateAddress(name []byte,m int, n int) (hash sc.IAddress, err error) {
    
    we := new(WalletEntry)
    
    nm := w.db.GetRaw([]byte("wallet.address.name"),name)
    if nm != nil {
        return nil, fmt.Errorf("Duplicate Name")
    }
    
    if m == 1 && n == 1 {
        pub,pri,err := w.generateKey()
        if err != nil { return nil, err  }
        we.AddKey(pub,pri)
        we.SetName(name)
        we.SetRCD(sc.NewRCD_1(pub))

        // If the name exists already, then we store this as the hash of the name.
        // If that exists, then we store it as the hash of the hash and so forth.
        // This way, we can get a list of addresses with the same name.
        //
        nm  := w.db.GetRaw([]byte("wallet.address.name"),name)
        switch {
            case nm == nil :       // New Name
                hash, _ = we.GetAddress()
                w.db.PutRaw([]byte(sc.W_ADDRESS_HASH),hash.Bytes(),we)
                w.db.PutRaw([]byte(sc.W_ADDRESS_PUB_KEY),pub,we)                
                w.db.PutRaw([]byte(sc.W_NAME_HASH),name,we)
            case nm != nil :       // Duplicate name.  We generate a new name, and recurse.
                return nil, fmt.Errorf("Should never get here!  This is disabled!")
                nh := sc.Sha(name)
                return w.GenerateAddress(nh.Bytes(),m, n)
            default :
                return nil, fmt.Errorf("Should never get here!  This isn't possible!")
        }
        
    } else {
        return nil, fmt.Errorf("Not this far yet!")
    }
    return
}

func (w *SCWallet) Init (a ...interface{}) {
    if w.r != nil { return }
    w.r = rand.New(rand.NewSource(13436523)) 
    w.db.Init()
}
    
func (w *SCWallet) Read(buf []byte) (int, error) {
    //if r==nil { r = rand.New(rand.NewSource(time.Now().Unix())) }
    w.Init()
    for i := range buf {
        buf[i] = byte(w.r.Int())
    }
    return len(buf), nil
}

func (w *SCWallet) generateKey() (public []byte,private []byte, err error){
    pub,pri,err := ed25519.GenerateKey(w)
    return pub[:], pri[:], err
}

func (w *SCWallet)  CreateTransaction() sc.ITransaction {
    return new(sc.Transaction)
}

func (w *SCWallet) AddInput(trans sc.ITransaction, hash sc.IAddress, amount uint64) error {
    v := w.db.GetRaw([]byte(sc.W_ADDRESS_HASH),hash.Bytes())
          
     if(v == nil) {
         return fmt.Errorf("Unknown address")
     }
     we := v.(*WalletEntry)
     adr, err := we.GetAddress()
     trans.AddInput(sc.CreateAddress(adr),amount)
     trans.AddRCD(we.GetRCD())
     if err != nil {
         return err
     }
     return nil
}     

func (w *SCWallet) UpdateInput(trans sc.ITransaction, 
                               index int, 
                               hash sc.IAddress, 
                               amount uint64) error {
    
    v := w.db.GetRaw([]byte(sc.W_ADDRESS_HASH),hash.Bytes())
    if v == nil { 
        return fmt.Errorf("Unknown address")
    }
    we := v.(*WalletEntry)
    adr, err := we.GetAddress()
    in,err := trans.GetInput(index)
    if err != nil {return err}
    in.SetAddress(adr)
    in.SetAmount(amount)
 
    return nil
}     

func (w *SCWallet) AddOutput(trans sc.ITransaction, hash sc.IAddress, amount uint64) error {
    v := w.db.GetRaw([]byte(sc.W_ADDRESS_HASH),hash.Bytes())
    if v == nil { 
        return fmt.Errorf("Unknown address")
    }
    we := v.(*WalletEntry)
    adr, err := we.GetAddress()
    trans.AddOutput(sc.CreateAddress(adr),amount)
    if err != nil {
        return err
    }
    return nil
}     
 
func (w *SCWallet) AddECOutput(trans sc.ITransaction, hash sc.IAddress, amount uint64) error {
    we := w.db.GetRaw([]byte(sc.W_ADDRESS_HASH),hash.Bytes()).(*WalletEntry)
    if we == nil { 
        return fmt.Errorf("Unknown address")
    }
    adr, err := we.GetAddress()
    trans.AddECOutput(sc.CreateAddress(adr),amount)
    if err != nil {
        return err
    }
    return nil
}

func (w *SCWallet) Validate(trans sc.ITransaction) (bool,error){
    valid := trans.Validate()
    if valid == sc.WELL_FORMED {
        return true, nil
    }
    fmt.Println("Validation Failed: ",valid)
    return false, nil
}    
 
 func (w *SCWallet) ValidateSignatures(trans sc.ITransaction) bool{
     valid := trans.ValidateSignatures()
     return valid
 } 
 
