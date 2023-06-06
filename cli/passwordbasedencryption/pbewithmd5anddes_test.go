package passwordbasedencryption

import "testing"

func TestBasicDecryption(t *testing.T) {
	cases := []struct {
		ciphertext, password, plaintext string
	}{
		{"u6ccN+pf88NQFo0p2W5HUgoJXW/iGZPt", "password", "plaintext"},
		{"nWUp2auqbcKucN6VBYkL8sQtYwyFc6dXjLLJjOhR4WTKS1XfMdmx0kkYBiD4sVDycSH1Vp5JDXqDLg74PSBQ8j5k5Ongvel2", "password", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{"TgLG/fANuEVycFMO6Ap7eA==", "password", ""},
		{"Wt9vfiouLnMHPEcSBx2ZUYpVYcSrmR9O1IAt7768VbK1DH5tZe3A2YNyqdHA0dLma3Hlwe3WeU4Ba32+RLG5dIH7KUrLlZH9", "password", "ðƏ kwɪk braʊn fƊks dʒʊmptƏʊvƏ ðƏ lɛɪzi: dƊgz"},
		{"inZQMiY+UsI5HLLifuvV2HxBhoj3nNNA", "g9Q95=yNVt7E?a+nDN=%", "plaintext"},
		{"1uurVxPzTV5KGuL1ZupT+e+K57KhfDdGjV/Ej+zWvZrajf5B/KfyoGBSiE3qSYX5iIZoPO/XIIFplaAtPwAI1eWsWx4NFHWM", "g9Q95=yNVt7E?a+nDN=%", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{"ygsi6PB2b6RcOIJeiFAcIg==", "g9Q95=yNVt7E?a+nDN=%", ""},
		{"4v7gZN8/e20qX7Nm5EVbRs84zZ7IkWt+GNi8q+4dETeJodVONdoF7jaXBl8qialZ5KIlvlDD04idlAVjqiY6H/HDxkWBcyTE", "g9Q95=yNVt7E?a+nDN=%", "ðƏ kwɪk braʊn fƊks dʒʊmptƏʊvƏ ðƏ lɛɪzi: dƊgz"},
	}
	for _, c := range cases {
		got, err := Decrypt(c.password, 1000, c.ciphertext)
		if err != nil {
			t.Errorf("Got error %q for password %q, ciphertext %q", err.Error(), c.password, c.ciphertext)
		}
		if got != c.plaintext {
			t.Errorf("Decrypt(%q, 1000, %q) == %q, want %q", c.password, c.ciphertext, got, c.plaintext)
		}
	}
}

func TestBasicEncryption(t *testing.T) {
	cases := []struct {
		plaintext  string
		password   string
		iterations int
	}{
		{"plaintext", "password", 1000},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "password", 1000},
		{"", "password", 1000},
		{"ðƏ kwɪk braʊn fƊks dʒʊmptƏʊvƏ ðƏ lɛɪzi: dƊgz", "password", 1000},
		{"plaintext", "g9Q95=yNVt7E?a+nDN=%", 1000},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "g9Q95=yNVt7E?a+nDN=%", 1000},
		{"", "g9Q95=yNVt7E?a+nDN=%", 1000},
		{"ðƏ kwɪk braʊn fƊks dʒʊmptƏʊvƏ ðƏ lɛɪzi: dƊgz", "g9Q95=yNVt7E?a+nDN=%", 1000},
		{"plaintext", "password", 5},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "password", 5},
		{"", "password", 5},
		{"ðƏ kwɪk braʊn fƊks dʒʊmptƏʊvƏ ðƏ lɛɪzi: dƊgz", "password", 5},
		{"plaintext", "g9Q95=yNVt7E?a+nDN=%", 5},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "g9Q95=yNVt7E?a+nDN=%", 5},
		{"", "g9Q95=yNVt7E?a+nDN=%", 5},
		{"ðƏ kwɪk braʊn fƊks dʒʊmptƏʊvƏ ðƏ lɛɪzi: dƊgz", "g9Q95=yNVt7E?a+nDN=%", 5},
	}
	for _, c := range cases {
		ciphertext, err := Encrypt(c.password, c.iterations, c.plaintext)
		if err != nil {
			t.Errorf("Got error %q for password %q, plaintext %q", err.Error(), c.password, c.plaintext)
		}

		plaintext, _ := Decrypt(c.password, c.iterations, ciphertext)
		if plaintext != c.plaintext {
			t.Errorf("Got %q, want %q", plaintext, c.plaintext)
		}
	}
}

func TestEncryptWithFixedSalt(t *testing.T) {
	cases := []struct {
		plaintext, password, fixedsalt string
		iterations                     int
	}{
		{"plaintext", "password", "fixed_salt", 1000},
		{"encryption test", "password", "fixed_salt", 1000},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "SoMePaSsWoRd", "FixedSalt", 1000},
		{"àéïûõç", "BfRK4TnM1zYj30amLjb3", "bCi@*5tX9Van", 1000},
		{"", "TO72&BjDpUYa", "u0@5#4Yj9LxI", 1000},
	}
	for _, c := range cases {
		ciphertext, err := EncryptWithFixedSalt(c.password, c.iterations, c.plaintext, c.fixedsalt)
		if err != nil {
			t.Errorf("Got error %q for password %q, plaintext %q, salt %q", err.Error(), c.password, c.plaintext, c.fixedsalt)
		}

		plaintext, _ := DecryptWithFixedSalt(c.password, c.iterations, ciphertext, c.fixedsalt)
		if plaintext != c.plaintext {
			t.Errorf("Got %q, expected %q", plaintext, c.plaintext)
		}
	}
}

func TestDecryptWithFixedSalt(t *testing.T) {
	cases := []struct {
		ciphertext, password, plaintext, fixedsalt string
	}{
		{"IcszAY8NRJf6ANt152Fifg==", "password", "encryption test", "fixed_salt"},
	}
	for _, c := range cases {
		got, err := DecryptWithFixedSalt(c.password, 1000, c.ciphertext, c.fixedsalt)
		if err != nil {
			t.Errorf("Got error %q for password %q, ciphertext %q, salt %q", err.Error(), c.password, c.ciphertext, c.fixedsalt)
		}
		if got != c.plaintext {
			t.Errorf("Decrypt(%q, 1000, %q, %q) == %q, want %q", c.password, c.ciphertext, c.fixedsalt, got, c.plaintext)
		}
	}
}
