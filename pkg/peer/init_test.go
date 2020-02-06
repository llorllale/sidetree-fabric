/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	cb "github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/extensions/collections/storeprovider"
	"github.com/hyperledger/fabric/extensions/gossip/blockpublisher"
	viper "github.com/spf13/viper2015"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/fabric-peer-ext/pkg/common/compositekey"
	extconfig "github.com/trustbloc/fabric-peer-ext/pkg/config"
	ledgercfg "github.com/trustbloc/fabric-peer-ext/pkg/config/ledgerconfig/config"
	ledgercfgmgr "github.com/trustbloc/fabric-peer-ext/pkg/config/ledgerconfig/mgr"
	statemocks "github.com/trustbloc/fabric-peer-ext/pkg/gossip/state/mocks"
	"github.com/trustbloc/fabric-peer-ext/pkg/mocks"
	extpeer "github.com/trustbloc/fabric-peer-ext/pkg/peer"
	"github.com/trustbloc/fabric-peer-ext/pkg/resource"
	extroles "github.com/trustbloc/fabric-peer-ext/pkg/roles"
	"github.com/trustbloc/sidetree-fabric/pkg/peer/config"
	"github.com/trustbloc/sidetree-fabric/pkg/role"
)

//go:generate counterfeiter -o ../mocks/peerledger.gen.go --fake-name Ledger github.com/hyperledger/fabric/core/ledger.PeerLedger

const (
	channelID = "mychannel"
	mspID     = "Org1MSP"
	peerID    = "peer1.example.com"
	configSCC = "configscc"

	listenPort = 43900

	docNotFound  = "document not found"
	pageNotFound = "404 page not found"

	v1   = "1"
	v0_4 = "0.4"
	v0_5 = "0.5"

	tx1 = "tx1"
	tx2 = "tx2"
	tx3 = "tx3"
	tx4 = "tx4"
	tx5 = "tx5"
	tx6 = "tx6"

	peerSidetreeCfgJson  = `{"Monitor":{"Period":"5s"},"Namespaces":[{"Namespace":"did:sidetree","BasePath":"/document"}]}`
	peerTrustblocCfgJson = `{"Monitor":{"Period":"5s"},"Namespaces":[{"Namespace":"did:bloc:trustbloc.dev","BasePath":"/trustbloc.dev"}]}`
	peerBothCfgJson      = `{"Monitor":{"Period":"5s"},"Namespaces":[{"Namespace":"did:sidetree","BasePath":"/document"},{"Namespace":"did:bloc:trustbloc.dev","BasePath":"/trustbloc.dev"}]}`
	peerNoneCfgJson      = `{}`

	didTrustblocNamespace             = "did:bloc:trustbloc.dev"
	didTrustblocBasePath              = "/trustbloc.dev"
	didTrustblocProtocol_V0_5_CfgJSON = `{"startingBlockchainTime":250000,"hashAlgorithmInMultihashCode":18,"maxOperationByteSize":20000,"maxOperationsPerBatch":200}`
	didTrustblocCfgYaml               = `batchWriterTimeout: 1s`
	didTrustblocCfgUpdateYaml         = `batchWriterTimeout: 5s`

	didSidetreeNamespace             = "did:sidetree"
	didSidetreeBasePath              = "/document"
	didSidetreeCfgJSON               = `{"batchWriterTimeout":"5s"}`
	didSidetreeProtocol_V0_4_CfgJSON = `{"startingBlockchainTime":200000,"hashAlgorithmInMultihashCode":18,"maxOperationByteSize":2000,"maxOperationsPerBatch":10}`
	didSidetreeProtocol_V0_5_CfgJSON = `{"startingBlockchainTime":500000,"hashAlgorithmInMultihashCode":18,"maxOperationByteSize":10000,"maxOperationsPerBatch":100}`
)

