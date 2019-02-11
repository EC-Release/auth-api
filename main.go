/*
 * Copyright (c) 2016 General Electric Company. All rights reserved.
 *
 * The copyright to the computer software herein is the property of
 * General Electric Company. The software may be used and/or copied only
 * with the written permission of General Electric Company or in accordance
 * with the terms and conditions stipulated in the agreement/contract
 * under which the software has been supplied.
 *
 * author: apolo.yasuda@ge.com
 */

package main

import (
	api "github.build.ge.com/212359746/wzapi"
	util "github.build.ge.com/212359746/wzutil"
	"net/http"
	"crypto/rsa"
	"strings"
	"errors"
	"encoding/base64"
	"os"
)

var (
	PVT_PWD = os.Getenv("EC_PRVT_PWD")
	PVT_KEY = os.Getenv("EC_PRVT_KEY")
	EC_CRT = os.Getenv("EC_PUB_KEY")
	ADMIN_USR = os.Getenv("ADMIN_USR")
	ADMIN_TKN = os.Getenv("ADMIN_TKN")
)

const (
	EC_HTTP_HEADER = "ec-options"
)
func main(){

	util.Branding("/.ec","ec-plugin","ec-config","https://ca-not-in-use.com","EC")
	util.Init("agent",true)

	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request){

		defer func(){
			if r:=recover();r!=nil{
				util.PanicRecovery(r)
				api.ErrResponse(w, 500, errors.New("internal error"), r.(string))
			}
		}()
		
		usr, tkn, ok:=r.BasicAuth()
		if !ok {
			api.ErrResponse(w, 401, errors.New("internal error"), "not a basic aithentication request.")
			return 
		}
	
		if usr!=ADMIN_USR || tkn!= ADMIN_TKN {
			api.ErrResponse(w, 401, errors.New("internal error"), "authentication failed.")
			return
		}
		
		w.Header().Set("Content-Type", "application/json")

		_opt := r.Header.Get(EC_HTTP_HEADER)

		crt,err:=util.NewCert("minota")
		if err!=nil{
			api.ErrResponse(w, 500, err, "")
			return	
		}
		
		pk,err:=crt.ParsePvtKey([]byte(PVT_KEY), PVT_PWD)
		if err!=nil{
			api.ErrResponse(w, 500, err, "")
			return
		}
		

		pp,err:=base64.StdEncoding.DecodeString(_opt)
		if err!=nil{
			api.ErrResponse(w, 500, err, "")
			return 
		}

		pq,err:=decrypt(string(pp),pk,crt)
		if err!=nil{
			api.ErrResponse(w, 500, err, "")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"`+base64.StdEncoding.EncodeToString([]byte(pq))+`"}`))

	})
	
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request){

		defer func(){
			if r:=recover();r!=nil{
				util.PanicRecovery(r)
				api.ErrResponse(w, 500, errors.New("internal error"), r.(string))
			}
		}()

		usr, tkn, ok:=r.BasicAuth()
		if !ok {
			api.ErrResponse(w, 401, errors.New("internal error"), "not a basic aithentication request.")
			return 
		}
		
		if usr!=ADMIN_USR || tkn!= ADMIN_TKN {
			api.ErrResponse(w, 401, errors.New("internal error"), "authentication failed.")
			return
			
		}

		w.Header().Set("Content-Type", "application/json")

		_opt := r.Header.Get(EC_HTTP_HEADER)
				
		op,err:=encrypt(_opt,[]byte(EC_CRT))
		if err!=nil{
			api.ErrResponse(w, 500, err, "")
			return 
		}
		
		op=base64.StdEncoding.EncodeToString([]byte(op))
		util.DbgLog(op)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"`+op+`","test":"`+PVT_PWD+`"}`))

	})

	util.InfoLog("decrypt api is up and running.")

	err:=http.ListenAndServe(":8990", nil)
	if err!=nil{
		panic(err)
	}	
}

func decrypt(d string, pk *rsa.PrivateKey, crt *util.Cert) (string, error){

	_s,err:=crt.Decrypt([]byte(strings.TrimSpace(d)),pk)
	if err!=nil{
		return "",err
	}
	
	return string(_s),nil
}

func encrypt(d string, pbk []byte) (string, error){

	crt,err:=util.NewCert("minota")
	if err==nil{
		return "",err
	}
	
	_s:=crt.EncryptV2(strings.TrimSpace(d), pbk)
	if _s==nil{
		return "",errors.New("encrypt failed.")
	}
	
	return string(_s),nil
}
