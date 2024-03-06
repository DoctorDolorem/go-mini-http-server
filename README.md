# Flags 

| Flag  | Description                  | Example                   |
|-------|------------------------------|---------------------------|
| -port | Define the port to listen on | -port 1234        |
| -dir  | Define the root directory  | -dir "C:\Users\Docs"
| -upload | Define upload directory | -upload myuploads 

# Usage
By default, the server listens on port 9000 and starts from the current directory it's executed from.

To upload files, you can use cURL. The upload  form is named "file".

curl -F "file=@/path/to/file" URL/upload_directory