var (
	// Ensure that the provider instances are instantiated and registered as a resource
	_ = blockpublisher.ProviderInstance
	_ = storeprovider.NewProviderFactory()
	_ = extroles.GetRoles()

	didTrustblocCfgKeyBytes         = ledgercfgmgr.MarshalKey(ledgercfg.NewAppKey(config.GlobalMSPID, didTrustblocNamespace, v1))
	didTrustblocCfgValueBytes       = marshalConfigValue(tx1, didTrustblocCfgYaml, "yaml")
	didTrustblocCfgUpdateValueBytes = marshalConfigValue(tx3, didTrustblocCfgUpdateYaml, "yaml")

	didTrustblocProtocol_v0_5_CfgKeyBytes   = ledgercfgmgr.MarshalKey(ledgercfg.NewComponentKey(config.GlobalMSPID, didTrustblocNamespace, v1, config.ProtocolComponentName, v0_5))
	didTrustblocProtocol_v0_5_CfgValueBytes = marshalConfigValue(tx1, didTrustblocProtocol_V0_5_CfgJSON, "json")

	didSidetreeCfgKeyBytes   = ledgercfgmgr.MarshalKey(ledgercfg.NewAppKey(config.GlobalMSPID, didSidetreeNamespace, v1))
	didSidetreeCfgValueBytes = marshalConfigValue(tx1, didSidetreeCfgJSON, "json")

	didSidetreeProtocol_v0_4_CfgKeyBytes   = ledgercfgmgr.MarshalKey(ledgercfg.NewComponentKey(config.GlobalMSPID, didSidetreeNamespace, v1, config.ProtocolComponentName, v0_4))
	didSidetreeProtocol_v0_4_CfgValueBytes = marshalConfigValue(tx1, didSidetreeProtocol_V0_4_CfgJSON, "json")

	didSidetreeProtocol_v0_5_CfgKeyBytes   = ledgercfgmgr.MarshalKey(ledgercfg.NewComponentKey(config.GlobalMSPID, didSidetreeNamespace, v1, config.ProtocolComponentName, v0_5))
	didSidetreeProtocol_v0_5_CfgValueBytes = marshalConfigValue(tx1, didSidetreeProtocol_V0_5_CfgJSON, "json")

	peerCfgKeyBytes            = ledgercfgmgr.MarshalKey(ledgercfg.NewPeerKey(mspID, peerID, config.SidetreeAppName, config.SidetreeAppVersion))
	peerSidetreeCfgValueBytes  = marshalConfigValue(tx1, peerSidetreeCfgJson, "json")
	peerTrustblocCfgValueBytes = marshalConfigValue(tx5, peerTrustblocCfgJson, "json")
	peerBothCfgValueBytes      = marshalConfigValue(tx2, peerBothCfgJson, "json")
	peerNoneCfgValueBytes      = marshalConfigValue(tx6, peerNoneCfgJson, "json")
)

