#! /bin/sh

git add .
if [ $# -ge 1 ]; then
    git commit . -m "$1"
else
    git commit .
fi
git push
