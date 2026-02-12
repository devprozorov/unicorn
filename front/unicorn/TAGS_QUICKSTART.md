# Обновление системы тегов - краткая инструкция

## Что изменилось

Реализована **иерархическая система тегов** с тремя уровнями:

1. **Категория** - IT, GameDev, Startup, Support, Marketing, Management, HR, Finance
2. **Специализация** - Backend, Unity, Product, Customer Support и т.д.
3. **Детали** - Backend+Node.js, Unity+C#, Support+Level 2 и т.д.

## Основные преимущества

✅ **Иерархия**: GameDev > Unity > Unity+C#  
✅ **Умная фильтрация**: Выбор "Backend" показывает все вакансии с Backend+Node.js, Backend+Python и т.д.  
✅ **Переиспользование**: C# может быть в Backend+C# (IT) и Unity+C# (GameDev)  
✅ **Приоритет**: Основные категории (IT, GameDev, Startup) показываются первыми  
✅ **Поиск**: Поиск по всем уровням иерархии с учетом searchTerms  

## Быстрый старт

### 1. Использование композабла

```typescript
import { useTags } from '~/composables/useTags'

const tags = useTags()

// Выбрать категорию
tags.selectCategory('IT')

// Выбрать специализацию  
tags.selectSpecialization('Backend')

// Добавить детальный тег
tags.toggleTag('Backend+Node.js')

// Получить теги для API
const selectedTags = tags.getSelectedTagsForApi()
// => ['IT', 'Backend', 'Backend+Node.js']
```

### 2. Использование компонента

```vue
<HierarchicalTagSelector
  :selected-tags="selectedTags"
  :enable-search="true"
  :show-recommendations="true"
  @toggle-tag="handleToggleTag"
  @select-category="handleSelectCategory"
  @clear-all="handleClearAll"
/>
```

### 3. Фильтрация вакансий

```typescript
import { matchesTagFilter } from '~/utils/tags'

const vacancy = {
  tags: ['IT', 'Backend', 'Backend+Node.js']
}

// Фильтр по IT - вакансия будет показана
matchesTagFilter(vacancy.tags, ['IT']) // => true

// Фильтр по Backend - вакансия будет показана
matchesTagFilter(vacancy.tags, ['Backend']) // => true

// Фильтр по Frontend - вакансия НЕ будет показана
matchesTagFilter(vacancy.tags, ['Frontend']) // => false
```

## Файлы

### Созданные файлы
- `unicorn/app/composables/useTags.ts` - композабл для работы с тегами
- `unicorn/app/components/tags/HierarchicalTagSelector.vue` - компонент выбора тегов
- `TAGS_SYSTEM.md` - подробная документация
- `TAGS_EXAMPLES.md` - примеры использования

### Обновленные файлы
- `unicorn/app/utils/tags.ts` - добавлена иерархия и функции работы с ней
- `unicorn/app/pages/jobs/index.vue` - обновлена фильтрация с учетом иерархии
- `unicorn/app/pages/company/vacancies.vue` - можно обновить аналогично

## Примеры иерархий

### IT разработчик на Node.js
```
IT > Backend > Backend+Node.js
```
Теги: `['IT', 'Backend', 'Backend+Node.js']`

### GameDev на UE5 с C++
```
GameDev > Unreal Engine > UE5+C++
```
Теги: `['GameDev', 'Unreal Engine', 'UE5+C++']`

### Support специалист 2 уровня
```
Support > Customer Support > Support+Level 2
```
Теги: `['Support', 'Customer Support', 'Support+Level 2']`

## Миграция

Старые теги остаются совместимыми. Рекомендации:

**Было**: `['Backend', 'Node.js', 'JavaScript']`  
**Стало**: `['IT', 'Backend', 'Backend+Node.js']`

**Было**: `['Unity', 'C#']`  
**Стало**: `['GameDev', 'Unity', 'Unity+C#']`

## API функций

| Функция | Описание |
|---------|----------|
| `getChildTags(tag)` | Получить все дочерние теги |
| `getParentTag(tag)` | Получить родительский тег |
| `getTagPath(tag)` | Получить путь от корня |
| `searchTags(query)` | Поиск тегов по запросу |
| `matchesTagFilter(vacancyTags, filterTags)` | Проверка соответствия фильтру |
| `formatTagDisplay(tag)` | Форматировать для отображения |
| `getTagFullDescription(tag)` | Получить полное описание |

## Следующие шаги

1. ✅ Базовая иерархия реализована
2. ✅ Поиск и фильтрация работают
3. ✅ Компоненты созданы
4. ⏳ Добавить локализацию для новых тегов
5. ⏳ Обновить страницу создания вакансий
6. ⏳ Мигрировать существующие вакансии

## Вопросы?

Смотрите подробную документацию:
- [TAGS_SYSTEM.md](./TAGS_SYSTEM.md) - полное описание системы
- [TAGS_EXAMPLES.md](./TAGS_EXAMPLES.md) - примеры использования
