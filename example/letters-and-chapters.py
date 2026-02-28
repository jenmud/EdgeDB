import json
import sys
import urllib.request


class Node:
    def __init__(self, node_id, label, properties):
        self.id = node_id
        self.label = label
        self.properties = properties

    def to_dict(self):
        return {
            "id": self.id,
            "label": self.label,
            "properties": self.properties
        }


def get_books():
    return [
        # Old Testament
        ("Genesis", "Gen", 50, "Old", 1),
        ("Exodus", "Exod", 40, "Old", 2),
        ("Leviticus", "Lev", 27, "Old", 3),
        ("Numbers", "Num", 36, "Old", 4),
        ("Deuteronomy", "Deut", 34, "Old", 5),
        ("Joshua", "Josh", 24, "Old", 6),
        ("Judges", "Judg", 21, "Old", 7),
        ("Ruth", "Ruth", 4, "Old", 8),
        ("1 Samuel", "1Sam", 31, "Old", 9),
        ("2 Samuel", "2Sam", 24, "Old", 10),
        ("1 Kings", "1Kgs", 22, "Old", 11),
        ("2 Kings", "2Kgs", 25, "Old", 12),
        ("1 Chronicles", "1Chr", 29, "Old", 13),
        ("2 Chronicles", "2Chr", 36, "Old", 14),
        ("Ezra", "Ezra", 10, "Old", 15),
        ("Nehemiah", "Neh", 13, "Old", 16),
        ("Esther", "Est", 10, "Old", 17),
        ("Job", "Job", 42, "Old", 18),
        ("Psalms", "Ps", 150, "Old", 19),
        ("Proverbs", "Prov", 31, "Old", 20),
        ("Ecclesiastes", "Eccl", 12, "Old", 21),
        ("Song of Solomon", "Song", 8, "Old", 22),
        ("Isaiah", "Isa", 66, "Old", 23),
        ("Jeremiah", "Jer", 52, "Old", 24),
        ("Lamentations", "Lam", 5, "Old", 25),
        ("Ezekiel", "Ezek", 48, "Old", 26),
        ("Daniel", "Dan", 12, "Old", 27),
        ("Hosea", "Hos", 14, "Old", 28),
        ("Joel", "Joel", 3, "Old", 29),
        ("Amos", "Amos", 9, "Old", 30),
        ("Obadiah", "Obad", 1, "Old", 31),
        ("Jonah", "Jonah", 4, "Old", 32),
        ("Micah", "Mic", 7, "Old", 33),
        ("Nahum", "Nah", 3, "Old", 34),
        ("Habakkuk", "Hab", 3, "Old", 35),
        ("Zephaniah", "Zeph", 3, "Old", 36),
        ("Haggai", "Hag", 2, "Old", 37),
        ("Zechariah", "Zech", 14, "Old", 38),
        ("Malachi", "Mal", 4, "Old", 39),

        # New Testament
        ("Matthew", "Matt", 28, "New", 40),
        ("Mark", "Mark", 16, "New", 41),
        ("Luke", "Luke", 24, "New", 42),
        ("John", "John", 21, "New", 43),
        ("Acts", "Acts", 28, "New", 44),
        ("Romans", "Rom", 16, "New", 45),
        ("1 Corinthians", "1Cor", 16, "New", 46),
        ("2 Corinthians", "2Cor", 13, "New", 47),
        ("Galatians", "Gal", 6, "New", 48),
        ("Ephesians", "Eph", 6, "New", 49),
        ("Philippians", "Phil", 4, "New", 50),
        ("Colossians", "Col", 4, "New", 51),
        ("1 Thessalonians", "1Thess", 5, "New", 52),
        ("2 Thessalonians", "2Thess", 3, "New", 53),
        ("1 Timothy", "1Tim", 6, "New", 54),
        ("2 Timothy", "2Tim", 4, "New", 55),
        ("Titus", "Titus", 3, "New", 56),
        ("Philemon", "Phlm", 1, "New", 57),
        ("Hebrews", "Heb", 13, "New", 58),
        ("James", "Jas", 5, "New", 59),
        ("1 Peter", "1Pet", 5, "New", 60),
        ("2 Peter", "2Pet", 3, "New", 61),
        ("1 John", "1John", 5, "New", 62),
        ("2 John", "2John", 1, "New", 63),
        ("3 John", "3John", 1, "New", 64),
        ("Jude", "Jude", 1, "New", 65),
        ("Revelation", "Rev", 22, "New", 66),
    ]


def main():
    if len(sys.argv) < 1:
        print("Usage: python letters-and-chapters.py [api]")
        sys.exit(1)

    api = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080/api/v1/nodes"

    books = get_books()
    nodes = []

    for name, abbr, chapters, testament, order in books:
        # Book node
        nodes.append(Node(
            0,
            "letter",
            {
                "name": name,
                "abbreviation": abbr,
                "testament": testament,
                "chapter_count": chapters,
                "order": order,
                "type": "book"
            }
        ).to_dict())

        # Chapter nodes
        for chapter in range(1, chapters + 1):
            nodes.append(Node(
                0,
                "chapter",
                {
                    "book": name,
                    "chapter": chapter,
                    "testament": testament,
                    "bookOrder": order,
                    "type": "chapter"
                }
            ).to_dict())

    print("uploading to api")
    body = {"nodes": nodes}
    payload = json.dumps(body).encode("utf-8")

    req = urllib.request.Request(
        url=api,
        data=payload,
        method="PUT",
        headers={"Content-Type": "application/json", "Content-Length": str(len(payload))},
    )

    with urllib.request.urlopen(req) as resp:
        print("Status:", resp.status)

    print(f"Done")


if __name__ == "__main__":
    main()