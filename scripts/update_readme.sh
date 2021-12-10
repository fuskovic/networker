#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)

echo "***** INSTALLING LATEST NETWORKER GLOBALLY FROM SOURCE *****"

# install
$PROJECT_ROOT/scripts/install.sh

echo "***** UPDATING README *****"

writeCmdUsage () {
    echo -e "\n## $1 \n" >> $PROJECT_ROOT/new_readme.md
    echo "\`\`\`"  >> $PROJECT_ROOT/new_readme.md
    networker $2 --help 2>> $PROJECT_ROOT/new_readme.md
    echo "\`\`\`"  >> $PROJECT_ROOT/new_readme.md
}

# create new readme
touch new_readme.md
echo "$(cat $PROJECT_ROOT/README.md | head -26)" > $PROJECT_ROOT/new_readme.md
echo -e "\n# Usage \n" >> $PROJECT_ROOT/new_readme.md
echo "\`\`\`"  >> $PROJECT_ROOT/new_readme.md
networker 2>> $PROJECT_ROOT/new_readme.md
echo "\`\`\`"  >> $PROJECT_ROOT/new_readme.md
echo -e "\n# Commands " >> $PROJECT_ROOT/new_readme.md

writeCmdUsage "List" "ls"
writeCmdUsage "Lookup" "lu"
writeCmdUsage "Scan" "s"
writeCmdUsage "Request" "r"

# delete old one
rm $PROJECT_ROOT/README.md

# rename
mv $PROJECT_ROOT/new_readme.md $PROJECT_ROOT/README.md

echo "***** UPDATED README *****"