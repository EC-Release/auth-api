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
	api "github.com/wzlib/wzapi"
	util "github.com/wzlib/wzutil"
	"net/http"
	"crypto/rsa"
	"strings"
	"errors"
	"encoding/base64"
	"os"
)

var (
	PVT_PWD = os.Getenv("EC_PRVT_PWD")
	PVT_KEY = "./service.key"
	EC_CRT = "./service.crt"
	ADMIN_USR = os.Getenv("ADMIN_USR")
	ADMIN_TKN = os.Getenv("ADMIN_TKN")
)

const (
	EC_LOGO = `
           ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄                                            
          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░
          ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀   
          ▐░▌          ▐░▌            
          ▐░█▄▄▄▄▄▄▄▄▄ ▐░▌            
          ▐░░░░░░░░░░░▌▐░▌            
          ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌            
          ▐░▌          ▐░▌            
          ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄   
          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌  
           ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  @Digital Connect 
`
	COPY_RIGHT = "Digital Connect,  @GE Corporate"
	ISSUE_TRACKER = "https://github.com/EC-Release/ec-sdk/issues"
	TC_HEADER = "X-Thread-Connect"
	XCALR_URL = "https://x-thread-connect.run.pcs.aws-usw02-dev.ice.predix.io"
	EC_HTTP_HEADER = "ec-options"
)
func main(){

	util.Branding("/.ec","ec-plugin","ec-config",TC_HEADER,"EC",EC_LOGO,COPY_RIGHT,XCALR_URL,ISSUE_TRACKER)
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
		
		pkey,err:=util.ReadFile(PVT_KEY)
		if err!=nil {
			api.ErrResponse(w, 500, err, "")
			return
		}
		
		pk,err:=crt.ParsePvtKey(pkey, PVT_PWD)
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
		_,err = w.Write([]byte(`{"data":"`+base64.StdEncoding.EncodeToString([]byte(pq))+`"}`))
		if err != nil {
			panic(err)
		}
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
		
		pbk,err:=util.ReadFile(EC_CRT)
		if err!=nil {
			api.ErrResponse(w, 500, err, "")
			return
		}
		
		op,err:=encrypt(_opt,pbk)
		if err!=nil{
			api.ErrResponse(w, 500, err, "")
			return 
		}
		
		op=base64.StdEncoding.EncodeToString([]byte(op))
		util.DbgLog(op)
		
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"data":"`+op+`","test":"`+PVT_PWD+`"}`))
		if err != nil {
			panic(err)
		}
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
