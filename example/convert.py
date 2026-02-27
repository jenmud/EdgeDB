import csv
import json
import sys
import urllib.request
from io import StringIO


def main():
    nodes = []

    if len(sys.argv) < 2:
        print("Usage: python transform_commandments.py <raw_csv_url> [api]")
        sys.exit(1)

    raw_url = sys.argv[1]
    api = sys.argv[2] if len(sys.argv) > 2 else "http://localhost:8080/api/v1/nodes"

    # Fetch CSV from URL
    try:
        with urllib.request.urlopen(raw_url) as response:
            csv_data = response.read().decode("utf-8")
    except Exception as e:
        print(f"Failed to fetch CSV: {e}")
        sys.exit(1)

    reader = csv.DictReader(StringIO(csv_data))

    ident = 1
    for row in reader:
        # All other columns become JSON properties
        props = {
            k: v for k, v in row.items()
            if k != "commandment_number"
        }

        nodes.append({"id": ident, "label": "commandment", "properties": props})
        ident += 1


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
        resp_body = resp.read().decode("utf-8")
        print("Status:", resp.status)

    print(f"Done")


if __name__ == "__main__":
    main()