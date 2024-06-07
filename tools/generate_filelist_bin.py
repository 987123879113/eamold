import glob
import hashlib
import os
import sys

def get_padded_string(input, padlen=16):
    output = input.encode('ascii')
    while (len(output) % padlen) != 0:
        output += b"\0"
    return output

filelist_data = bytearray()
for path in sorted(glob.glob(os.path.join(sys.argv[1], "*"))):
    filename = os.path.basename(path).lower()

    if filename in ["filelist.bin", "data_ver.bin"]:
        continue

    with open(path, "rb") as infile:
        data = infile.read()
        filelen = len(data)
        md5 = hashlib.md5(data).digest()

    print(path, filename, md5, filelen)

    filelist_data += get_padded_string(filename[:16])
    filelist_data += int.to_bytes(filelen, 4, 'little')
    filelist_data += md5

filelist_data += get_padded_string("filelist.bin")
filelist_data += int.to_bytes(len(filelist_data) + 0x10 + 4, 4, 'little')
filelist_data += hashlib.md5(filelist_data).digest()

open(os.path.join(sys.argv[1], "filelist.bin"), "wb").write(filelist_data)
