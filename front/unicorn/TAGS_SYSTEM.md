# Иерархическая система тегов для вакансий

## Обзор

Новая система тегов использует канонический подход с трехуровневой иерархией:

1. **Категория** (Primary) - основные сферы: IT, GameDev, Startup
2. **Специализация** (Secondary) - конкретные направления: Backend, Unity, Product
3. **Детали** (Tertiary) - уточняющие теги: Backend+Node.js, Unity+C#, Product+Manager

## Структура иерархии

### IT
```
IT
├── Backend
│   ├── Backend+Node.js
│   ├── Backend+Python
│   ├── Backend+Go
│   ├── Backend+C#
│   ├── Backend+Java
│   ├── Backend+PHP
│   └── Backend+Ruby
├── Frontend
│   ├── Frontend+React
│   ├── Frontend+Vue
│   ├── Frontend+Angular
│   ├── Frontend+Next.js
│   └── Frontend+Nuxt
├── Fullstack
│   ├── Fullstack+MERN
│   ├── Fullstack+MEAN
│   └── Fullstack+Django
├── Mobile
│   ├── Mobile+iOS
│   ├── Mobile+Android
│   ├── Mobile+React Native
│   └── Mobile+Flutter
├── DevOps
│   ├── DevOps+AWS
│   ├── DevOps+Azure
│   ├── DevOps+GCP
│   ├── DevOps+Kubernetes
│   └── DevOps+Docker
├── Data
│   ├── Data+Analytics
│   ├── Data+Science
│   ├── Data+Engineering
│   └── Data+ML/AI
├── QA
│   ├── QA+Manual
│   ├── QA+Automation
│   └── QA+Performance
├── Security
│   ├── Security+InfoSec
│   ├── Security+AppSec
│   └── Security+DevSecOps
└── Design
    ├── Design+UX
    ├── Design+UI
    └── Design+Product
```

### GameDev
```
GameDev
├── Unity
│   ├── Unity+C#
│   ├── Unity+Mobile
│   └── Unity+PC
├── Unreal Engine
│   ├── UE5+C++
│   ├── UE5+C#
│   └── UE5+Blueprints
├── Game Design
│   ├── Game Design+Level
│   ├── Game Design+Narrative
│   └── Game Design+Systems
├── Game Art
│   ├── Game Art+3D
│   ├── Game Art+2D
│   ├── Game Art+Character
│   ├── Game Art+Environment
│   ├── Game Art+VFX
│   └── Game Art+UI
├── Game Audio
│   ├── Game Audio+Sound Design
│   ├── Game Audio+Music
│   └── Game Audio+Technical
└── Technical Art
    ├── Technical Art+Shaders
    ├── Technical Art+Pipeline
    └── Technical Art+Tools
```

### Startup
```
Startup
├── Product
│   ├── Product+Manager
│   ├── Product+Owner
│   └── Product+Analyst
├── Growth
│   ├── Growth+Marketing
│   ├── Growth+Product
│   └── Growth+Sales
└── Operations
    ├── Operations+Manager
    └── Operations+Analyst
```

### Support
```
Support
├── Customer Support
│   ├── Support+Chat
│   ├── Support+Phone
│   ├── Support+Email
│   ├── Support+Level 1
│   ├── Support+Level 2
│   └── Support+Level 3
└── Technical Support
    ├── Tech Support+IT
    ├── Tech Support+Hardware
    └── Tech Support+Software
```

## Примеры использования

### Пример 1: GameDev разработчик на UE5 с C#
Теги: `['GameDev', 'Unreal Engine', 'UE5+C#']`

Путь: GameDev > Unreal Engine > UE5+C#

### Пример 2: Backend разработчик на Node.js
Теги: `['IT', 'Backend', 'Backend+Node.js']`

Путь: IT > Backend > Backend+Node.js

### Пример 3: Support специалист 2 уровня с телефонией
Теги: `['Support', 'Customer Support', 'Support+Phone', 'Support+Level 2']`

Путь: Support > Customer Support > Support+Phone, Support+Level 2

## API функций

### `getChildTags(tagValue: string): string[]`
Получает все дочерние теги для указанного тега.

