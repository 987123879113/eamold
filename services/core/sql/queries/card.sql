-- name: GetCardNumberStatus :one
SELECT COUNT(*)
FROM core_assigned_card_numbers
WHERE label = ?
AND number = ?
LIMIT 1;

-- name: AddUsedCardNumber :exec
INSERT INTO core_assigned_card_numbers
(label, number)
VALUES
(?, ?);
