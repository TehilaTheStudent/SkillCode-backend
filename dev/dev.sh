export MODE_ENV=development
# Dynamically set PROJECT_ROOT to the current directory
export PROJECT_ROOT=$(pwd)

echo "MODE_ENV set to $MODE_ENV"

pip install -r template-assets/python/requirements.txt

npm install --prefix template-assets/javascript/

# Delete files with names longer than 32 characters in both directories
for dir in template-assets/python/ template-assets/javascript/; do
  find "$dir" -type f -name '????????????????????????????????*' -exec rm {} \;
done


