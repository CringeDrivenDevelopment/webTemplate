#!/bin/bash

echo "Setting up hook"
cp ./pre-commit.hook .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
echo "Hook set up"