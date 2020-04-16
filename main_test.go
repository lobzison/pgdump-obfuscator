package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var testConfig01 = &Configuration{
	Obfuscations: []TargetedObfuscation{
		TargetedObfuscation{
			Target{Table: "auth_user", Column: "email"},
			ScrambleEmail,
		},
		TargetedObfuscation{
			Target{Table: "auth_user", Column: "password"},
			GenScrambleBytes(7),
		},
		TargetedObfuscation{
			Target{Table: "accounts_profile", Column: "phone"},
			ScrambleDigits,
		},
	},
}

const testInput01 = `
--

SELECT pg_catalog.setval('auth_user_id_seq', 123111, true);


--
-- Data for Name: auth_user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY auth_user (id, username, first_name, last_name, email, password, is_staff, is_active, is_superuser, last_login, date_joined) FROM stdin;
123123111	499964777.sdsad	testname	testname	\N	!	f	t	f	2011-02-07 12:08:30+00	2010-11-22 19:27:12.31832+00
333114441	testT1@bing.com			testT1@bing.com	!	f	t	f	2011-06-08 12:57:36+00	2011-06-08 12:50:25.206298+00
515131311	whoisthere			noreply781134796251@bing.com	pbkdf2_sha256$10000$qweqweqweqwe$cThxOHE4	f	t	f	2012-11-16 18:27:43.673889+00	2012-11-16 18:27:43.229281+00
\.

COPY accounts_profile (id, user_id, opted_in, next_break, status, phone, last_visited, come_from, cs_letter, city_id, budget_range_id, prefs_opt) FROM stdin;
6161	12113	f	\N	0	+74991002000	2011-07-04 12:28:33.895325+00	\N	f	\N	\N	\N
1223	1321	f	\N	0	666666666	2011-09-28 09:37:20.83051+00	\N	f	\N	\N	\N
4423	55512	f	\N	0		\N	\N	f	\N	\N	\N
\.
`

func assertScramble(t *testing.T, scramble func([]byte) []byte, s, expected string) {
	assertString(t, string(scramble([]byte(s))), expected)
}

func assertString(t *testing.T, s, expected string) {
	if s != expected {
		t.Fatalf("Expected: '%s'(%d) received: '%s'(%d)", expected, len(expected), s, len(s))
	}
}

func TestProcess01(t *testing.T) {
	input := bufio.NewReader(bytes.NewBufferString(testInput01))
	output := new(bytes.Buffer)
	process(testConfig01, input, output)
	outString := output.String()
	if outString == testInput01 {
		t.Fatal("Outputs are equal")
	}
	if !strings.Contains(outString, "COPY auth_user (id, username, first_name, last_name, email, password, is_staff, is_active, is_superuser, last_login, date_joined) FROM stdin;") ||
		!strings.Contains(outString, "COPY accounts_profile (id, user_id, opted_in, next_break, status, phone, last_visited, come_from, cs_letter, city_id, budget_range_id, prefs_opt) FROM stdin;") {
		t.Fatal("Changed SQL")
	}
	if strings.Contains(outString, "pbkdf2_sha256$10000$qweqweqweqwe$cThxOHE4") ||
		strings.Contains(outString, "+3801445223001") {
		t.Fatal("Did not scramble sensitive data")
	}
	if !strings.Contains(outString, "515131311	whoisthere") ||
		!strings.Contains(outString, `	2011-07-04 12:28:33.895325+00	\N	f	\N	\N	\N`) ||
		!strings.Contains(outString, `1223	1321	f	\N	0	`) {
		t.Fatal("Changed other data")
	}
	lines := strings.Split(outString, "\n")
	if !strings.HasPrefix(lines[11], "123123111") {
		t.Fatal("Line 12 invalid:", lines[11])
	}
	fields := strings.Split(lines[11], "\t")
	assertString(t, fields[4], "\\N")
}

