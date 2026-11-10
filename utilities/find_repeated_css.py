import re
from pathlib import Path
from collections import Counter, defaultdict

import tinycss2

root = Path(__file__).resolve().parents[1] / "front" / "css"
files = sorted(root.rglob("*.css"))

selector_counts = Counter()
selector_files = defaultdict(set)
block_counts = Counter()
property_counts = Counter()


def norm_selector(selector: str) -> str:
    return re.sub(r"\s+", " ", selector).strip().lower()


for css_file in files:
    text = css_file.read_text(encoding="utf-8", errors="ignore")
    stylesheet = tinycss2.parse_stylesheet(text, skip_comments=True, skip_whitespace=True)

    def walk(rules):
        for rule in rules:
            at = getattr(rule, "at_keyword", None)
            if at:
                if at in {"media", "supports", "container", "layer", "scope"}:
                    if rule.content:
                        walk(rule.content)
                elif at == "keyframes":
                    continue
                else:
                    if rule.content:
                        walk(rule.content)
                continue

            if getattr(rule, "type", None) not in {"qualified-rule", "style-rule"}:
                continue

            selector_text = tinycss2.serialize(rule.prelude).strip()
            if not selector_text:
                continue

            selectors = [s.strip() for s in selector_text.split(",") if s.strip()]
            if not selectors:
                continue

            declaration_pairs = []
            decl_tokens = tinycss2.parse_declaration_list(rule.content)
            for token in decl_tokens:
                if getattr(token, "type", None) == "declaration":
                    name = token.name.strip().lower()
                    value = tinycss2.serialize(token.value).strip().lower()
                    if name and value:
                        declaration_pairs.append((name, value))
                        property_counts[(name, value)] += 1

            if declaration_pairs:
                decl_text = "; ".join(f"{name}: {value}" for name, value in declaration_pairs)
                block_text = f"{selector_text} {{ {decl_text}; }}"
                block_counts[block_text.lower()] += 1

            for selector in selectors:
                sel = norm_selector(selector)
                selector_counts[sel] += 1
                selector_files[sel].add(str(css_file.relative_to(root.parent.parent)))

    walk(stylesheet)

print(f"Scanned {len(files)} CSS files under {root}")
print("\nRepeated selectors across files or multiple occurrences:")
repeated = [(sel, count, sorted(selector_files[sel])) for sel, count in selector_counts.items() if count > 1]
for sel, count, file_list in sorted(repeated, key=lambda x: (-x[1], x[0]))[:40]:
    print(f"{count:>2}x  {sel}")
    print(f"     files: {', '.join(file_list[:10])}")

print("\nRepeated full rule blocks:")
found_block = False
for block_key, count in sorted(block_counts.items(), key=lambda x: (-x[1], x[0])):
    if count <= 1:
        continue
    found_block = True
    print(f"{count:>2}x  {block_key[:220]}")
    break
if not found_block:
    print("No repeated full rule blocks found.")

print("\nMost repeated property/value pairs:")
for (prop, value), count in sorted(property_counts.items(), key=lambda x: (-x[1], x[0]))[:30]:
    if count > 1:
        print(f"{count:>2}x  {prop}: {value}")
