#!/bin/zsh




# 增加文件到缓存区
git add .

# 写注释
echo "输入新的注释后【回车】确认:"

read newCommit

git commit -m "${newCommit}"
# 推送
git push

oldVersion=""

for k in $(git tag -l); do
  # 循环删除旧版本
  git tag -d "${k}"
  # 删除远程版本
  git push origin :refs/tags/"${k}"
  #
  oldVersion="${k}"
  echo "${k} success Delete! "
done


arr=$(echo "${oldVersion:1}" | awk 'BEGIN{FS=".";OFS=" "} {print $1,$2,$3,$4}');


hundred=$((arr[1]))
ten=$((arr[2]))
one=$((arr[3]))

if [ $one -lt 50 ] ;then
  one=$((one + 1));
else
  one=0;
  if [ $ten -lt 50 ]; then
    ten=$((ten + 1));
  else
    ten=0;
    hundred=$((hundred + 1))
  fi;

fi;

version="v$hundred.$ten.$one"

echo "旧版本【${oldVersion}】>> 新版本 【${version}】"

# 新版本名
git tag "$version"
  # 推送tag
git push origin --tags











####*********************************************************************************************
#
## 1.为避免冲突需要先同步下远程仓库
#git pull
#
## 2.在本地项目目录下删除缓存
#git rm -r --cached .
#
## 3.再次add所有文件，输入以下命令，再次将项目中所有文件添加到本地仓库缓存中
#git add .
#
## 4.添加commit，提交到远程库
#git commit -m "filter new files"
#git push
