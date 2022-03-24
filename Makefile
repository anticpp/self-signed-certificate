.PHONY: ca certs runs runc clean help

.DEFAULT_GOAL := help

help:
	@echo "Usage:"
	@echo " make ca      				# Create CA key/certificate."
	@echo " make certs   				# Create server key/certificate."
	@echo ""
	@echo " make runs    				# Run s_server"
	@echo " make runs-verify    		# Run s_server, verify client certificate"
	@echo ""
	@echo " make runc    				# Run s_client, expecting succ"
	@echo " make runc-err-cert 		   	# Run s_client, epxecting certificate verifed error "
	@echo " make runc-err-cn 	     	# Run s_client, epxecting common name verifed error "
	@echo " make runc-with-cert    		# Run s_client, with client certificate"
	@echo ""
	@echo " make clean   				# Clean all generated key/certificate"
	@echo " make help    				# Help message"

ca:
	@sh makeca

certs:
	@sh makecerts

runc:
	@echo "Expecting succ"
	@sudo openssl s_client -quiet -verify_hostname test-server -verify_return_error -CAfile CA_test/cacert.pem

runc-err-cert:
	@echo "No cacert.pem is provided, expectiong certificate verified error"
	@sudo openssl s_client -quiet -verify_return_error

runc-err-cn:
	@echo "Set \"-verify_hostname\" with an error name, expecting hostname verified error"
	@sudo openssl s_client -quiet -verify_hostname error-server-name -verify_return_error -CAfile CA_test/cacert.pem

runc-with-cert:
	@echo "With client certificate"
	@sudo openssl s_client -quiet -verify_hostname test-server -verify_return_error -CAfile CA_test/cacert.pem -cert client/test-client.crt -key client/test-client.key

runs:
	@sudo openssl s_server -cert server/test-server.crt  -key server/test-server.key

runs-verify:
	@echo "Server verify client certificate ON"
	@sudo openssl s_server -quiet -Verify 1 -verify_hostname test-client -verify_return_error  -cert server/test-server.crt  -key server/test-server.key -CAfile CA_test/cacert.pem

clean:
	@sudo rm -rvf /etc/pki/CA_test
	@ rm -rfv ./CA_test ./server ./client
