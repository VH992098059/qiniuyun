from flask import Flask, request, jsonify
from paddleocr import PaddleOCR
import os
import uuid
import json

app = Flask(__name__)
ocr = PaddleOCR(
    use_doc_orientation_classify=False,
    use_doc_unwarping=False,
    use_textline_orientation=False
)

OUTPUT_DIR = "output"
os.makedirs(OUTPUT_DIR, exist_ok=True)

@app.route('/ocr', methods=['POST'])
def run_ocr():
    if 'image' not in request.files:
        return jsonify({'error': 'No image uploaded'}), 400

    image_file = request.files['image']
    filename = f"new_orc.png"
    image_path = os.path.join(OUTPUT_DIR, filename)
    image_file.save(image_path)

    result = ocr.predict(input=image_path)

    json_results = []
    for idx, res in enumerate(result):
        json_path = os.path.join(OUTPUT_DIR, f"{'new_orc'}_res_{idx}.json")
        res.save_to_img(os.path.join(OUTPUT_DIR, f"{'new_orc'}_res_{idx}.jpg"))
        res.save_to_json(json_path)

        with open(json_path, 'r', encoding='utf-8') as f:
            json_data = json.load(f)
            json_results.append(json_data)

    return jsonify(json_results)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)