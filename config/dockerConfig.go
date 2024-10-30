package config

type Docker struct {
	Version string // 1.18
}

const GenHostScript = `
#!/bin/bash
CERT_PATH=/etc/docker/certs/${HOST}
CERT_PERIOD=100
function gen_tls() {
    cd ${CERT_PATH} || exit
    openssl genrsa -aes256 -passout pass:"${PASSWORD}" -out ca-key.pem 4096
    openssl req -new -x509 -days ${CERT_PERIOD} -key ca-key.pem -passin pass:"${PASSWORD}" -sha256 -out ca.pem -subj "/C=NL/ST=./L=./O=./CN=$IP"
    openssl genrsa -out server-key.pem 4096
    openssl req -subj "/CN=$IP" -sha256 -new -key server-key.pem -out server.csr
    echo subjectAltName = IP:$IP,IP:0.0.0.0 >>extfile.cnf
    echo extendedKeyUsage = serverAuth >>extfile.cnf
    openssl x509 -req -days ${CERT_PERIOD} -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -passin "pass:${PASSWORD}" -CAcreateserial -out server-cert.pem -extfile extfile.cnf
    openssl genrsa -out key.pem 4096
    openssl req -subj '/CN=client' -new -key key.pem -out client.csr
    echo extendedKeyUsage = clientAuth >>extfile.cnf
    echo extendedKeyUsage = clientAuth >extfile-client.cnf
    openssl x509 -req -days ${CERT_PERIOD} -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -passin "pass:${PASSWORD}" -CAcreateserial -out cert.pem -extfile extfile-client.cnf
    rm -f -v client.csr server.csr extfile.cnf extfile-client.cnf
    chmod -v 0400 ca-key.pem key.pem server-key.pem
    chmod -v 0444 ca.pem server-cert.pem cert.pem
}

function init_dir() {
  if [ ! -d ${CERT_PATH} ] ; then
    mkdir -p ${CERT_PATH}
    fi
}

function main() {
    echo "docker 证书生成工具,需要输入以下参数完成"
    read -p "请输入IP地址:" IP
    read -p "请输入证书密码:" PASSWORD
    read -p "请输入主机名称:" HOST
    init_dir
    gen_tls
}

main
`