func TestInitialize(t *testing.T) {
	defer removeDBPath(t)

	viper.Set("sidetree.port", listenPort)

	restore := setRoles(role.Observer, role.BatchWriter, role.Resolver)
	defer restore()

	qe := mocks.NewQueryExecutor().
		WithState(configSCC, didTrustblocCfgKeyBytes, didTrustblocCfgValueBytes).
		WithState(configSCC, getIndexKey(didTrustblocCfgKeyBytes, []string{config.GlobalMSPID}), []byte("{}")).
		WithState(configSCC, didSidetreeCfgKeyBytes, didSidetreeCfgValueBytes).
		WithState(configSCC, getIndexKey(didSidetreeCfgKeyBytes, []string{config.GlobalMSPID}), []byte("{}")).
		WithState(configSCC, peerCfgKeyBytes, peerNoneCfgValueBytes).
		WithState(configSCC, getIndexKey(peerCfgKeyBytes, []string{mspID}), []byte("{}")).
		WithState(configSCC, didSidetreeProtocol_v0_4_CfgKeyBytes, didSidetreeProtocol_v0_4_CfgValueBytes).
		WithState(configSCC, getIndexKey(didSidetreeProtocol_v0_4_CfgKeyBytes, []string{config.GlobalMSPID}), []byte("{}")).
		WithState(configSCC, didSidetreeProtocol_v0_5_CfgKeyBytes, didSidetreeProtocol_v0_5_CfgValueBytes).
		WithState(configSCC, getIndexKey(didSidetreeProtocol_v0_5_CfgKeyBytes, []string{config.GlobalMSPID}), []byte("{}")).
		WithState(configSCC, didTrustblocProtocol_v0_5_CfgKeyBytes, didTrustblocProtocol_v0_5_CfgValueBytes).
		WithState(configSCC, getIndexKey(didTrustblocProtocol_v0_5_CfgKeyBytes, []string{config.GlobalMSPID}), []byte("{}"))

	req := &Require{require.New(t)}

	req.NotPanics(extpeer.Initialize)
	req.NotPanics(Initialize)

	req.NoError(
		resource.Mgr.Initialize(
			blockpublisher.ProviderInstance,
			newMockLedgerProvider(qe),
			newMockPeerConfigPrivider(),
			&mocks.GossipProvider{},
			&mocks.IdentityDeserializerProvider{},
			&mocks.IdentifierProvider{},
			&mocks.IdentityProvider{},
			&statemocks.CCEventMgrProvider{},
			&mockBatchWriterConfig{batchTimeout: time.Second},
			&mockRESTServiceConfig{listenURL: "localhost:8978"},
		),
	)

	defer resource.Mgr.Close()

	req.NotPanics(func() { resource.Mgr.ChannelJoined(channelID) })

	// Give the services a chance to startup
	time.Sleep(200 * time.Millisecond)

	// The REST service shouldn't be started with no namespaces configured on the peer
	req.ConnectionRefused()

	t.Run("Update peer config with only did:sidetree namespace", func(t *testing.T) {
		blockBuilder := mocks.NewBlockBuilder(channelID, 1000)
		blockBuilder.Transaction(tx2, pb.TxValidationCode_VALID).
			ChaincodeAction(configSCC).
			Write(peerCfgKeyBytes, peerSidetreeCfgValueBytes)
		blockpublisher.ForChannel(channelID).Publish(blockBuilder.Build())

		time.Sleep(200 * time.Millisecond)

		req.Response(didSidetreeBasePath, didSidetreeNamespace, http.StatusNotFound, docNotFound)
		req.Response(didTrustblocBasePath, didTrustblocNamespace, http.StatusNotFound, pageNotFound)
	})

	t.Run("Update peer config with did:sidetree and did:bloc:trustbloc.dev namespaces", func(t *testing.T) {
		blockBuilder := mocks.NewBlockBuilder(channelID, 1000)
		blockBuilder.Transaction(tx2, pb.TxValidationCode_VALID).
			ChaincodeAction(configSCC).
			Write(peerCfgKeyBytes, peerBothCfgValueBytes)
		blockpublisher.ForChannel(channelID).Publish(blockBuilder.Build())

		time.Sleep(200 * time.Millisecond)

		req.Response(didSidetreeBasePath, didSidetreeNamespace, http.StatusNotFound, docNotFound)
		req.Response(didTrustblocBasePath, didTrustblocNamespace, http.StatusNotFound, docNotFound)
	})

	t.Run("Update consortium config", func(t *testing.T) {
		blockBuilder := mocks.NewBlockBuilder(channelID, 1001)
		blockBuilder.Transaction(tx3, pb.TxValidationCode_VALID).
			ChaincodeAction(configSCC).
			Write(didTrustblocCfgKeyBytes, didTrustblocCfgUpdateValueBytes)
		blockpublisher.ForChannel(channelID).Publish(blockBuilder.Build())

		time.Sleep(200 * time.Millisecond)

		req.Response(didSidetreeBasePath, didSidetreeNamespace, http.StatusNotFound, docNotFound)
		req.Response(didTrustblocBasePath, didTrustblocNamespace, http.StatusNotFound, docNotFound)
	})

	t.Run("Update irrelevant config", func(t *testing.T) {
		// Arbitrary config update should be ignored
		someCfgKeyBytes := ledgercfgmgr.MarshalKey(ledgercfg.NewPeerKey(mspID, peerID, "some-app", "1"))
		someCfgValueBytes := marshalConfigValue(tx4, "some-config", "other")

		blockBuilder := mocks.NewBlockBuilder(channelID, 1002)
		blockBuilder.Transaction(tx4, pb.TxValidationCode_VALID).
			ChaincodeAction(configSCC).
			Write(someCfgKeyBytes, someCfgValueBytes)
		blockpublisher.ForChannel(channelID).Publish(blockBuilder.Build())

		time.Sleep(200 * time.Millisecond)

		req.Response(didSidetreeBasePath, didSidetreeNamespace, http.StatusNotFound, docNotFound)
		req.Response(didTrustblocBasePath, didTrustblocNamespace, http.StatusNotFound, docNotFound)
	})

	t.Run("Update peer config with only did:bloc:trustbloc.dev namespace", func(t *testing.T) {
		blockBuilder := mocks.NewBlockBuilder(channelID, 1003)
		blockBuilder.Transaction(tx5, pb.TxValidationCode_VALID).
			ChaincodeAction(configSCC).
			Write(peerCfgKeyBytes, peerTrustblocCfgValueBytes)
		blockpublisher.ForChannel(channelID).Publish(blockBuilder.Build())

		time.Sleep(200 * time.Millisecond)

		req.Response(didSidetreeBasePath, didSidetreeNamespace, http.StatusNotFound, pageNotFound)
		req.Response(didTrustblocBasePath, didTrustblocNamespace, http.StatusNotFound, docNotFound)
	})

	t.Run("Update peer config with no namespaces", func(t *testing.T) {
		// This update removes all of the Sidetree peer config, which should cause all of the Sidetree services to be stopped
		blockBuilder := mocks.NewBlockBuilder(channelID, 1003)
		blockBuilder.Transaction(tx6, pb.TxValidationCode_VALID).
			ChaincodeAction(configSCC).
			Write(peerCfgKeyBytes, peerNoneCfgValueBytes)
		blockpublisher.ForChannel(channelID).Publish(blockBuilder.Build())

		time.Sleep(200 * time.Millisecond)

		// The REST service shouldn't even be started
		req.ConnectionRefused()
	})
}

