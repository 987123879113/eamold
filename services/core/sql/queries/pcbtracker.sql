-- name: GetPcbidStatus :one
SELECT status
FROM core_pcbid
WHERE pcbid = ?
LIMIT 1;
