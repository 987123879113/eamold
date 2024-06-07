import hashlib

from collections import deque


class EamuseRand:
    def __init__(self):
        self.seed = 0


    def srand(self, val):
        self.seed = val


    def rand(self):
        self.seed = (self.seed * 0x41c64e6d + 0x3039) & 0xffffffff
        return (self.seed >> 16) & 0x7fff


class EamuseRC5(object):
    @staticmethod
    def _generate_key(k):
        key = bytearray("#2-ngEn5000nEngA10?nOuChIwo9rAbUreBa-uEyEmOMorIn5nO5oTo97rI@\0\0\0\0".encode('ascii'))

        if type(k) is str:
            k = k.encode('ascii')

        k = bytearray(k)

        krand = EamuseRand()
        krand.srand(len(k))
        r = krand.rand()

        for i in range(len(k)):
            idx = (i + r) % len(key)
            key[idx] = (key[idx] + k[i]) & 0xff

        return bytearray(hashlib.md5(key).digest())


    @staticmethod
    def _update_counter(expanded_key, counter):
        expanded_key[counter[0] & 0x0f] = (expanded_key[counter[0] & 0x0f] + expanded_key[(counter[0] + 1) & 0x0f]) & 0xffffffff
        counter[0] += 1


    @staticmethod
    def _rotate_left(val, r_bits, max_bits):
        v1 = (val << r_bits % max_bits) & (2 ** max_bits - 1)
        v2 = ((val & (2 ** max_bits - 1)) >> (max_bits - (r_bits % max_bits)))
        return v1 | v2


    @staticmethod
    def _rotate_right(val, r_bits, max_bits):
        v1 = ((val & (2 ** max_bits - 1)) >> r_bits % max_bits)
        v2 = (val << (max_bits - (r_bits % max_bits)) & (2 ** max_bits - 1))
        return v1 | v2


    @staticmethod
    def _encrypt_block(data, expanded_key, blocksize, counter):
        w = blocksize // 2
        b = blocksize // 8
        mod = 2 ** w

        A = int.from_bytes(data[:b // 2], byteorder='little')
        B = int.from_bytes(data[b // 2:], byteorder='little')

        A = (A + expanded_key[counter[0] & 0xf]) % mod
        B = (B + expanded_key[(counter[0] + 1) & 0xf]) % mod

        for _ in range(3):
            A = (EamuseRC5._rotate_left((A ^ B), B, w) + expanded_key[counter[0] & 0xf]) % mod
            B = (EamuseRC5._rotate_left((A ^ B), A, w) + expanded_key[(counter[0] + 1) & 0xf]) % mod

        res = A.to_bytes(b // 2, byteorder='little') + B.to_bytes(b // 2, byteorder='little')

        EamuseRC5._update_counter(expanded_key, counter)

        return bytearray(res)


    @staticmethod
    def _decrypt_block(data, expanded_key, blocksize, counter):
        w = blocksize // 2
        b = blocksize // 8
        mod = 2 ** w

        A = int.from_bytes(data[:b // 2], byteorder='little')
        B = int.from_bytes(data[b // 2:], byteorder='little')

        for _ in range(3):
            B = EamuseRC5._rotate_right(B - expanded_key[(counter[0] + 1) & 0xf], A, w) ^ A
            A = EamuseRC5._rotate_right(A - expanded_key[counter[0] & 0xf], B, w) ^ B

        B = (B - expanded_key[(counter[0] + 1) & 0xf]) % mod
        A = (A - expanded_key[counter[0] & 0xf]) % mod

        res = A.to_bytes(b // 2, byteorder='little') + B.to_bytes(b // 2, byteorder='little')

        EamuseRC5._update_counter(expanded_key, counter)

        return bytearray(res)


    @staticmethod
    def _expand_key(key, wordsize, rounds):
        def _align_key(key, align_val):
            key2 = key.ljust(len(key) + (align_val - (len(key) % (align_val))), b'\x00')
            L = [int.from_bytes(key[i:i + align_val], byteorder='little') for i in range(0, len(key), align_val)]
            return L


        def _extend_key(w, r):
            P, Q = (0xB7E15163, 0x9E3779B9)
            t = 2 * (r + 1)

            S = [P]
            for i in range(1, t):
                S.append((S[i - 1] + Q) % 2 ** w)

            return S


        def _mix(L, S, r, w, c):
            t = 2 * (r + 1)
            m = max(c, t)
            A = B = i = j = 0

            for k in range(m):
                for _ in range(3):
                    A = S[i] = EamuseRC5._rotate_left(S[i] + A + B, 3, w)
                    B = L[j] = EamuseRC5._rotate_left(L[j] + A + B, L[j] + A + B, w)

                i = (i + 1) % t
                j = (j + 1) % c

            return S

        aligned = _align_key(key, wordsize // 8)
        extended = _extend_key(wordsize, rounds)
        S = _mix(aligned, extended, rounds, wordsize, len(aligned))

        return S


    @staticmethod
    def encrypt_file(infile, outfile, key, blocksize=64, rounds=7):
        key = EamuseRC5._generate_key(key)

        rand = EamuseRand()
        w = blocksize // 2
        b = blocksize // 8

        infile.seek(0, 2)
        file_len = infile.tell()
        infile.seek(0)

        expanded_key = EamuseRC5._expand_key(key, w, rounds)

        counter = [sum(expanded_key) & 0xffffffff]
        rand.srand(counter[0])

        # Generate IV based on counter and file length
        iv = bytearray([ord('r'), ord('C'), ord('?'), ord('0')])
        iv += bytearray(int.to_bytes(file_len, 4, 'little'))
        iv = deque(iv)
        iv.rotate(counter[0] & 7)
        iv = bytearray(iv)

        iv = bytearray(map(lambda x: (x + rand.rand()) & 0xff, iv))
        iv = EamuseRC5._encrypt_block(
            iv,
            expanded_key,
            blocksize,
            counter
        )
        outfile.write(iv)

        while True:
            chunk = infile.read(b)

            if not chunk:
                break

            encrypted_chunk = EamuseRC5._encrypt_block(
                chunk.ljust(b, b'\x00'),
                expanded_key,
                blocksize,
                counter
            )

            outfile.write(encrypted_chunk)
            iv = encrypted_chunk


    @staticmethod
    def decrypt_file(infile, outfile, key, blocksize=64, rounds=7):
        key = EamuseRC5._generate_key(key)

        rand = EamuseRand()
        w = blocksize // 2
        b = blocksize // 8

        expanded_key = EamuseRC5._expand_key(key, w, rounds)
        counter = [sum(expanded_key) & 0xffffffff]
        rand.srand(counter[0])

        # Decrypt IV and pull file size from it
        iv = infile.read(b)
        dec_iv = EamuseRC5._decrypt_block(
            iv,
            expanded_key,
            blocksize,
            counter
        )

        dec_iv = map(lambda x: (x - rand.rand()) & 0xff, dec_iv)
        dec_iv = deque(dec_iv)
        dec_iv.rotate(-(counter[0] & 7)+1)
        dec_iv = bytearray(dec_iv)

        if dec_iv[:4] != b"rC?0":
            print("Invalid IV, can't decrypt")
            return

        output_len = int.from_bytes(dec_iv[4:8], 'little')

        while True:
            chunk = infile.read(b)

            if not chunk:
                break

            decrypted_chunk = EamuseRC5._decrypt_block(
                chunk,
                expanded_key,
                blocksize,
                counter
            )

            iv = chunk

            if output_len - len(decrypted_chunk) < 0:
                decrypted_chunk = decrypted_chunk[:output_len]

            output_len -= len(decrypted_chunk)
            outfile.write(decrypted_chunk)
