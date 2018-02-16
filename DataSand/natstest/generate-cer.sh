export NATS=/home/saichler/nats/gnatsd-v1.0.4-linux-amd64
export HOST=saichler
openssl genrsa -aes256 -out $NATS/ca-key.pem 4096
openssl req -new -x509 -days 365 -key $NATS/ca-key.pem -sha256 -out $NATS/ca.pem
openssl genrsa -out $NATS/server-key.pem 4096
openssl req -subj "/CN=$HOST" -sha256 -new -key $NATS/server-key.pem -out $NATS/server.csr
echo subjectAltName = DNS:$HOST,IP:10.157.157.6,IP:127.0.0.1 >> extfile.cnf
echo extendedKeyUsage = serverAuth >> $NATS/extfile.cnf
openssl x509 -req -days 365 -sha256 -in $NATS/server.csr -CA $NATS/ca.pem -CAkey $NATS/ca-key.pem \
	  -CAcreateserial -out $NATS/server-cert.pem -extfile $NATS/extfile.cnf
