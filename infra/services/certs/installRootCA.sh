CA_NAME="knh.cloud"
# CERT_FILE="$CA_NAME.crt"
CERT_FILE="RootCA.crt"
CERT_INSTALL_DIR="/usr/local/share/ca-certificates"
CERT_PATH="${CERT_INSTALL_DIR}/${CERT_FILE}"

#1. Generate root CA certificate
sudo openssl s_client -connect ${CA_NAME}:443 < /dev/null > /tmp/temporary.out
sudo mkdir -p "$CERT_INSTALL_DIR"
sudo openssl x509 -outform PEM < /tmp/temporary.out > "$CERT_PATH"
rm /tmp/temporary.out

#2. Install the root CA certificate
sudo ln -s "$CERT_PATH" "/etc/ssl/certs/$CERT_FILE"
sudo chmod -x "$CERT_PATH"

#3. Add the root CA certificate to the system's trust store
# Generate the hash
HASH="$(openssl x509 -hash -noout -in $CERT_PATH).0"
#Note: If another cert has the same hash use suffix .1 or .2 instead of .0.
# Display the hash value
echo "$HASH"
# Link the hash to the certificate
sudo ln -s "/etc/ssl/certs/$CERT_FILE" "/etc/ssl/certs/$HASH"
ls -al "/etc/ssl/certs/$HASH"
