if [[ -z "${TF_VERSION}" ]]; then
  echo `Installing Terraform version ${TF_VERSION}`
  TF_VERSION="1.0.0"
fi

cd /usr/local/bin
curl https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip -o terraform_${TF_VERSION}_linux_amd64.zip
unzip terraform_${TF_VERSION}_linux_amd64.zip
rm terraform_${TF_VERSION}_linux_amd64.zip