.PHONY: ca certs runs runc clean help

.DEFAULT_GOAL := help

help:
	@echo "Usage: make ca|certs|runc|runs|clean|help"
	@echo " make ca      				# Create CA key/certificate."
	@echo " make certs   				# Create server key/certificate."
	@echo " make runs    				# Run s_server"
	@echo " make runc    				# Run s_client, expecting succ"
	@echo " make runc-err-cert 		   	# Run s_client, epxecting certificate verifed error "
	@echo " make runc-err-cn 	     	# Run s_client, epxecting common name verifed error "
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

runs:
	@sudo openssl s_server -cert server/test-server.crt  -key server/test-server.key

clean:
	@sudo rm -rvf /etc/pki/CA_test
	@ rm -rfv ./CA_test ./server