func removeDBPath(t testing.TB) {
	removePath(t, extconfig.GetTransientDataLevelDBPath())
}

func removePath(t testing.TB, path string) {
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf(err.Error())
	}
}

type mockBatchWriterConfig struct {
	batchTimeout time.Duration
}

func (m *mockBatchWriterConfig) GetBatchTimeout() time.Duration {
	return m.batchTimeout
}

type mockRESTServiceConfig struct {
	listenURL string
}

func (m *mockRESTServiceConfig) GetListenURL() string {
	return m.listenURL
}

func getIndexKey(key string, fields []string) string {
	return compositekey.Create("cfgmgmt-mspid", append(fields, key))
}

func marshalConfigValue(txID, cfg string, format ledgercfg.Format) []byte {
	bytes, err := json.Marshal(&ledgercfg.Value{TxID: txID, Format: format, Config: cfg})
	if err != nil {
		panic(err)
	}
	return bytes
}

func newMockLedgerProvider(qe *mocks.QueryExecutor) *mocks.LedgerProvider {
	ledgerProvider := &mocks.LedgerProvider{}
	ledgerProvider.GetLedgerReturns(
		&mocks.Ledger{
			BlockchainInfo: &cb.BlockchainInfo{Height: 1000},
			QueryExecutor:  qe,
		},
	)

	return ledgerProvider
}

func newMockPeerConfigPrivider() *mocks.PeerConfig {
	p := &mocks.PeerConfig{}
	p.MSPIDReturns(mspID)
	p.PeerIDReturns(peerID)

	return p
}

func setRoles(roles ...extroles.Role) func() {
	rolesValue := make(map[extroles.Role]struct{})
	for _, role := range roles {
		rolesValue[role] = struct{}{}
	}
	extroles.SetRoles(rolesValue)

	return func() {
		extroles.SetRoles(nil)
	}
}

func httpGet(url string) (status int, payload string, err error) {
	client := &http.Client{}

	resp, err := invokeWithRetry(
		func() (response *http.Response, e error) {
			return client.Get(url)
		},
	)
	if err != nil {
		return 0, "", err
	}

	return handleHttpResp(resp)
}

func invokeWithRetry(invoke func() (*http.Response, error)) (*http.Response, error) {
	remainingAttempts := 10
	for {
		resp, err := invoke()
		if err == nil {
			return resp, err
		}
		remainingAttempts--
		if remainingAttempts == 0 {
			return nil, err
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func handleHttpResp(resp *http.Response) (status int, payload string, err error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, strings.ReplaceAll(string(respBytes), "\n", ""), err
}

// Require extends require with additional assertions
type Require struct {
	*require.Assertions
}

// ConnectionRefused requires that the HTTP response returns the error "connection refused"
func (r *Require) ConnectionRefused() {
	_, _, err := httpGet(fmt.Sprintf("http://localhost:%d/document/some-did", listenPort))
	r.Error(err)
	r.Contains(err.Error(), "connection refused")
}

// Response requires that the HTTP response's status and response are the expected status and response respectfully
func (r *Require) Response(basePath, namespace string, expectedStatus int, expectedResponse string) {
	const urlTemplate = "http://localhost:%d%s/%s:some-did"

	status, payload, err := httpGet(fmt.Sprintf(urlTemplate, listenPort, basePath, namespace))
	r.NoError(err)
	r.Contains(payload, expectedResponse)
	r.Equal(expectedStatus, status)
}