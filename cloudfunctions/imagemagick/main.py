import os
import tempfile

from google.cloud import storage
from wand.image import Image

storage_client = storage.Client()

def resize_image(data, context):
    file_data = data

    file_name = file_data['name']
    bucket_name = file_data['bucket']

    blob = storage_client.bucket(bucket_name).get_blob(file_name)
    blob_uri = f'gs://{bucket_name}/{file_name}'
    blob_source = {'source': {'image_uri': blob_uri}}

    if file_name.startswith('thumb-'):
        print(f'The image {file_name} is already resized.')
        return

    _, temp_local_filename = tempfile.mkstemp()

    blob.download_to_filename(temp_local_filename)
    print(f'Image {file_name} was downloaded to {temp_local_filename}.')

    with Image(filename=temp_local_filename) as image:
        image.transform(resize='x128')
        image.save(filename=temp_local_filename)

        print(f'Image {file_name} ({image.mimetype}) was resized.')

        new_file_name = f'thumb-{file_name}'
        new_blob = blob.bucket.blob(new_file_name)
        new_blob.upload_from_filename(temp_local_filename, content_type=image.mimetype)
        print(f'resized image was uploaded to {new_file_name}.')

    os.remove(temp_local_filename)