func TestProcess02(t *testing.T) {
	testConfig02 := &Configuration{
		Obfuscations: []TargetedObfuscation{
			TargetedObfuscation{
				Target{Table: "with_emails", Column: "id"},
				ScrambleDigits,
			},
			TargetedObfuscation{
				Target{Table: "with_emails", Column: "emails_list"},
				ScrambleEmail,
			},
		},
	}
	const testInput02 = `COPY with_emails (id, emails_list) FROM stdin;
41	e1@hbo.com,info@o2.co.uk
52	e3@mail.ru,e4@gmail.com
\.
`
	const expected = `COPY with_emails (id, emails_list) FROM stdin;
09	7f@example.com
02	M7@example.com
\.
`
	input := bufio.NewReader(bytes.NewBufferString(testInput02))
	output := new(bytes.Buffer)
	Salt = []byte("test-salt")
	process(testConfig02, input, output)
	assertString(t, output.String(), expected)
}

func TestScrambleBytes(t *testing.T) {
	Salt = []byte("test-salt")
	assertScramble(t, ScrambleBytes, "everyone lies", "oSE0Sm0yioFSJ")
	assertScramble(t, ScrambleBytes, "very long line very long line very long line very",
		"4ce6EsWcmziuUzpEtV0rGiZAOtiHprwB0wWWWuOYrHkqHQtAN")
	assertScramble(t, ScrambleBytes, "{item1,\"item space 2\"}",
		"{yho3y,rEZwPM7FVuVf1S}")
}

func TestScrambleBytesUtf8(t *testing.T) {
	Salt = []byte("test-salt")
	// Output must be of same length as input
	assertScramble(t, ScrambleBytes, "also русский and 你好",
		"emEY0UP-gkC2kV+J6pK")
	assertScramble(t, ScrambleBytes, "{\"array z\",руки,你好}",
		"{LoKXy6kRZ,uefS,G1}")
}

func TestScrambleDigits(t *testing.T) {
	Salt = []byte("test-salt")
	assertScramble(t, ScrambleDigits, "+7(876) 123-0011 или 99999999999;",
		"+1(584) 047-9250 или 22280031035;")
	assertScramble(t, ScrambleDigits, "5", "7")
}

func TestScrambleEmail(t *testing.T) {
	Salt = []byte("test-salt")
	assertScramble(t, ScrambleEmail, "solar.sultan@emerginspaceagency.com",
		"lxtTUsvMGzRo@example.com")
	assertScramble(t, ScrambleEmail, "{foo@bar.com,test@example.com}",
		"{DK3@example.com,LDVR@example.com}")
	assertScramble(t, ScrambleEmail, "унеун@mail.ru", "gfpFV@example.com")
	assertScramble(t, ScrambleEmail, "multiple@emails.com,but.not@array.in",
		"+QWUPnIS@example.com")
}

func TestScrambleUniqueEmail(t *testing.T) {
	Salt = []byte("test-salt")
	assertScramble(t, ScrambleUniqueEmail, "solar.sultan@emerginspaceagency.com",
		"lxtTUsvMGzRo@uoievjjkleb8gh7a2ayyie.example")
	assertScramble(t, ScrambleUniqueEmail, "{foo@bar.com,test@example.com}",
		"{DK3@1u7t913.example,LDVR@1o4wthaadn4.example}")
	assertScramble(t, ScrambleUniqueEmail, "унеун@mail.ru", "gfpFV@2ps2nuk.example")
	assertScramble(t, ScrambleUniqueEmail, "multiple@emails.com,but.not@array.in",
		"+QWUPnIS@7x8f15oletv2wrbhq8mcriros04.example")
}

func TestScrambleBindUrl(t *testing.T) {
	Salt = []byte("test-salt")
	assertScramble(t, ScrambleBindUrls,
		"http://bind.com?carrier=safeco&some_gid=2495330c-5d-afdb1f0845f2e9f943f7f6",
		"https://example.com?quote_gid=DJv1RQElgiRH79Gu5oOBHwml1kUglJIJqL")
	assertScramble(t, ScrambleBindUrls, "", "")
	assertScramble(t, ScrambleBindUrls, "https://example.com", "https://example.com")
}

