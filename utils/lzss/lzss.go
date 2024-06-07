package lzss

func Decompress(input []byte) (output []byte, err error) {
	offset := 0

	for offset < len(input) {
		flag := input[offset]
		offset += 1

		for bit := range 8 {
			if (flag & (1 << bit)) != 0 {
				output = append(output, input[offset])
				offset += 1
			} else {
				if offset >= len(input) {
					break
				}

				lookback_flag := (int(input[offset]) << 8) | int(input[offset+1])
				offset += 2

				lookback_length := int(lookback_flag&0xf) + 3
				lookback_offset := int(lookback_flag >> 4)

				if lookback_flag == 0 {
					break
				}

				for range lookback_length {
					loffset := len(output) - lookback_offset
					if loffset <= 0 || loffset >= len(output) {
						output = append(output, 0)
					} else {
						output = append(output, output[loffset])
					}
				}
			}
		}
	}

	return
}
