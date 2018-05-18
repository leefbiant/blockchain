curl -XPOST 127.0.0.1:8000/nodes/register -d '{"nodes":["http://127.0.0.1:8000/chain", "http://127.0.0.1:8001/chain"]}'
curl -XPOST 127.0.0.1:8001/nodes/register -d '{"nodes":["http://127.0.0.1:8000/chain", "http://127.0.0.1:8001/chain"]}'

curl -XPOST 127.0.0.1:8000/txion -d '{"from":"leef", "to":"biant", "amount":1}' | python -m json.tool
curl -XPOST 127.0.0.1:8001/txion -d '{"from":"leef", "to":"biant", "amount":2}' | python -m json.tool
curl -XPOST 127.0.0.1:8001/txion -d '{"from":"leef", "to":"biant", "amount":3}' | python -m json.tool
curl -XPOST 127.0.0.1:8000/mine  | python -m json.tool
curl -XPOST 127.0.0.1:8001/mine  | python -m json.tool
curl 127.0.0.1:8000/chain | python -m json.tool
curl 127.0.0.1:8001/chain | python -m json.tool



