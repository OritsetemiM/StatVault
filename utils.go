package main

func normalizeText(s string) string {
	replacements := map[rune]rune{
		'Š': 'S', 'š': 's', 'Ž': 'Z', 'ž': 'z',
		'Č': 'C', 'č': 'c', 'Ć': 'C', 'ć': 'c',
		'Đ': 'D', 'đ': 'd', 'Ñ': 'N', 'ñ': 'n',
		'Ü': 'U', 'ü': 'u', 'Ö': 'O', 'ö': 'o',
		'Ä': 'A', 'ä': 'a', 'É': 'E', 'é': 'e',
		'È': 'E', 'è': 'e', 'Ê': 'E', 'ê': 'e',
		'Á': 'A', 'á': 'a', 'À': 'A', 'à': 'a',
		'Í': 'I', 'í': 'i', 'Ó': 'O', 'ó': 'o',
		'Ú': 'U', 'ú': 'u', 'Ý': 'Y', 'ý': 'y',
		'Ø': 'O', 'ø': 'o', 'Å': 'A', 'å': 'a',
		'Æ': 'A', 'æ': 'a', 'Ğ': 'G', 'ğ': 'g',
		'Ş': 'S', 'ş': 's', 'İ': 'I', 'ı': 'i',
		'Ā': 'A', 'ā': 'a', 'Ē': 'E', 'ē': 'e',
		'Ī': 'I', 'ī': 'i', 'Ō': 'O', 'ō': 'o',
		'Ū': 'U', 'ū': 'u', 'Ņ': 'N', 'ņ': 'n',
		'Ķ': 'K', 'ķ': 'k', 'Ļ': 'L', 'ļ': 'l',
		'Ģ': 'G', 'ģ': 'g',
	}

	result := make([]rune, 0, len(s))
	for _, r := range s {
		if rep, ok := replacements[r]; ok {
			result = append(result, rep)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}