func TestScrambleInet(t *testing.T) {
	Salt = []byte("test-salt")
	assertScramble(t, ScrambleInet, "142.34.56.78", "56.42.246.77")
	assertScramble(t, ScrambleInet, "97.34.0.18", "e4ed:d550:209d:9f10:f690:953:5d4f:c0d6")
}

func TestScrambleIBAN(t *testing.T) {
	assertScramble(t, ScrambleIBAN,
		"DE1234567890",
		"DE75512108001245126199")
	assertScramble(t, ScrambleIBAN, "", "DE75512108001245126199")
	assertScramble(t, ScrambleIBAN, "qwd8989ajn", "DE75512108001245126199")
}

func BenchmarkProcessShort(b *testing.B) {
	b.StopTimer()
	config := &Configuration{
		Obfuscations: []TargetedObfuscation{
			TargetedObfuscation{
				Target{Table: "simple", Column: "id"},
				ScrambleDigits,
			},
			TargetedObfuscation{
				Target{Table: "simple", Column: "email"},
				ScrambleEmail,
			},
			TargetedObfuscation{
				Target{Table: "simple", Column: "password"},
				ScrambleBytes,
			},
		},
	}
	const s = `--
-- Useless comments
--

select 'and other statements';

create table simple (
	id serial not null,
	email text not null,
	password text not null
);


COPY simple (id, email, password) FROM stdin;
13	e1zz@hbo.com	12345
27	e3@mail.ru	password
28	wha@who.net	strongPassw0rd
33	admin@ad.co.uk	allyourbase
41	martin.martin@drop.tv	belongto
42	xxx@strong.net	usususus
43	e1@hbo.com	12345
44	e3qq@mail.ru	password
48	wha@who.net	strongPassw0rd
49	admin@ad.co.uk	allyourbase
50	martin.martin@drop.tv	belongto
121	xxx@strong.net	usususus
122	e1www@hbo.com	12345
123	e3gggggg@mail.ru	password
124	wha@who.net	strongPassw0rd
125	admin@ad.co.uk	allyourbase
126	martin.martin@drop.tv	belongto
127	xxx@strong.net	usususus
\.
`
	var input *bufio.Reader
	var output *bytes.Buffer
	for i := 0; i < b.N; i++ {
		input = bufio.NewReader(bytes.NewBufferString(s))
		output = new(bytes.Buffer)
		b.StartTimer()
		process(config, input, output)
		b.StopTimer()
	}
}

func BenchmarkScrambleBytes(b *testing.B) {
	Salt = []byte("test-salt")
	s := []byte("everybody lies many times")
	for i := 0; i < b.N; i++ {
		ScrambleBytes(s)
	}
}

func BenchmarkScrambleBytesArray(b *testing.B) {
	Salt = []byte("test-salt")
	s := []byte("{everybody,lies,\"many many many\",times}")
	for i := 0; i < b.N; i++ {
		ScrambleBytes(s)
	}
}

func BenchmarkScrambleDigits(b *testing.B) {
	Salt = []byte("test-salt")
	digits := "+7(876) 123-0011"
	for i := 0; i < b.N; i++ {
		ScrambleDigits([]byte(digits))
	}
}

func BenchmarkScrambleEmail(b *testing.B) {
	Salt = []byte("test-salt")
	email := "igor.igorev@igorooking.mx1.uk.gb"
	for i := 0; i < b.N; i++ {
		ScrambleEmail([]byte(email))
	}
}

func BenchmarkScrambleEmailArray(b *testing.B) {
	Salt = []byte("test-salt")
	email := "{admin@slapmode.de,trisha@moodnight.com.co,wilf@wolf.herztner.jp}"
	for i := 0; i < b.N; i++ {
		ScrambleEmail([]byte(email))
	}
}

func BenchmarkScrambleInet(b *testing.B) {
	Salt = []byte("test-salt")
	s := []byte("23.11.1.239")
	for i := 0; i < b.N; i++ {
		ScrambleInet(s)
	}
}
