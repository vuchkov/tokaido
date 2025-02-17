/*

`tls` is a third-generation TLS certificate library for Tokaido.

It is responsible for generating a local Certificate Authority which, on macOS, it can
then install into the trusted keychain automatically.

This CA is then used to sign the *.local.tokaido.io wildcard certificate for the Proxy
server and also to sign the 'haproxy' certificate for each Tokaido instance.

*/

package tls

// ConfigureTLS is the principal entry point into the TLS library, and is responsible for the
// entire TLS security workflow
func ConfigureTLS() (err error) {

	// Remove up any legacy SSL certificates from previous versions
	removeLegacyCertificate()

	// Generate a new Certificate Authority
	createCA()

	// Sign a new Certificate for *.local.tokaido.io
	createWildcardCertificate()
	return nil

}

// SignCertificate creates a singular certificate request based only on the CN supplied
// it then saves the key and certificate to disk at ~/.tok/tls/{project}/{cn}.[crt|key]
// while also returning both as byte slices
func SignCertificate(cn string, sans []string) (key, cert []byte, err error) {
	return nil, nil, nil
}
