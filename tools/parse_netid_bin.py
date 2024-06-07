import json
import sys

with open(sys.argv[1], "rb") as infile:
    data = bytearray(infile.read())

output = {
    "version": [int(data[0]), int(data[1])],
}

netids_count = int.from_bytes(data[2:4], 'little')
netids_offset = int.from_bytes(data[4:6], 'little')
netids_len = int.from_bytes(data[6:8], 'little')
checksum_len = int.from_bytes(data[8:10], 'little')
full_file_len = int.from_bytes(data[10:12], 'little')

timestamp_values = [int(x) for x in data[12:12+6]]
output['timestamp'] = {
    'year': timestamp_values[0] + 2000,
    'month': timestamp_values[1],
    'day': timestamp_values[2],
    'hour': timestamp_values[3],
    'minute': timestamp_values[4],
    'second': timestamp_values[5],
}

output['song_ids'] = []
for i in range(netids_offset, netids_offset + netids_count * 2, 2):
    output['song_ids'].append(int.from_bytes(data[i:i+2], 'little'))

print(json.dumps(output, indent=4))
