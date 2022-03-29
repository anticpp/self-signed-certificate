.PHONY: ca certs \
		runs runs-verify \
		runc-no-verify runc-verify-no-ca runc-verify-with-ca runc-verify-host-err runc-verify-succ runc-with-cert \
		clean help

.DEFAULT_GOAL := help

help:
	@echo "Usage:"
	@echo " make ca      				# Create CA key/certificate."
	@echo " make certs   				# Create server key/certificate."
	@echo ""
	@echo " make runs    				# Run s_server"
	@echo " make runs-verify    		# Run s_server, verify client certificate"
	@echo ""
	@echo " make runc-no-verify    	    # Run s_client, no verify"
	@echo " make runc-verify-no-ca    	# Run s_client, no cafile is provided."
	@echo " make runc-verify-with-ca    # Run s_client, cafile is provided."
	@echo " make runc-verify-host-err   # Run s_client, hostname mismatched."
	@echo " make runc-verify-succ       # Run s_client, certificate verification succ"
	@echo " make runc-with-cert    		# Run s_client, with client certificate"
	@echo ""
	@echo " make clean   				# Clean all generated key/certificate"
	@echo " make help    				# Help message"

ca:
	@sh makeca

certs:
	@sh makecerts

# Run server.
runs:
	@echo "Run s_server."
	@sudo openssl s_server -cert server/test-server.pem  -key server/test-server-key.pem

# Run server, with client certificate verification ON.
runs-verify:
	@echo "Run s_server, client verification ON."
	@sudo openssl s_server -quiet -Verify 1 -verify_hostname test-client -verify_return_error -cert server/test-server.pem -key server/test-server-key.pem -CAfile CA_test/cacert.pem


# Test client without verification, the TLS connection will still go on with encrypted transport.
runc-no-verify:
	@echo "Run s_client"
	@echo "CAfile...             [NO]"
	@echo "Verify return error...[NO]"
	@echo "Verify hostname...    [NO]"
	@echo ""
	@sudo openssl s_client -quiet

# Test client with verification, the TLS connection will fail because no CAfile is provided.
runc-verify-no-ca:
	@echo "Run s_client"
	@echo "CAfile...             [NO]"
	@echo "Verify return error...[YES]"
	@echo "Verify hostname...    [NO]"
	@echo ""
	@sudo openssl s_client -quiet -verify_return_error

# Test client with verification, the TLS connection will succ because the CAfile is provided.
runc-verify-with-ca:
	@echo "Run s_client"
	@echo "CAfile...             [YES]"
	@echo "Verify return error...[YES]"
	@echo "Verify hostname...    [NO]"
	@echo ""
	@sudo openssl s_client -quiet -verify_return_error -CAfile CA_test/cacert.pem

# Test client with verification, the TLS connection will succ because the hostname is mismatched with server certificate's CN/DNS.
runc-verify-host-err:
	@echo "Run s_client"
	@echo "CAfile...             [YES]"
	@echo "Verify return error...[YES]"
	@echo "Verify hostname...    [YES] ... with error name"
	@echo ""
	@sudo openssl s_client -quiet -verify_hostname error-server-name -verify_return_error -CAfile CA_test/cacert.pem

# Test success, client open TLS connection on server, with server certificate verified success.
runc-verify-succ:
	@echo "Run s_client"
	@echo "CAfile...             [YES]"
	@echo "Verify return error...[YES]"
	@echo "Verify hostname...    [YES]"
	@echo ""
	@sudo openssl s_client -quiet -verify_hostname test-server -verify_return_error -CAfile CA_test/cacert.pem

# Test client with client certificate/key, for server side verification.
runc-with-cert:
	@echo "Run s_client"
	@echo "CAfile...             [YES]"
	@echo "Verify return error...[YES]"
	@echo "Verify hostname...    [YES]"
	@echo "Use client certificate... [YES]"
	@echo ""
	@sudo openssl s_client -quiet -verify_hostname test-server -verify_return_error -CAfile CA_test/cacert.pem -cert client/test-client.pem -key client/test-client-key.pem

clean:
	@sudo rm -rvf /etc/pki/CA_test
	@rm -rfv ./CA_test ./server ./client
