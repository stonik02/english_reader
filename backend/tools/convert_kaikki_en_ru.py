#!/usr/bin/env python3
"""Convert Kaikki/Wiktextract English data to the backend dictionary JSONL.

The source archive is read line-by-line, so its 2.6 GB gzip file is never
unpacked and never loaded entirely into memory. The result is compatible with
`make dictionary-import`.
"""

from __future__ import annotations

import argparse
import gzip
import json
import re
import sys
import unicodedata
from collections import defaultdict
from pathlib import Path
from typing import Any, Iterator, TextIO
from urllib.parse import quote


ATTRIBUTION = "Wiktionary contributors, CC BY-SA 4.0"
SOURCE_URL = "https://en.wiktionary.org/wiki/{}"
PARTS_OF_SPEECH = {
    "noun": "noun",
    "name": "proper noun",
    "verb": "verb",
    "adj": "adjective",
    "adjective": "adjective",
    "adv": "adverb",
    "adverb": "adverb",
    "pron": "pronoun",
    "pronoun": "pronoun",
    "prep": "preposition",
    "preposition": "preposition",
    "conj": "conjunction",
    "conjunction": "conjunction",
    "intj": "interjection",
    "interjection": "interjection",
    "det": "determiner",
    "determiner": "determiner",
    "num": "numeral",
    "numeral": "numeral",
    "article": "article",
    "phrase": "phrase",
    "prep_phrase": "prepositional phrase",
    "particle": "particle",
    "postp": "postposition",
    "prefix": "prefix",
    "suffix": "suffix",
    "contraction": "contraction",
    "proverb": "proverb",
    "symbol": "symbol",
    "punct": "punctuation",
}


def arguments() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", required=True, type=Path, help="Kaikki raw .jsonl or .jsonl.gz file")
    parser.add_argument("--output", required=True, type=Path, help="Target JSONL for dictionary-import")
    parser.add_argument("--version", required=True, help="Immutable source version, for example 2026-08-05")
    parser.add_argument("--max-input-lines", type=int, default=0, help="Process only N input lines; useful for a dry preview")
    parser.add_argument("--include-phrases", action="store_true", help="Keep multi-word English entries")
    return parser.parse_args()


def open_source(path: Path) -> TextIO:
    if path.suffix == ".gz":
        return gzip.open(path, "rt", encoding="utf-8", errors="replace")
    return path.open("r", encoding="utf-8", errors="replace")


def valid_word(word: Any, include_phrases: bool) -> bool:
    if not isinstance(word, str) or not word or len(word) > 128:
        return False
    if not include_phrases and any(character.isspace() for character in word):
        return False
    return all(character.isalpha() or character in "-'" for character in word)


def normalize_russian(value: str) -> str:
    """Remove Wiktionary's optional acute/grave stress marks from Russian text.

    We remove only U+0301 and U+0300, rather than every combining mark: that
    preserves letters such as `й`, whose Unicode decomposition includes a
    combining breve.
    """
    value = unicodedata.normalize("NFC", value)
    return unicodedata.normalize("NFC", value.replace("\u0301", "").replace("\u0300", ""))


def first_example(sense: dict[str, Any]) -> str:
    for example in sense.get("examples", []):
        text = example.get("text")
        if isinstance(text, str) and text:
            return text
    return ""


def sense_key(sense: dict[str, Any], index: int) -> str:
    glosses = sense.get("glosses")
    if isinstance(glosses, list) and glosses and isinstance(glosses[0], str):
        return glosses[0]
    return f"__unscoped_{index}"


def translation_groups(entry: dict[str, Any]) -> dict[str, list[str]]:
    groups: dict[str, list[str]] = defaultdict(list)
    for translation in entry.get("translations", []):
        if not isinstance(translation, dict) or translation.get("lang_code") != "ru":
            continue
        word = translation.get("word")
        if not isinstance(word, str) or not word:
            continue
        word = normalize_russian(word)
        key = translation.get("sense")
        if not isinstance(key, str) or not key:
            key = "__unscoped"
        if word not in groups[key]:
            groups[key].append(word)
    return groups


def terms(value: str) -> set[str]:
    return set(re.findall(r"[a-z]{3,}", value.lower()))


def source_sense(label: str, senses: list[dict[str, Any]], fallback_index: int) -> dict[str, Any]:
    """Find the closest Wiktionary gloss for Kaikki's translation label.

    Kaikki's `translation.sense` is normally a short paraphrase, while the
    sense gloss is a full definition. They often do not match literally. A
    lexical-overlap match keeps examples where possible; order is a stable
    fallback so valid Russian translations are never silently discarded.
    """
    if label == "__unscoped":
        return senses[0]

    label_terms = terms(label)
    best_index = 0
    best_score = 0.0
    for index, sense in enumerate(senses):
        glosses = sense.get("glosses", [])
        text = " ".join(gloss for gloss in glosses if isinstance(gloss, str))
        if label == text or label in text or text in label:
            return sense
        gloss_terms = terms(text)
        if not label_terms or not gloss_terms:
            continue
        score = len(label_terms & gloss_terms) / len(label_terms | gloss_terms)
        if score > best_score:
            best_index = index
            best_score = score
    if best_score > 0:
        return senses[best_index]
    return senses[min(fallback_index, len(senses) - 1)]


def records(entry: dict[str, Any], version: str, include_phrases: bool) -> Iterator[dict[str, Any]]:
    part_of_speech = PARTS_OF_SPEECH.get(entry.get("pos"))
    if entry.get("lang_code") != "en" or part_of_speech is None:
        return
    word = entry.get("word")
    if not valid_word(word, include_phrases):
        return
    groups = translation_groups(entry)
    if not groups:
        return
    senses = [value for value in entry.get("senses", []) if isinstance(value, dict)]
    if not senses:
        return

    position = 0
    for index, (label, translations) in enumerate(groups.items()):
        sense = source_sense(label, senses, index)
        yield {
            "lemma": word.lower(),
            "language": "en",
            "part_of_speech": part_of_speech,
            "translations": translations,
            "example_en": first_example(sense),
            "example_ru": "",
            "source_url": SOURCE_URL.format(quote(word.replace(" ", "_"), safe="")),
            "attribution": ATTRIBUTION,
            "position": position,
            "source_version": version,
        }
        position += 1


def main() -> int:
    args = arguments()
    if not args.input.is_file():
        print(f"input file does not exist: {args.input}", file=sys.stderr)
        return 2
    args.output.parent.mkdir(parents=True, exist_ok=True)

    input_lines = 0
    malformed = 0
    written = 0
    with open_source(args.input) as source, args.output.open("w", encoding="utf-8") as target:
        for line in source:
            if args.max_input_lines and input_lines >= args.max_input_lines:
                break
            input_lines += 1
            try:
                entry = json.loads(line)
            except json.JSONDecodeError:
                malformed += 1
                continue
            if not isinstance(entry, dict):
                continue
            for record in records(entry, args.version, args.include_phrases):
                target.write(json.dumps(record, ensure_ascii=False, separators=(",", ":")) + "\n")
                written += 1
            if input_lines % 100_000 == 0:
                print(f"processed={input_lines} written={written} malformed={malformed}", file=sys.stderr)

    print(f"completed: processed={input_lines} written={written} malformed={malformed}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
