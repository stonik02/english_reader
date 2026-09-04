package morphology

import "strings"

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) Lemma(word string) (string, bool) {
	candidates := s.Candidates(word)
	if len(candidates) == 1 {
		return word, false
	}
	return candidates[1], true
}

func (s *Service) Candidates(word string) []string {
	result := []string{word}
	if lemma, ok := irregular[word]; ok {
		return appendCandidate(result, lemma)
	}
	if strings.HasSuffix(word, "ies") && len(word) > 3 {
		return appendCandidate(result, strings.TrimSuffix(word, "ies")+"y")
	}
	if strings.HasSuffix(word, "ing") && len(word) > 4 {
		return inflectionCandidates(result, strings.TrimSuffix(word, "ing"))
	}
	if strings.HasSuffix(word, "ed") && len(word) > 3 {
		return inflectionCandidates(result, strings.TrimSuffix(word, "ed"))
	}
	if strings.HasSuffix(word, "es") && len(word) > 3 {
		result = appendCandidate(result, strings.TrimSuffix(word, "es"))
	}
	if strings.HasSuffix(word, "s") && len(word) > 2 {
		result = appendCandidate(result, strings.TrimSuffix(word, "s"))
	}
	return result
}

func inflectionCandidates(result []string, stem string) []string {
	if len(stem) > 2 && stem[len(stem)-1] == stem[len(stem)-2] {
		result = appendCandidate(result, stem[:len(stem)-1])
	}
	result = appendCandidate(result, stem)
	if strings.HasSuffix(stem, "i") {
		result = appendCandidate(result, strings.TrimSuffix(stem, "i")+"y")
	}
	return appendCandidate(result, stem+"e")
}

func appendCandidate(values []string, value string) []string {
	for _, current := range values {
		if current == value {
			return values
		}
	}
	return append(values, value)
}

var irregular = map[string]string{
	"am": "be", "are": "be", "been": "be", "did": "do", "done": "do",
	"is": "be", "was": "be", "were": "be", "has": "have", "had": "have",
	"arose": "arise", "arisen": "arise", "awoke": "awake", "awoken": "awake",
	"became": "become", "began": "begin", "begun": "begin", "bent": "bend",
	"bit": "bite", "bitten": "bite", "bled": "bleed", "blew": "blow", "blown": "blow",
	"broke": "break", "broken": "break", "brought": "bring", "built": "build", "bought": "buy",
	"caught": "catch", "chose": "choose", "chosen": "choose", "came": "come", "crept": "creep",
	"dealt": "deal", "dug": "dig", "drew": "draw", "drawn": "draw", "drank": "drink", "drunk": "drink",
	"drove": "drive", "driven": "drive", "ate": "eat", "eaten": "eat", "fed": "feed", "fell": "fall",
	"felt": "feel", "flew": "fly", "flown": "fly", "found": "find", "forgot": "forget", "forgotten": "forget",
	"forgave": "forgive", "forgiven": "forgive", "froze": "freeze", "frozen": "freeze", "gave": "give", "given": "give",
	"gone": "go", "got": "get", "grew": "grow", "grown": "grow", "heard": "hear", "held": "hold",
	"hid": "hide", "hidden": "hide", "hit": "hit", "hurt": "hurt", "kept": "keep", "knew": "know", "known": "know",
	"laid": "lay", "lay": "lie", "lain": "lie", "led": "lead", "left": "leave", "lent": "lend", "let": "let",
	"lost": "lose", "made": "make", "meant": "mean", "met": "meet", "paid": "pay", "put": "put", "ran": "run",
	"read": "read", "rode": "ride", "ridden": "ride", "rang": "ring", "rung": "ring", "rose": "rise", "risen": "rise",
	"said": "say", "saw": "see", "seen": "see", "sold": "sell", "sent": "send", "shook": "shake", "shaken": "shake",
	"shot": "shoot", "shut": "shut", "sang": "sing", "sung": "sing", "sank": "sink", "sunk": "sink", "slept": "sleep",
	"spoke": "speak", "spoken": "speak", "spent": "spend", "stood": "stand", "stole": "steal", "stolen": "steal",
	"stuck": "stick", "swam": "swim", "swum": "swim", "swept": "sweep", "swore": "swear", "sworn": "swear",
	"taught": "teach", "tore": "tear", "torn": "tear", "thought": "think", "threw": "throw", "thrown": "throw",
	"took": "take", "told": "tell", "understood": "understand", "went": "go", "woke": "wake", "woken": "wake",
	"wore": "wear", "worn": "wear", "won": "win", "wrote": "write", "written": "write",
	"children": "child", "mice": "mouse", "better": "good",
}
