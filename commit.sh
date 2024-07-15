#! /bin/sh

if [ $# -ne 1 ]; then
    echo "expected '[sh ./commit.sh | ./commit.sh] [commit message]'"
    exit 1
fi

git add .
if [ "$1" = "-m" ]; then
    git commit . -m "$1"
fi
git push