#!/usr/bin/env bash
set -euo pipefail

# Usage: bash scripts/squash-pr.sh [branch] [base] [commit message]
# Rewrites the branch and immediately updates its PR on origin.
if [[ ${1:-} == "--help" || $# -gt 3 ]]; then
  echo 'Usage: bash scripts/squash-pr.sh [branch=tasks] [base=main] [commit message]'
  exit 0
fi

branch=${1:-tasks}
base=${2:-main}
message=${3:-feat: get tasks with filtering and pagination}

cd "$(git rev-parse --show-toplevel)"

if [[ "$branch" == "$base" ]]; then
  echo 'Ошибка: ветка PR и базовая ветка должны различаться.' >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo 'Ошибка: сначала закоммить или сохрани в stash текущие изменения.' >&2
  exit 1
fi

git switch "$branch"
git fetch origin

base_sha=$(git rev-parse --verify "refs/remotes/origin/$base^{commit}")
remote_sha=$(git rev-parse --verify "refs/remotes/origin/$branch^{commit}")

if ! git merge-base --is-ancestor "$remote_sha" HEAD; then
  echo 'Ошибка: сначала объедини изменения из удалённой ветки PR с локальной.' >&2
  exit 1
fi

if git diff --quiet "$base_sha" HEAD; then
  echo 'Нет изменений относительно базовой ветки для сквоша.' >&2
  exit 1
fi

backup="backup/$branch-before-squash-$(date +%Y%m%d-%H%M%S)-$$"
git branch "$backup"
echo "Резервная копия: $backup"

git reset --soft "$base_sha"
git commit -m "$message"
git push "--force-with-lease=refs/heads/$branch:$remote_sha" origin "HEAD:refs/heads/$branch"

echo "Готово: ветка $branch содержит один коммит поверх origin/$base; PR обновлён."
