#curl --location https://api.ionos.com/databases/quota \
#  -H "Authorization: Bearer $IONOS_TOKEN" | jq
#
#curl --location https://dns.de-fra.ionos.com/quota \
#  -H "Authorization: Bearer $IONOS_TOKEN" | jq


curl -v localhost:8080/quotas -H "Authorization: Bearer $IONOS_TOKEN" | jq