```typescript
getChildTags('IT') 
// => ['Backend', 'Frontend', 'Backend+Node.js', 'Frontend+React', ...]
```

### `getParentTag(tagValue: string): string | null`
Получает родительский тег для указанного тега.

```typescript
getParentTag('Backend+Node.js') // => 'Backend'
getParentTag('Backend') // => 'IT'
```

### `getTagPath(tagValue: string): string[]`
Получает полный путь от корня до указанного тега.

```typescript
getTagPath('UE5+C#') // => ['GameDev', 'Unreal Engine', 'UE5+C#']
```

### `searchTags(query: string): string[]`
Поиск тегов по ключевым словам с учетом иерархии.

```typescript
searchTags('C#') 
// => ['Backend+C#', 'Unity+C#', 'UE5+C#']
```

### `matchesTagFilter(vacancyTags: string[], filterTags: string[]): boolean`
Проверяет, соответствуют ли теги вакансии фильтру с учетом иерархии.

```typescript
// Если в фильтре выбран 'Backend', вакансии с 'Backend+Node.js' тоже пройдут
matchesTagFilter(['Backend+Node.js'], ['Backend']) // => true
matchesTagFilter(['Frontend+React'], ['Backend']) // => false
```

## Компоненты

### `<HierarchicalTagSelector>`
Компонент для выбора тегов с иерархией.

```vue
<HierarchicalTagSelector
  :selected-tags="selectedTags"
  :selected-category="selectedCategory"
  :selected-specialization="selectedSpecialization"
  :enable-search="true"
  :show-recommendations="true"
  @toggle-tag="handleToggleTag"
  @select-category="handleSelectCategory"
  @select-specialization="handleSelectSpecialization"
  @clear-all="handleClearAll"
/>
```

### Композабл `useTags()`
Хук для работы с системой тегов.

```typescript
const tags = useTags()

// Выбрать категорию
tags.selectCategory('IT')

// Выбрать специализацию
tags.selectSpecialization('Backend')

// Добавить/удалить тег
tags.toggleTag('Backend+Node.js')

// Очистить все теги
tags.clearAllTags()

// Получить отображаемое имя
tags.getTagDisplayName('Backend+Node.js') // => 'Node.js'

// Проверить соответствие
tags.checkTagMatch(['Backend+Node.js']) // => boolean
```

## Фильтрация вакансий

### Иерархическая фильтрация

Если пользователь выбирает родительский тег, показываются все вакансии с дочерними тегами:

- Выбран `IT` → показываются вакансии с `Backend`, `Frontend`, `Backend+Node.js`, `Frontend+React` и т.д.
- Выбран `Backend` → показываются вакансии с `Backend+Node.js`, `Backend+Python`, `Backend+Go` и т.д.
- Выбран `Backend+Node.js` → показываются только вакансии с этим конкретным тегом

### Комбинирование тегов

Теги могут повторяться в разных сферах:
- `C#` может быть в `Backend+C#` (IT) и `UE5+C#` (GameDev)
- При поиске по `C#` показываются все вакансии с любым вариантом

## Приоритет отображения

Группы тегов имеют приоритет (priority):
1. **Основные категории** (priority: 1): IT, GameDev, Startup
2. **Специализации** (priority: 2): Backend, Unity, Product
3. **Детальные теги** (priority: 3): Backend+Node.js, Unity+C#

Фильтры отображаются в порядке приоритета для лучшего UX.

## Локализация

Все теги поддерживают i18n через labelKey:

```json
{
  "tags": {
    "categories": {
      "it": "IT",
      "gamedev": "GameDev",
      "startup": "Стартап"
    },
    "specializations": {
      "backend": "Backend-разработка",
      "unity": "Unity"
    },
    "details": {
      "backendNodejs": "Backend (Node.js)",
      "unityCsharp": "Unity (C#)"
    }
  }
}
```

## Миграция со старой системы

Старые теги остаются совместимыми. Рекомендуется постепенно переводить вакансии на новую иерархическую систему.

Пример миграции:
- `Backend` → `IT > Backend`
- `Node.js` → `IT > Backend > Backend+Node.js`
- `Unity` → `GameDev > Unity`
- `C#` + `Unity` → `GameDev > Unity > Unity+C#`
