package converter

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	azureAuthProvider = "azure"
	cfgClientID       = "client-id"
	cfgApiserverID    = "apiserver-id"
	cfgTenantID       = "tenant-id"
	cfgEnvironment    = "environment"
	cfgConfigMode     = "config-mode"

	argClientID     = "--client-id"
	argServerID     = "--server-id"
	argTenantID     = "--tenant-id"
	argEnvironment  = "--environment"
	argClientSecret = "--client-secret"
	argIsLegacy     = "--legacy"
	argUsername     = "--username"
	argPassword     = "--password"
	argLoginMethod  = "--login"

	flagClientID     = "client-id"
	flagServerID     = "server-id"
	flagTenantID     = "tenant-id"
	flagEnvironment  = "environment"
	flagClientSecret = "client-secret"
	flagIsLegacy     = "legacy"
	flagUsername     = "username"
	flagPassword     = "password"
	flagLoginMethod  = "login"

	execName        = "kubelogin"
	getTokenCommand = "get-token"
	execAPIVersion  = "client.authentication.k8s.io/v1beta1"
)

func Convert(o Options) error {
	config, err := o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return fmt.Errorf("unable to load kubeconfig: %s", err)
	}

	for _, authInfo := range config.AuthInfos {
		if authInfo != nil {
			if authInfo.AuthProvider == nil || authInfo.AuthProvider.Name != azureAuthProvider {
				continue
			}
			exec := &api.ExecConfig{
				Command: execName,
				Args: []string{
					getTokenCommand,
				},
				APIVersion: execAPIVersion,
			}
			if o.isSet(flagEnvironment) {
				exec.Args = append(exec.Args, argEnvironment)
				exec.Args = append(exec.Args, o.TokenOptions.Environment)
			} else if authInfo.AuthProvider.Config[cfgEnvironment] != "" {
				exec.Args = append(exec.Args, argEnvironment)
				exec.Args = append(exec.Args, authInfo.AuthProvider.Config[cfgEnvironment])
			}
			if o.isSet(flagServerID) {
				exec.Args = append(exec.Args, argServerID)
				exec.Args = append(exec.Args, o.TokenOptions.ServerID)
			} else if authInfo.AuthProvider.Config[cfgApiserverID] != "" {
				exec.Args = append(exec.Args, argServerID)
				exec.Args = append(exec.Args, authInfo.AuthProvider.Config[cfgApiserverID])
			}
			if o.isSet(flagClientID) {
				exec.Args = append(exec.Args, argClientID)
				exec.Args = append(exec.Args, o.TokenOptions.ClientID)
			} else if authInfo.AuthProvider.Config[cfgClientID] != "" {
				exec.Args = append(exec.Args, argClientID)
				exec.Args = append(exec.Args, authInfo.AuthProvider.Config[cfgClientID])
			}
			if o.isSet(flagClientID) {
				exec.Args = append(exec.Args, argTenantID)
				exec.Args = append(exec.Args, o.TokenOptions.TenantID)
			} else if authInfo.AuthProvider.Config[cfgClientID] != "" {
				exec.Args = append(exec.Args, argTenantID)
				exec.Args = append(exec.Args, authInfo.AuthProvider.Config[cfgTenantID])
			}
			if o.isSet(flagIsLegacy) && o.TokenOptions.IsLegacy {
				exec.Args = append(exec.Args, argIsLegacy)
			} else if authInfo.AuthProvider.Config[cfgConfigMode] == "" || authInfo.AuthProvider.Config[cfgConfigMode] == "0" {
				exec.Args = append(exec.Args, argIsLegacy)
			}
			if o.isSet(flagClientSecret) {
				exec.Args = append(exec.Args, argClientSecret)
				exec.Args = append(exec.Args, o.TokenOptions.ClientSecret)
			}
			if o.isSet(flagUsername) {
				exec.Args = append(exec.Args, argUsername)
				exec.Args = append(exec.Args, o.TokenOptions.Username)
			}
			if o.isSet(flagPassword) {
				exec.Args = append(exec.Args, argPassword)
				exec.Args = append(exec.Args, o.TokenOptions.Password)
			}
			if o.isSet(flagLoginMethod) {
				exec.Args = append(exec.Args, argLoginMethod)
				exec.Args = append(exec.Args, o.TokenOptions.LoginMethod)
			}
			authInfo.Exec = exec
			authInfo.AuthProvider = nil
		}
	}

	clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)

	return nil
}
