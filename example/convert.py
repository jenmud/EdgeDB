import csv
import json
import sys
import urllib.request
from io import StringIO


def main():
    if len(sys.argv) < 2:
        print("Usage: python transform_commandments.py <raw_csv_url> [output_file]")
        sys.exit(1)

    raw_url = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else "transformed.csv"

    # Fetch CSV from URL
    try:
        with urllib.request.urlopen(raw_url) as response:
            csv_data = response.read().decode("utf-8")
    except Exception as e:
        print(f"Failed to fetch CSV: {e}")
        sys.exit(1)

    reader = csv.DictReader(StringIO(csv_data))

    with open(output_file, "w", newline="", encoding="utf-8") as f_out:
        writer = csv.writer(f_out)

        # Header (adjust/remove if your importer doesn't want one)
        writer.writerow(["id", "label", "properties"])

        ident = 1
        for row in reader:
            # All other columns become JSON properties
            props = {
                k: v for k, v in row.items()
                if k != "commandment_number"
            }

            props_json = json.dumps(props, ensure_ascii=False)

            writer.writerow([ident, "command", props_json])
            ident += 1

    print(f"Done. Wrote transformed CSV to: {output_file}")


if __name__ == "__main__":
    main()