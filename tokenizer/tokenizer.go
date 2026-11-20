package tokenizer

func Tokenizing(data []byte) [][]byte {
	var next byte
	var lines [][]byte
	var flag bool
	var line []byte
	for j := 0; j < len(data); j++ {
		flag = false
		current := data[j]
		if j+1 < len(data) {
			next = data[j+1]
			flag = true
		}
		if current >= 0x30 && current <= 0x39 {
			line = append(line, current)
			continue
		} else if current >= 0x61 && current <= 0x7A {
			line = append(line, current)
			continue
		} else if current >= 0x41 && current <= 0x5A {
			line = append(line, current+0x20)
			continue
		}
		if flag == true && current == 0xD0 && next >= 0xB0 && next <= 0xBF {
			j++
			line = append(line, current)
			line = append(line, next)
			continue
		} else if flag == true && current == 0xD1 && next >= 0x80 && next <= 0x8F {
			j++
			line = append(line, current)
			line = append(line, next)
			continue
		} else if flag == true && current == 0xD0 && next >= 0x90 && next <= 0x9F {
			j++
			next += 32
			line = append(line, current)
			line = append(line, next)
			continue
		} else if flag == true && current == 0xD0 && next >= 0xA0 && next <= 0xAF {
			j++
			current += 1
			next -= 32
			line = append(line, current)
			line = append(line, next)
			continue
		} else if (flag == true && current == 0xD0 && next == 0x81) || (flag == true && current == 0xD1 && next == 0x91) {
			j++
			current = 0xD1
			next = 0x91
			line = append(line, current)
			line = append(line, next)
			continue
		}
		if len(line) > 0 {
			lines = append(lines, line)
			line = nil
		}
	}
	if len(line) > 0 {
		lines = append(lines, line)
	}
	return lines
}

func SplitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i := 0; i <= len(data); i++ {
		if i == len(data) || data[i] == '\n' {
			line := data[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			line = stripSpace(line)
			lines = append(lines, line)
			start = i + 1
		}
	}
	return lines
}

func stripSpace(b []byte) []byte {
	i, j := 0, len(b)
	for i < j && (b[i] == ' ' || b[i] == '\t') {
		i++
	}
	for j > i && (b[j-1] == ' ' || b[j-1] == '\t') {
		j--
	}
	return b[i:j]
}
