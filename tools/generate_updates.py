import glob
import hashlib
import json
import math
import os
import string

import lxml
from lxml import etree as ET
from lxml.builder import E

from eamuserc5 import EamuseRC5




def create_update_archive(filenames):
    filetable_data = bytearray()
    main_data = bytearray()

    cur_offset = 0
    for filename in filenames:
        base_filename = os.path.basename(filename)
        base_filename += " " * (32 - len(base_filename))
        filesize = os.path.getsize(filename)

        filetable_data += base_filename.encode('ascii')
        filetable_data += int.to_bytes(cur_offset, length=4, byteorder="little")
        filetable_data += int.to_bytes(filesize, length=4, byteorder="little")

        main_data += open(filename, "rb").read()

        cur_offset += filesize

    output = bytearray()

    output += int.to_bytes(len(filenames), length=4, byteorder="little")
    output += int.to_bytes(len(filetable_data), length=4, byteorder="little")
    output += int.to_bytes(0, length=4, byteorder="little")
    output += int.to_bytes(len(main_data), length=4, byteorder="little")
    output += filetable_data
    output += main_data

    return output



def replace_bytes(buffer, offset, data):
    return buffer[:offset] + bytearray(data) + buffer[offset+len(data):]


def xml2str(xml, xml_declaration=False):
    if type(xml) == lxml.etree._Element:
        xml = ET.tostring(xml, encoding="euc-jp", xml_declaration=xml_declaration)

    return xml


def generate_file_xml(path, target_game_codes):
    output_filename = os.path.join(path, "file.xml")
    os.makedirs(os.path.dirname(output_filename), exist_ok=True)

    target_files = []
    for game_code in target_game_codes:
        target_files.append("%s/file.xml" % game_code)

    file_list_info = []

    for target_file in target_files:
        target_filepath = os.path.join(path, os.path.normpath(target_file))
        assert(os.path.exists(target_filepath))

        file_xml = open(target_filepath, "rb").read()

        info = E.list(
            path=target_file,
            size=str(len(file_xml)),
            sum=hashlib.md5(file_xml).hexdigest().upper(),
        )

        file_list_info.append(info)

    xml = E.datasync(*file_list_info)

    with open(output_filename, "wb") as outfile:
        outfile.write(xml2str(xml))

    return [output_filename]


def generate_sum_file_xml(path):
    output_filename = os.path.join(path, "sum", "file.xml")
    os.makedirs(os.path.dirname(output_filename), exist_ok=True)

    main_file_xml = open(os.path.join(path, "file.xml"), "rb").read()

    xml = E.datasync(
        E.list(
            path="file.xml",
            size=str(len(main_file_xml)),
            sum=hashlib.md5(main_file_xml).hexdigest().upper(),
        )
    )

    with open(output_filename, "wb") as outfile:
        outfile.write(xml2str(xml))

    return [output_filename]


def generate_net_id(net_id_info):
    netids_count = len(net_id_info['song_ids'])
    netids_offset = 0x18
    netids_len = netids_count * 2
    checksum_len = 0x10

    output = bytearray()
    output += int.to_bytes(net_id_info['version'][0], 1, 'little')
    output += int.to_bytes(net_id_info['version'][1], 1, 'little')
    output += int.to_bytes(netids_count, 2, 'little')
    output += int.to_bytes(netids_offset, 2, 'little') # Offset to net IDs
    output += int.to_bytes(netids_len, 2, 'little') # Size of net ID table
    output += int.to_bytes(checksum_len, 2, 'little') # Size of checksum?
    output += int.to_bytes(netids_offset + netids_len + checksum_len, 2, 'little') # End of file

    output += int.to_bytes(net_id_info['timestamp']['year'] - 2000, 1, 'little')
    output += int.to_bytes(net_id_info['timestamp']['month'], 1, 'little')
    output += int.to_bytes(net_id_info['timestamp']['day'], 1, 'little')
    output += int.to_bytes(net_id_info['timestamp']['hour'], 1, 'little')
    output += int.to_bytes(net_id_info['timestamp']['minute'], 1, 'little')
    output += int.to_bytes(net_id_info['timestamp']['second'], 1, 'little')

    padding = netids_offset - len(output)
    if padding > 0:
        output += bytearray([0] * padding)

    for net_id in net_id_info['song_ids']:
        output += int.to_bytes(net_id, 2, 'little')

    output += hashlib.md5(output).digest()
    return output


