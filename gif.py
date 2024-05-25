import os
import logging
import sys
from PIL import Image

logging.basicConfig(level=logging.DEBUG, stream=sys.stdout)

NUM_IMAGES = 36
INDEX_LENGTH = 2

def make_gif(uuid, image_path):
    gif_folder_path = os.path.join(image_path, uuid, "gif")
    img_folder_path = os.path.join(image_path, uuid, "img")
    if os.path.exists(gif_folder_path):
        return False
    logging.info(f"Creating gif for {uuid}.")
    img_files = [f for f in os.listdir(img_folder_path) if f.endswith(".jpg")]
    num_images = len(img_files)
    frames = [
        Image.open(os.path.join(img_folder_path, f"{str(i).zfill(INDEX_LENGTH)}.jpg"))
        for i in range(1, num_images + 1)
    ]

    os.makedirs(gif_folder_path, exist_ok=True)
    gif_path = os.path.join(gif_folder_path, f"{uuid}.gif")
    frames[0].save(
        gif_path,
        format="GIF",
        append_images=frames[1:],
        save_all=True,
        duration=100,
        loop=0,
    )
    logging.info(f"Created gif for {uuid}.")
    return True

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <uuid> <image_path>")
        sys.exit(1)

    uuid = sys.argv[1]
    image_path = sys.argv[2]

    # Example of calling make_gif using system arguments
    make_gif(uuid, image_path)
    print("GIF creation complete.")
