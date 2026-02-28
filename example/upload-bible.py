import requests
import json
import time

KJV_URL = "https://raw.githubusercontent.com/thiagobodruk/bible/master/json/en_kjv.json"
API_URL = "http://localhost:8080/api/v1/nodes"

BATCH_SIZE = 500


def download_bible():
    print("Downloading KJV Bible...")
    response = requests.get(KJV_URL)
    response.raise_for_status()

    # Remove UTF-8 BOM safely
    text = response.content.decode("utf-8-sig")
    return json.loads(text)


def generate_nodes(bible_json):
    node_id = 0
    for book_order, book in enumerate(bible_json, start=1):
        book_name = book["name"]
        testament = "Old" if book_order <= 39 else "New"

        # Book node
        yield {
            "id": node_id,
            "label": "letter",
            "properties": {
                "name": book_name,
                "testament": testament,
                "order": book_order,
                "chapter_count": len(book["chapters"]),
                "type": "book",
            },
        }
        node_id += 1

        for chapter_index, chapter in enumerate(book["chapters"], start=1):
            # Chapter node
            yield {
                "id": node_id,
                "label": "chapter",
                "properties": {
                    "book": book_name,
                    "chapter": chapter_index,
                    "testament": testament,
                    "bookOrder": book_order,
                    "verse_count": len(chapter),
                    "reference": f"{book_name} {chapter_index}",
                    "type": "chapter",
                },
            }
            node_id += 1

            for verse_index, verse_text in enumerate(chapter, start=1):
                yield {
                    "id": node_id,
                    "label": "verse",
                    "properties": {
                        "book": book_name,
                        "chapter": chapter_index,
                        "verse": verse_index,
                        "testament": testament,
                        "bookOrder": book_order,
                        "reference": f"{book_name} {chapter_index}:{verse_index}",
                        "content": verse_text,
                        "word_count": len(verse_text.split()),
                        "char_count": len(verse_text),
                        "type": "verse",
                    },
                }
                node_id += 1


def upload_batch(batch):
    payload = {"nodes": batch}

    response = requests.put(
        API_URL,
        headers={"Content-Type": "application/json"},
        data=json.dumps(payload),
        timeout=60,
    )

    if response.status_code >= 400:
        print("Upload failed:", response.status_code)
        print(response.text)
        raise Exception("Upload failed")

    print(f"Uploaded batch of {len(batch)} nodes")


def upload_in_batches(node_generator):
    batch = []

    for node in node_generator:
        batch.append(node)

        if len(batch) >= BATCH_SIZE:
            upload_batch(batch)
            batch.clear()
            time.sleep(0.2)  # small pause to avoid overwhelming server

    if batch:
        upload_batch(batch)


def main():
    bible_json = download_bible()
    node_generator = generate_nodes(bible_json)
    upload_in_batches(node_generator)
    print("Done uploading full Bible.")


if __name__ == "__main__":
    main()