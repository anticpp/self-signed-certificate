.PHONY: ca certs runs runc clean help

.DEFAULT_GOAL := help

help:
	@echo "Usage: make ca|certs|runc|runs|clean|help"
	@echo " make ca      # Create CA key/certificate."
	@echo " make certs   # Create server key/certificate."
	@echo " make runs    # Run s_server"
	@echo " make runc    # Run s_client"
	@echo " make clean   # Clean all generated key/certificate"
	@echo " make help    # Help message"

ca:
	@sh makeca

certs:
	@sh makecerts

runc:
	@sudo openssl s_client -quiet -verify_return_error -CAfile CA_test/cacert.pem

runs:
	@sudo openssl s_server -cert server/test-server.crt  -key server/test-server.key

clean:
	@sudo rm -rvf /etc/pki/CA_test
	@ rm -rfv ./CA_test ./server