def generate_gamelist_file_xml(path, game_id, output_path):
    output_filenames = []

    output_filename = os.path.join(output_path, game_id, "file.xml")
    os.makedirs(os.path.dirname(output_filename), exist_ok=True)

    net_id_data = generate_net_id(
        json.load(open(os.path.join(path, game_id, "net_id.json")))
    )

    netid_path = os.path.join(output_path, game_id, "net_id.bin")
    with open(netid_path, "wb") as outfile:
        outfile.write(net_id_data)

    output_filenames.append(netid_path)

    target_files = [
        "%s/net_id.bin" % game_id,
        "%s/list.xml" % game_id,
    ]

    file_list_info = []

    for target_file in target_files:
        target_filepath = os.path.join(output_path, os.path.normpath(target_file))
        print(target_filepath)

        if not os.path.exists(target_filepath):
            continue

        file_xml = open(target_filepath, "rb").read()

        info = E.list(
            path=target_file,
            size=str(len(file_xml)),
            sum=hashlib.md5(file_xml).hexdigest().upper(),
        )

        file_list_info.append(info)

    xml = E.datasync(*file_list_info)

    with open(output_filename, "wb") as outfile:
        outfile.write(xml2str(xml))

    output_filenames.append(output_filename)

    return output_filenames


def generate_gamelist_list_xml(path, game_id, output_path):
    net_ids = json.load(open(os.path.join(path, game_id, "net_id.json")))

    output_filenames = []

    output_list_filename = os.path.join(output_path, game_id, "list.xml")
    os.makedirs(os.path.dirname(output_list_filename), exist_ok=True)

    songid_folders = [x for x in glob.glob(os.path.join(path, game_id, "*")) if os.path.isdir(x)]

    xml_info = []

    for songid_folder in songid_folders:
        song_id = "".join([x for x in os.path.basename(songid_folder) if x in string.digits])

        data_ver = "%d.%02d" % (net_ids['data_ver'][0], net_ids['data_ver'][1])

        fcn_data = create_update_archive(sorted(glob.glob(os.path.join(songid_folder, "*"))))

        # open("tmp_%s.fcn" % song_id, "wb").write(fcn_data)

        from io import BytesIO

        part_idx = 0
        song_info = []

        CHUNK_SIZE = 0x30000

        assert(math.ceil(len(fcn_data) / CHUNK_SIZE) <= 0x1f)

        for i in range(0, len(fcn_data), CHUNK_SIZE):
            path = "%s/%s_%06x" % (game_id, song_id, part_idx)
            enc_filename = os.path.join(output_path, path)

            os.makedirs(os.path.dirname(enc_filename), exist_ok=True)

            with BytesIO(fcn_data[i:i+CHUNK_SIZE]) as input_file:
                with open(enc_filename, "wb") as output_file:
                    EamuseRC5.encrypt_file(input_file, output_file, os.path.basename(path))
                    # output_file.write(input_file.read())
                    output_filenames.append(enc_filename)

            with open(enc_filename, "rb") as infile:
                file_data = infile.read()
                filesize = len(file_data)
                checksum = hashlib.md5(file_data).hexdigest()

            data_xml = E.data(
                path=path,
                size=str(filesize),
                sum=checksum.upper(),
                ver=data_ver,
                music_id=song_id
            )

            print(xml2str(data_xml, False).decode('ascii'))

            song_info.append(data_xml)

            part_idx += 1

        output_filename = os.path.join(output_path, game_id, "i%s.xml" % song_id)
        xml = E.info(*song_info)
        with open(output_filename, "wb") as outfile:
            outfile.write(xml2str(xml))
        output_filenames.append(output_filename)

        with open(output_filename, "rb") as infile:
            file_data = infile.read()
            filesize = len(file_data)
            checksum = hashlib.md5(file_data).hexdigest()

        xml_info.append(E.data(
            path=os.path.join(game_id, "i%s.xml" % song_id),
            size=str(filesize),
            sum=checksum.upper(),
            ver=data_ver,
            music_id=song_id
        ))
        print(output_filename)

    output_filename = os.path.join(output_path, game_id, "list.xml")

    xml = E.info(*xml_info)
    with open(output_filename, "wb") as outfile:
        outfile.write(xml2str(xml))
    print(output_filename)

    output_filenames.append(output_filename)

    return output_filenames



# Add the game codes to generate updates for here
target_game_codes = [
    "D09",
    "D10",
    "D39",
    "D40",
]

generated_files = []
for game_code in target_game_codes:
    generated_files += generate_gamelist_list_xml("updates_work", game_code, "updates")
    generated_files += generate_gamelist_file_xml("updates_work", game_code, "updates")

generated_files += generate_file_xml("updates", target_game_codes)
generated_files += generate_sum_file_xml("updates")

print(generated_files)

filemap = {
    "files": [],
}
for filename in generated_files:
    prefix = os.path.abspath("updates")
    path = os.path.abspath(filename)

    if path.startswith(prefix):
        path = path[len(prefix):]

    path = path.lstrip("/")

    filemap['files'].append({
        "path": path,
    })

# Dump file for server to be able to route the static files only
json.dump(filemap, open("static.json", "w"))
