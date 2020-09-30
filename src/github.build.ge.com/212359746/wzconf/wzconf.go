//go:binary-only-package

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

package wzconf
import (
	"os/signal"
	"time"
	"encoding/base64"
	//"crypto/tls"
	//"crypto/rand"
	"fmt"
	"os"
	"errors"
	//"os/signal"
	//"crypto/x509"
	"net"
	"strconv"
	"strings"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
	"github.com/pborman/uuid"
	model "github.build.ge.com/212359746/wzschema"
	util "github.build.ge.com/212359746/wzutil"
	plugin "github.build.ge.com/212359746/wzplugin"
	api "github.build.ge.com/212359746/wzapi"
)
