#!/usr/bin/env bash
# scripts/release.sh — 一键发版脚本。
#
# 工作流：
#   1. 校验工作区干净、当前在 main/feat 分支
#   2. 校验传入的版本号格式
#   3. 同步更新 package.json 的 version 与 helios.binaryVersion
#   4. 提交版本 bump、打 tag、推送
#   5. 远端 GitHub Actions 会捕获 tag 并跑 GoReleaser
#
# 用法：
#   ./scripts/release.sh 0.1.0
#   ./scripts/release.sh 0.2.0 --dry-run    # 只演练，不推送

set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <version> [--dry-run]" >&2
  echo "  version: semver without leading 'v', e.g. 0.1.0" >&2
  exit 1
fi

VERSION="$1"
DRY_RUN="${2:-}"
TAG="v${VERSION}"

if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
  echo "error: version must be semver (e.g. 0.1.0 or 0.1.0-rc.1), got: $VERSION" >&2
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

# 1. 工作区必须干净
if [[ -n "$(git status --porcelain)" ]]; then
  echo "error: working tree not clean; commit or stash first" >&2
  git status --short >&2
  exit 1
fi

# 2. tag 不能已经存在
if git rev-parse -q --verify "refs/tags/${TAG}" >/dev/null; then
  echo "error: tag ${TAG} already exists" >&2
  exit 1
fi

# 3. 同步 package.json 中的两个版本字段
PKG="$REPO_ROOT/package.json"
echo "==> updating package.json version -> ${VERSION}"
node - "$PKG" "$VERSION" <<'NODE'
const fs = require('fs');
const [, , file, ver] = process.argv;
const pkg = JSON.parse(fs.readFileSync(file, 'utf8'));
pkg.version = ver;
pkg.helios = pkg.helios || {};
pkg.helios.binaryVersion = ver;
fs.writeFileSync(file, JSON.stringify(pkg, null, 2) + '\n');
NODE

# 4. 提交 + 打 tag（如果 package.json 没有变化就跳过 commit）
git add package.json
if git diff --cached --quiet -- package.json; then
  echo "==> package.json already at ${VERSION}, skipping commit"
else
  git commit -m "chore(release): ${TAG}"
fi
git tag -a "${TAG}" -m "${TAG}"

# 5. 推送（除非 dry-run）
if [[ "$DRY_RUN" == "--dry-run" ]]; then
  echo "==> dry-run mode, NOT pushing. Inspect with:"
  echo "    git log -1"
  echo "    git show ${TAG}"
  echo "    # to undo tag:    git tag -d ${TAG}"
  echo "    # to undo commit: git reset --hard HEAD~1   (only if a release commit was created)"
  exit 0
fi

CURRENT_BRANCH="$(git branch --show-current)"
echo "==> pushing branch ${CURRENT_BRANCH} and tag ${TAG}"
git push origin "${CURRENT_BRANCH}"
git push origin "${TAG}"

echo
echo "==> done. GitHub Actions will now build and publish release ${TAG}."
echo "    https://github.com/auXiaoYuan/helios-cli/actions"
echo "    https://github.com/auXiaoYuan/helios-cli/releases/tag/${TAG}"
