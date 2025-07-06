# viberules

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey.svg)]()

[English](README.md) | [한국어](README.ko.md)

> 심볼릭 링크를 활용한 AI 어시스턴트 규칙 관리 도구

⚠️ **참고**: Windows는 심볼릭 링크 제한으로 인해 지원되지 않습니다. Windows에서는 WSL2 사용을 고려하세요.

viberules는 심볼릭 링크를 사용하여 AI 코딩 어시스턴트 규칙(Claude Code, Amazon Q Developer, Gemini Code Assist 등)을 실시간으로 동기화하는 CLI 도구입니다.

## ✨ 주요 기능

- 🎯 **통합 관리**: 모든 AI 도구 규칙을 .viberules/ 폴더의 단일 파일로 관리
- 🔄 **실시간 동기화**: 심볼릭 링크를 통해 변경사항 자동 반영
- 🛠️ **개별 타겟 제어**: 특정 AI 도구만 활성화/비활성화 가능
- 🏠 **유연한 모드**: 팀 공유를 위한 public 모드 또는 개인 규칙을 위한 local 모드
- 📝 **스마트 .gitignore**: 모드 인식 정책으로 gitignore 자동 관리
- 🌍 **크로스 플랫폼**: macOS와 Linux 지원

## 🚀 빠른 시작

### 설치

```bash
# GitHub에서 설치
go install github.com/sky1core/viberules@latest
```

### 프로젝트 초기화

```bash
# 프로젝트 루트에서 실행
viberules init

# 강제 재초기화 (기존 rules.md 보존)
viberules init --force
```

다음 파일들이 생성됩니다:
- `.viberules/rules.md` - 모든 AI 도구를 위한 단일 규칙 파일
- 각 AI 도구용 심볼릭 링크 (CLAUDE.md, GEMINI.md, AGENTS.md, .amazonq/rules/AMAZONQ.md)
- 모드 인식 정책이 적용된 `.gitignore` 업데이트

### 타겟 관리

```bash
# 활성화된 타겟 목록 보기
viberules list

# 불필요한 타겟 제거
viberules remove amazonq

# 타겟 다시 추가
viberules add amazonq

# 프로젝트 모드 설정
viberules mode public   # .viberules를 팀과 공유
viberules mode local    # .viberules를 비공개로 유지
```

## 📋 지원 AI 도구

| AI 도구 | 타겟 이름 | 출력 파일 |
|---------|-------------|--------------|
| Claude Code | `claude` | `CLAUDE.md` |
| Amazon Q Developer | `amazonq` | `.amazonq/rules/AMAZONQ.md` |
| Gemini Code Assist | `gemini` | `GEMINI.md` |
| 범용 AI 도구/Codex | `codex` | `AGENTS.md` |

## 🛠️ 명령어

```bash
# 프로젝트 초기화
viberules init

# 기존 프로젝트 재초기화 (rules.md 보존)
viberules init --force

# 활성화된 타겟 목록
viberules list

# 타겟 추가/제거
viberules add claude
viberules remove amazonq

# 프로젝트 모드 관리
viberules mode          # 현재 모드 표시
viberules mode public   # public 모드로 설정 (팀 공유)
viberules mode local    # local 모드로 설정 (비공개)

# 도움말
viberules --help
```

## 📁 프로젝트 구조

`viberules init` 실행 후:

```
your-project/
├── .viberules/              # 설정 디렉토리
│   ├── rules.md             # 모든 AI 도구를 위한 단일 규칙 파일
│   └── .config.yaml         # 설정 파일 (모드 & 타겟, git에서 무시됨)
├── .gitignore               # 모드에 따라 자동 업데이트
├── CLAUDE.md                # .viberules/rules.md로의 심볼릭 링크
├── GEMINI.md                # .viberules/rules.md로의 심볼릭 링크
├── AGENTS.md                # .viberules/rules.md로의 심볼릭 링크
└── .amazonq/
    └── rules/
        └── AMAZONQ.md       # ../../.viberules/rules.md로의 심볼릭 링크
```

### 모드별 .gitignore 동작

viberules는 규칙 공유 방식을 제어하는 두 가지 모드를 지원합니다:

**Local 모드** (기본값, 개인 규칙용):
- 전체 `.viberules/` 디렉토리가 git에서 무시됨
- 모든 규칙이 로컬 머신에만 유지됨
- 출력 파일(CLAUDE.md 등)이 무시됨
- 개인 설정이나 민감한 정보가 포함된 규칙에 사용

**Public 모드** (팀 협업용):
- `.viberules/rules.md`가 git에서 추적됨 (팀과 공유)
- `.viberules/.config.yaml`은 항상 무시됨 (개인 설정)
- 출력 파일(CLAUDE.md 등)이 무시됨
- AI 어시스턴트 규칙을 팀과 공유하고 싶을 때 사용

## ⚙️ 작동 원리

1. **규칙 편집**: `.viberules/rules.md` 수정 (단일 소스)
2. **즉시 동기화**: 심볼릭 링크를 통해 모든 AI 도구에 자동 반영
3. **모드 인식 Git**: Public 모드는 팀과 규칙 공유, local 모드는 모든 것을 비공개로 유지
4. **스마트 타겟팅**: 사용하는 AI 도구만 활성화

## 🔧 고급 사용법

### 프로젝트 모드

**Public 모드** (팀 프로젝트에 권장):
```bash
viberules mode public
```
- `.viberules/rules.md`가 커밋되고 팀과 공유됨
- 개인 설정(타겟 구성)은 로컬에 유지됨

**Local 모드** (개인 프로젝트용):
```bash
viberules mode local
```
- 전체 `.viberules/` 디렉토리가 git에서 무시됨
- 규칙이 완전히 비공개로 유지됨

### 효과적인 규칙 작성

`.viberules/rules.md` 편집:

```markdown
# AI 어시스턴트 규칙

## 프로젝트 개요
Next.js와 Tailwind CSS를 사용하는 TypeScript React 프로젝트입니다.

## 코딩 표준
- TypeScript를 strict 모드로 사용
- ESLint 설정 준수
- 모든 함수에 대한 단위 테스트 작성
- 설명적인 변수명 사용

## 아키텍처 가이드라인
- 클린 아키텍처 원칙 준수
- 비즈니스 로직과 UI 컴포넌트 분리
- 상태 관리를 위한 커스텀 훅 사용

## API 가이드라인
- 적절한 HTTP 메서드로 REST API 사용
- 적절한 에러 처리 구현
- API 응답에 TypeScript 인터페이스 사용
```

### 타겟 관리

```bash
# Claude만으로 시작
viberules remove amazonq
viberules remove gemini

# 나중에 다른 도구 추가
viberules add gemini
```

## 🧪 개발

### 필요 조건

- Go 1.21 이상
- macOS 또는 Linux (Windows 미지원)

### 빌드

```bash
# 클론 및 빌드
git clone https://github.com/sky1core/viberules.git
cd viberules
go build .
```

### 테스트

```bash
# 모든 테스트 실행
go test ./...

# 커버리지와 함께 실행
go test ./... -cover

# 특정 테스트 실행
go test -v -run TestCompleteViberulesWorkflow .
```

---

<p align="center">
  Created by <a href="https://github.com/sky1core">sky1core</a>
</p>