###
Read all projects from gitlab group and generate a shell file.
```

 git init
 git remote add main-repo git@github.com:qs-wang/main-repo.git
 git remote add repo1  git@gitlab.com/repo1.git
 git remote add repo2  git@gitlab.com/repo2.git
 git fetch --all --no-tags

 ~/git/monorepo-tools/monorepo_build.sh  main-repo repo1 repo2

 ```